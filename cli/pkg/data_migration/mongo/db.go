package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Mdb *mongo.Client

func Setup(mongouri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error

	Mdb, err = mongo.Connect(ctx, options.Client().ApplyURI(mongouri))
	if err != nil {
		return err
	}

	// ping
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = Mdb.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	return nil
}

// get data from mongo
func GetAllData(result *[]bson.D) error {
	ctx := context.Background()

	log.Print("getting data from mongo...")

	// get all data from collection with 1000 limit each time
	const LIMIT = 10000
	for i := 0; i < LIMIT; i++ {
		cur, err := Mdb.Database("rss3").Collection("files").Find(ctx, bson.D{}, options.Find().SetLimit(LIMIT).SetSkip(int64(i)*LIMIT))
		if err != nil {
			return err
		}

		var res []bson.D

		err = cur.All(ctx, &res)
		if err != nil {
			return err
		}

		*result = append(*result, res...)

		log.Print("got data:", len(*result))

		if length := len(res); length < 1000 {
			break
		}
	}

	return nil
}
