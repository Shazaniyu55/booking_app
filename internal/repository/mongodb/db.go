package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func Connect(ctx context.Context, uri, dbName string) (*mongo.Client, *mongo.Database, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, fmt.Errorf("mongo connect: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		return nil, nil, fmt.Errorf("mongo ping: %w", err)
	}
	return client, client.Database(dbName), nil
}

func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection(bookingsColl).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "resource_id", Value: 1},
			{Key: "start_time", Value: 1},
			{Key: "end_time", Value: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("create booking index: %w", err)
	}
	return nil
}
