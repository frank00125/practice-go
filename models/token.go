package models

import (
	"context"
	"fmt"
	"practice-go/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetRefreshTokenFromInMemoryDB(ctx context.Context, redisClient *redis.Client, userId string) (string, time.Time, error) {
	refreshTokenKey := utils.GetRefreshTokenKey(userId)
	refreshTokenInRedis, error := redisClient.Get(ctx, refreshTokenKey).Result()
	now := time.Now().UTC()
	if error != nil {
		return "", now, error
	}

	ttl, error := redisClient.TTL(ctx, refreshTokenKey).Result()
	if error != nil {
		return "", now, error
	}

	fmt.Println(now, ttl)

	return refreshTokenInRedis, now.Add(ttl), nil
}
