package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Public config
type pConfig struct {
	addr string
	pk   *[32]byte
}

func (s *Server) Bootstrap(ctx context.Context) error {
	s.rC = redis.NewClient(&redis.Options{
		Addr: s.c.RAddr,
	})
	err := s.rC.Ping(ctx).Err()
	if err != nil {
		return err
	}

	mC, err := mongo.NewClient(options.Client().ApplyURI(s.c.MURI))
	if err != nil {
		return err
	}
	s.mDB = mC.Database("karaoke")

	ip, err := publicIP()
	if err != nil {
		return err
	}
	s.addr = fmt.Sprintf("%s:%d", ip, s.c.Port)

	coll := s.mDB.Collection("bootstrap-pconfig")
	coll.InsertOne(ctx, map[string]any{
		"addr": s.addr,
		"pk":   s.c.PK,
	})

	sub := s.rC.Subscribe(ctx, "karaoke/bootstrap-pconfig-ok")
	ch := sub.Channel()
	select {
	case ok := <-ch:
		if ok.Payload != "1" {
			return ErrInvalidBootstrapPConfigOKSignal
		}
	case <-ctx.Done():
		return ctx.Err()
	}
	sub.Close()

	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return err
	}

	var cs []bson.M
	err = cur.All(ctx, &cs)
	if err != nil {
		return err
	}

	for _, c := range cs {
		addr := c["addr"].(string)
		id := c["id"].(string)
		if addr == s.addr {
			s.id = id
		} else {
			s.pCs[id] = pConfig{
				addr: addr,
				pk:   c["pk"].(*[32]byte),
			}
		}
	}
	return nil
}

var ErrInvalidBootstrapPConfigOKSignal = errors.New("invalid bootstrap-pconfig-ok signal: not 1")

// From Cloudflare API
func publicIP() (string, error) {
	res, err := http.Get("https://cloudflare.com/cdn-cgi/trace")
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
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
