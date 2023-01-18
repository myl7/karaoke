package coordinator

import (
	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Coordinator struct {
	c CoordinatorConfig

	rC  *redis.Client
	mDB *mongo.Database
}

type CoordinatorConfig struct {
	ServerN int

	// Redis addr
	RAddr string
	// MongoDB URI
	MURI string
}

func NewCoordinator(c CoordinatorConfig) *Coordinator {
	return &Coordinator{c: c}
}
