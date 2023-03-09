package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetDatabaseConnection(databaseName string) (context.Context, *mongo.Database, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@localhost:27017/?directConnection=true"))
	if err != nil {
		panic(err)
	}
	userDatabase := client.Database(databaseName)

	return ctx, userDatabase, cancel
}
