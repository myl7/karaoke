package coordinator

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Coordinator struct {
	c CoordinatorConfig

	round int

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

func (co *Coordinator) Run(ctx context.Context) error {
	for {
		co.round += 1
		err := co.rC.Publish(ctx, "karaoke/round", strconv.Itoa(co.round)).Err()
		if err != nil {
			return err
		}

		sub := co.rC.Subscribe(ctx, "karaoke/round_ok")
		ch := sub.Channel()
		select {
		case roundMsg := <-ch:
			round, err := strconv.Atoi(roundMsg.Payload)
			if err != nil {
				panic(err)
			}

			if round != co.round {
				panic(ErrRoundNotMatch)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
		err = sub.Close()
		if err != nil {
			panic(err)
		}
	}
}

var ErrRoundNotMatch = errors.New("start and end round not match")
