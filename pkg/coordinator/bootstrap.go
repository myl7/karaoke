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

	mC, err := mongo.NewClient(options.Client().ApplyURI(co.c.MURI))
	if err != nil {
		return err
	}
	co.mDB = mC.Database("karaoke")

	coll := co.mDB.Collection("bootstrap_pconfig")
	cStream, err := coll.Watch(ctx, mongo.Pipeline{bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "operationType", Value: "insert"},
		}},
	}})
	if err != nil {
		return err
	}

	n := co.c.ServerN
	for {
		if !cStream.Next(ctx) {
			return cStream.Err()
		}

		n -= 1
		if n == 0 {
			break
		}
	}
	cStream.Close(ctx)

	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return err
	}

	var cs []bson.M
	err = cur.All(ctx, &cs)
	if err != nil {
		return err
	}

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
