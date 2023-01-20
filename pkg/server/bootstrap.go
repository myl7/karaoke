package server

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// PConfig is server publicly available config
type PConfig struct {
	Addr  string
	PK    *[32]byte
	Layer int
}

func (s *Server) Bootstrap(ctx context.Context) error {
	s.rC = redis.NewClient(&redis.Options{
		Addr: s.c.RAddr,
	})
	err := s.rC.Ping(ctx).Err()
	if err != nil {
		return err
	}

	mC, err := mongo.Connect(ctx, options.Client().ApplyURI(s.c.MURI))
	if err != nil {
		return err
	}
	err = mC.Ping(ctx, nil)
	if err != nil {
		return err
	}
	s.mDB = mC.Database("karaoke")

	// TODO: Dead drop external persistence
	s.deadDrop = make(map[string][]byte)

	if s.c.Addr == "" {
		ip, err := publicIP()
		if err != nil {
			return err
		}
		s.addr = ip + s.c.LAddr
	} else {
		s.addr = s.c.Addr
	}

	coll := s.mDB.Collection("bootstrap_pconfig")
	// Initially use addr as unique ID
	_, err = coll.DeleteOne(ctx, bson.D{
		{Key: "addr", Value: s.addr},
	})
	if err != nil {
		return err
	}
	_, err = coll.InsertOne(ctx, map[string]any{
		"addr":  s.addr,
		"pk":    s.c.PK,
		"layer": s.c.Layer,
	})
	if err != nil {
		return err
	}

	sub := s.rC.Subscribe(ctx, "karaoke/bootstrap_pconfig_ok")
	ch := sub.Channel()

	err = s.rC.Publish(ctx, "karaoke/bootstrap_pconfig_n_ok", s.addr).Err()
	if err != nil {
		return err
	}

	select {
	case ok := <-ch:
		if ok.Payload != "1" {
			panic(ErrBootstrapInvalidPConfigOKSignal)
		}
	case <-ctx.Done():
		return ctx.Err()
	}
	err = sub.Close()
	if err != nil {
		panic(err)
	}

	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return err
	}

	var cs []bson.M
	err = cur.All(ctx, &cs)
	if err != nil {
		return err
	}

	s.pCs = make(map[string]PConfig, len(cs))
	s.layerIdx = make(map[int][]string)
	for _, c := range cs {
		addr := c["addr"].(string)
		id := c["id"].(string)
		layer := int(c["layer"].(int32))
		if addr == s.addr {
			s.id = id
		}
		s.pCs[id] = PConfig{
			Addr:  addr,
			PK:    (*[32]byte)(c["pk"].(primitive.Binary).Data),
			Layer: layer,
		}
		s.layerIdx[layer] = append(s.layerIdx[layer], id)
	}

	zLCs := make(map[string]RPCClient)
	pLCs := make(map[string]RPCClient)
	var lock sync.Mutex
	g, _ := errgroup.WithContext(ctx)
	if s.c.Layer > 0 {
		s.poolMask = make(map[string]bool, len(s.layerIdx[s.c.Layer-1]))
		s.poolFullCh = make(chan bool)

		for _, id := range s.layerIdx[0] {
			i := id
			g.Go(func() error {
				conn, err := grpc.Dial(s.pCs[i].Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					return err
				}
				client := NewRPCClient(conn)

				lock.Lock()
				defer lock.Unlock()
				zLCs[i] = client
				return nil
			})
		}
	}
	if len(s.layerIdx[s.c.Layer+1]) > 0 {
		for _, id := range s.layerIdx[s.c.Layer+1] {
			i := id
			g.Go(func() error {
				conn, err := grpc.Dial(s.pCs[i].Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					return err
				}
				client := NewRPCClient(conn)

				lock.Lock()
				defer lock.Unlock()
				pLCs[i] = client
				return nil
			})
		}
	}
	err = g.Wait()
	if err != nil {
		return err
	}
	s.zeroLayerClients = zLCs
	s.postLayerClients = pLCs

	// TODO: Replace dummy clients
	for i := 0; i < 10; i++ {
		s.clients = append(s.clients, *NewClient(ClientConfig{
			PCs:      s.pCs,
			LayerIdx: s.layerIdx,
		}))
	}

	return nil
}

var ErrBootstrapInvalidPConfigOKSignal = errors.New("invalid bootstrap_pconfig_ok signal: not 1")

func (s *Server) Close(ctx context.Context) error {
	err := s.rC.Close()
	if err != nil {
		return err
	}

	err = s.mDB.Client().Disconnect(ctx)
	if err != nil {
		return err
	}

	return nil
}

// publicIP gets the public IP from Cloudflare API
func publicIP() (string, error) {
	res, err := http.Get("https://cloudflare.com/cdn-cgi/trace")
	if err != nil {
		return "", err
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			panic(err)
		}
	}()
	bodyB, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	body := string(bodyB)
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ip=") {
			return strings.TrimSpace(strings.TrimPrefix(line, "ip=")), nil
		}
	}
	return "", ErrPublicIPNotFound
}

var ErrPublicIPNotFound = errors.New("public IP not found in body returned by Cloudflare")

func (s *Server) Listen(ctx context.Context) error {
	l, err := net.Listen("tcp", s.c.LAddr)
	if err != nil {
		return err
	}
	gs := grpc.NewServer()
	RegisterRPCServer(gs, s)
	return gs.Serve(l)
}
