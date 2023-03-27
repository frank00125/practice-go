package models

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetRefreshTokenFromInMemoryDB(ctx context.Context, redisClient *redis.Client, refreshToken string) (string, time.Time, error) {
	value, error := redisClient.Get(ctx, refreshToken).Result()
	if error != nil {
		return "", time.Now(), error
	}

	ttl, error := redisClient.TTL(ctx, refreshToken).Result()
	if error != nil {
		return "", time.Now(), error
	}

	return value, time.Now().Add(time.Duration(ttl)), nil
}

func RemoveRefreshTokenFromInMemoryDB(ctx context.Context, redisClient *redis.Client, refreshToken string) error {
	_, error := redisClient.Del(ctx, refreshToken).Result()

	return error
}
