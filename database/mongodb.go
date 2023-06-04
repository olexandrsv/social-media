package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MI MongoInstance

type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func InitMongoDatabase(mongoURI, name string) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI + "/" + name))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(name)

	if err != nil {
		return err
	}

	MI = MongoInstance{
		Client: client,
		DB:     db,
	}
	return nil
}
