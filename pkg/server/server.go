// Copyright (C) 2023 myl7
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"log"
	"sync"

	"github.com/go-redis/redis/v9"
	"github.com/myl7/karaoke/pkg/rpc"
	"github.com/myl7/karaoke/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	rpc.UnimplementedServerRPCServer

	c ServerConfig

	// Set in bootstrap
	addr     string
	id       string
	pCs      map[string]rpc.PConfig
	layerIdx map[int][]string

	round int

	// pool contains encoded onion bytes received
	pool [][]byte
	// poolMask indicate which server has sent onions here in the round
	poolMask map[string]bool
	poolLock sync.Mutex
	// poolFullCh notifys runRound to process onions in pool
	poolFullCh chan bool

	// TODO: Dead drop external persistence
	deadDrop     map[string][]byte
	deadDropLock sync.Mutex

	rC  *redis.Client
	mDB *mongo.Database

	postLayerClients map[string]rpc.ServerRPCClient
	zeroLayerClients map[string]rpc.ServerRPCClient

	clients []Client
}

type ServerConfig struct {
	Layer int
	// Addr is the endpoint accessible by other servers.
	// If empty, use the port of [LAddr] and the public IP by detection.
	// At that time [LAddr] must be in ":%d" format.
	Addr string
	// LAddr is for gRPC Listening
	LAddr string
	// PK, SK can be left empty to atomically gen
	// TODO: PK/SK rotation
	PK, SK *[32]byte

	// RAddr is Redis addr
	RAddr string
	// MURI is MongoDB URI
	MURI string
}

func NewServer(c ServerConfig) *Server {
	if c.PK == nil || c.SK == nil {
		var err error
		c.PK, c.SK, err = box.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
	}

	return &Server{c: c}
}

func (s *Server) Run(ctx context.Context) error {
	sub := s.rC.Subscribe(ctx, "karaoke/round")
	defer func() {
		err := sub.Close()
		if err != nil {
			panic(err)
		}
	}()
	ch := sub.Channel()
L:
	for {
		select {
		case rSB := <-ch:
			var rS rpc.RoundStartMsg
			err := json.Unmarshal([]byte(rSB.Payload), &rS)
			if err != nil {
				panic(err)
			}
			round := rS.Round

			if round <= s.round {
				log.Println("invalid round num")
				continue L
			} else {
				s.round = round
			}

			err = s.runRound(ctx)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *Server) runRound(ctx context.Context) error {
	log.Println("start round")

	var cipherBs [][]byte
	if s.c.Layer > 0 {
		<-s.poolFullCh
		s.poolLock.Lock()
		cipherBs = s.pool
		s.pool = nil
		s.poolLock.Unlock()
	} else {
		for _, c := range s.clients {
			msg, err := c.RunRound()
			if err != nil {
				return err
			}
			cipherBs = append(cipherBs, msg)
		}
	}

	log.Println("get to-process onions")
	// Do onion processing here other than in handler for easier time monitoring

	var plainBs [][]byte
	for _, b := range cipherBs {
		plainB, ok := utils.BoxDec(b, s.c.PK, s.c.SK)
		if ok {
			plainBs = append(plainBs, plainB)
		}
	}

	// Dedup
	m := make(map[string][]byte)
	for _, b := range plainBs {
		m[string(utils.Hash(b))] = b
	}
	plainBs = nil
	for _, b := range m {
		plainBs = append(plainBs, b)
	}

	var os []*rpc.Onion
	for _, b := range plainBs {
		o := new(rpc.Onion)
		err := proto.Unmarshal(b, o)
		if err != nil {
			continue
		}

		ok := true
		if o.NextHop != "" && o.DeadDrop != "" {
			ok = false
		}
		if o.NextHop != "" {
			pC, ok := s.pCs[o.NextHop]
			if !ok {
				ok = false
			}
			if pC.Layer != s.c.Layer+1 {
				ok = false
			}
		}
		if !ok {
			log.Println("invalid received onion")
			continue
		}

		os = append(os, o)
	}

	// Avoid empty Bloom
	if len(os) == 0 {
		return nil
	}

	// TODO: Gen noise or check Bloom

	log.Println("get to-send onions")

	oM := make(map[string][][]byte)
	s.deadDropLock.Lock()
	for _, o := range os {
		if o.DeadDrop != "" {
			s.deadDrop[o.DeadDrop] = o.Body
		} else {
			oM[o.NextHop] = append(oM[o.NextHop], o.Body)
		}
	}
	s.deadDropLock.Unlock()

	g, gCtx := errgroup.WithContext(ctx)
	for _, id := range s.layerIdx[s.c.Layer+1] {
		i := id
		c := s.postLayerClients[id]
		g.Go(func() error {
			cc, err := c.FwdOnions(gCtx)
			if err != nil {
				return err
			}

			err = cc.Send(&rpc.OnionMsg{
				Msg: &rpc.OnionMsg_Meta{
					Meta: &rpc.OnionMsgMeta{
						Id: s.id,
					},
				},
			})
			if err != nil {
				return err
			}
			for _, b := range oM[i] {
				err = cc.Send(&rpc.OnionMsg{
					Msg: &rpc.OnionMsg_Body{
						Body: b,
					},
				})
				if err != nil {
					return err
				}
			}
			_, err = cc.CloseAndRecv()
			return err
		})
	}
	err := g.Wait()
	if err != nil {
		return err
	}

	s.poolLock.Lock()
	s.poolMask = make(map[string]bool)
	s.poolLock.Unlock()

	rE := rpc.RoundEndMsg{
		Round: s.round,
		ID:    s.id,
	}
	rEB, err := json.Marshal(rE)
	if err != nil {
		panic(err)
	}
	err = s.rC.Publish(ctx, "karaoke/round_ok", rEB).Err()
	if err != nil {
		return err
	}
	log.Println("end round")
	return nil
}

var ErrBloomCheckFailure = errors.New("fail to check Bloom, which means noise loss and adversary possibility")

func (s *Server) FwdOnions(ss rpc.ServerRPC_FwdOnionsServer) error {
	metaMsg, err := ss.Recv()
	if err != nil {
		return err
	}
	meta := metaMsg.Msg.(*rpc.OnionMsg_Meta).Meta

	ok, err := func() (bool, error) {
		s.poolLock.Lock()
		defer s.poolLock.Unlock()

		if s.poolMask[meta.Id] {
			return false, nil
		}
		s.poolMask[meta.Id] = true

		var bs [][]byte
		for {
			msg, err := ss.Recv()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return false, err
				}
			}
			bs = append(bs, msg.Msg.(*rpc.OnionMsg_Body).Body)
		}

		ok := true
		for _, id := range s.layerIdx[s.c.Layer-1] {
			if !s.poolMask[id] {
				ok = false
				break
			}
		}

		s.pool = append(s.pool, bs...)
		return ok, nil
	}()
	if ok {
		s.poolFullCh <- true
	}
	if err != nil {
		return err
	}
	return ss.SendAndClose(&rpc.FwdOnionsRes{})
}
