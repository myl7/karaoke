// Copyright (C) 2023 myl7
// SPDX-License-Identifier: Apache-2.0

package coordinator

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (co *Coordinator) Bootstrap(ctx context.Context) error {
	co.rC = redis.NewClient(&redis.Options{
		Addr: co.c.RAddr,
	})
	err := co.rC.Ping(ctx).Err()
	if err != nil {
		return err
	}

	mC, err := mongo.Connect(ctx, options.Client().ApplyURI(co.c.MURI))
	if err != nil {
		return err
	}
	err = mC.Ping(ctx, nil)
	if err != nil {
		return err
	}
	co.mDB = mC.Database("karaoke")

	sub := co.rC.Subscribe(ctx, "karaoke/bootstrap_pconfig_n_ok")
	ch := sub.Channel()
	m := make(map[string]bool)
L:
	for {
		select {
		case addr := <-ch:
			m[addr.Payload] = true
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

	coll := co.mDB.Collection("bootstrap_pconfig")
	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return err
	}

	var cs []bson.M
	err = cur.All(ctx, &cs)
	if err != nil {
		return err
	}

	// Set server IDs as continuous 0-started int
	for i, c := range cs {
		_, err = coll.UpdateByID(ctx, c["_id"], bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "id", Value: fmt.Sprintf("%d", i)},
			}},
		})
		if err != nil {
			return err
		}
	}

	err = co.rC.Publish(ctx, "karaoke/bootstrap_pconfig_ok", "1").Err()
	if err != nil {
		return err
	}

	return nil
}

func (co *Coordinator) Close(ctx context.Context) error {
	err := co.rC.Close()
	if err != nil {
		return err
	}

	err = co.mDB.Client().Disconnect(ctx)
	if err != nil {
		return err
	}

	return nil
}
