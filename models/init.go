package models

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func GetDatabaseConnection(databaseName string) (context.Context, *mongo.Database, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@localhost:27017/?directConnection=true"))
	if err != nil {
		panic(err)
	}
	userDatabase := client.Database(databaseName)

	return ctx, userDatabase, cancel
}

func GetRedisConnection() (context.Context, *redis.Client) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})
	// test connection for redis
	_, error := redisClient.Ping(ctx).Result()
	if error != nil {
		panic(error)
	}

	return ctx, redisClient
}
