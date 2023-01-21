// Copyright (C) 2023 myl7
// SPDX-License-Identifier: Apache-2.0

package coordinator

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/myl7/karaoke/pkg/rpc"
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
		rS := rpc.RoundStartMsg{Round: co.round}
		rSB, err := json.Marshal(rS)
		if err != nil {
			panic(err)
		}
		err = co.rC.Publish(ctx, "karaoke/round", rSB).Err()
		if err != nil {
			return err
		}

		sub := co.rC.Subscribe(ctx, "karaoke/round_ok")
		ch := sub.Channel()
		m := make(map[string]bool, co.c.ServerN)
	L:
		for {
			select {
			case rEB := <-ch:
				var rE rpc.RoundEndMsg
				err = json.Unmarshal([]byte(rEB.Payload), &rE)
				if err != nil {
					panic(err)
				}
				round := rE.Round
				id := rE.ID

				if round != co.round {
					panic(ErrRoundNotMatch)
				}
				m[id] = true
				if len(m) >= co.c.ServerN {
					break L
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		err = sub.Close()
		if err != nil {
			panic(err)
		}
	}
}

var ErrRoundNotMatch = errors.New("start and end round not match")
