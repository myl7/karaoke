package server

import (
	"crypto/rand"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/nacl/box"
)

type Server struct {
	c ServerConfig

	// By detection
	addr string
	// Set in bootstrap
	id  string
	pCs map[string]pConfig

	rC  *redis.Client
	mDB *mongo.Database
}

type ServerConfig struct {
	Layer int
	// Server must have a public IP, so only port is needed
	Port int
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
