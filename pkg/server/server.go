package server

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v9"
	"github.com/myl7/karaoke/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative server.proto

type Server struct {
	UnimplementedRPCServer

	c ServerConfig

	// Set in bootstrap
	addr     string
	id       string
	pCs      map[string]PConfig
	layerIdx map[int][]string

	round int

	// Pool of encoded onion bytes received
	pool [][]byte
	// Indicate which other servers have sent onions here
	poolMask map[string]bool
	poolLock sync.Mutex
	// Notify [runRound] to process onions in pool
	poolFullCh chan bool

	// TODO: Dead drop new persistent storage
	deadDrop     map[string][]byte
	deadDropLock sync.Mutex

	rC  *redis.Client
	mDB *mongo.Database

	postLayerClients map[string]RPCClient
	zeroLayerClients map[string]RPCClient

	clients []Client
}

type ServerConfig struct {
	Layer int
	// Endpoint accessible by other servers.
	// If empty, use the port of [LAddr] and the public IP by detection.
	// At that time [LAddr] must be in ":%d" format.
	Addr string
	// For gRPC Listening
	LAddr string
	// Can leave empty to atomically gen
	// TODO: PK/sk rotation
	PK, SK *[32]byte

	// Redis addr
	RAddr string
	// MongoDB URI
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
L:
	for {
		sub := s.rC.Subscribe(ctx, "karaoke/round")
		ch := sub.Channel()
		select {
		case roundMsg := <-ch:
			round, err := strconv.Atoi(roundMsg.Payload)
			if err != nil {
				panic(err)
			}

			if round <= s.round {
				continue L
			}

			err = s.runRound(ctx)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
		err := sub.Close()
		if err != nil {
			panic(err)
		}
	}
}

func (s *Server) runRound(ctx context.Context) error {
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

	var os []*Onion
	for _, b := range plainBs {
		var o *Onion
		err := proto.Unmarshal(b, o)
		if err != nil {
			continue
		}

		if o.NextHop != "" && o.DeadDrop != "" {
			continue
		}
		if o.NextHop != "" {
			_, ok := s.pCs[o.NextHop]
			if !ok {
				continue
			}
		}

		os = append(os, o)
	}

	// Avoid empty Bloom
	if len(os) == 0 {
		return nil
	}

	// TODO: Gen noise or check Bloom

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

			err = cc.Send(&OnionMsg{
				Msg: &OnionMsg_Meta{
					Meta: &OnionMsgMeta{
						Id: s.id,
					},
				},
			})
			if err != nil {
				return err
			}
			for _, b := range oM[i] {
				err = cc.Send(&OnionMsg{
					Msg: &OnionMsg_Body{
						Body: b,
					},
				})
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return err
	}

	s.poolLock.Lock()
	s.poolMask = make(map[string]bool)
	s.poolLock.Unlock()
	return nil
}

var ErrBloomCheckFailure = errors.New("bloom check failure, which means noise loss and adversary possibility")

func (s *Server) FwdOnions(ss RPC_FwdOnionsServer) error {
	metaMsg, err := ss.Recv()
	if err != nil {
		return err
	}
	meta := metaMsg.Msg.(*OnionMsg_Meta).Meta

	var bs [][]byte
	for {
		msg, err := ss.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		bs = append(bs, msg.Msg.(*OnionMsg_Body).Body)
	}

	s.poolLock.Lock()
	if s.poolMask[meta.Id] {
		return nil
	}
	s.poolMask[meta.Id] = true

	ok := true
	for _, id := range s.layerIdx[s.c.Layer-1] {
		if !s.poolMask[id] {
			ok = false
			break
		}
	}

	s.pool = append(s.pool, bs...)
	s.poolLock.Unlock()
	if ok {
		s.poolFullCh <- true
	}
	return nil
}
