package coordinator

import (
	"context"
	"errors"
	"strconv"
	"time"

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
	// To make sure all servers are waiting rounds but still separate Run & Bootstrap,
	// poll subscriber number until enough.
	for {
		nM, err := co.rC.PubSubNumSub(ctx, "karaoke/round").Result()
		if err != nil {
			return err
		}
		if nM["karaoke/round"] >= int64(co.c.ServerN) {
			break
		}
		select {
		case <-time.After(1 * time.Second):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

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
