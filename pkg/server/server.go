package server

import (
	"crypto/rand"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/nacl/box"
)

type Server struct {
	c ServerConfig

	// Set in bootstrap
	addr string
	id   string
	pCs  map[string]pConfig

	rC  *redis.Client
	mDB *mongo.Database
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
