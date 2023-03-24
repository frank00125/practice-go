package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func GenerateJWT(userId string) string {
	jwtSecret := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Unix() + 2*60*60,
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		panic(err)
	}

	return tokenString
}

func GenerateRefreshToken(userId string, ctx context.Context, redisClient *redis.Client) (string, time.Time) {
	token := make([]byte, 32)
	_, error := rand.Read(token)
	if error != nil {
		panic(error)
	}

	refreshToken := base64.RawURLEncoding.EncodeToString(token)
	durationToExpire := time.Duration(24 * 30 * int(time.Hour))
	userRefreshTokenKey := userId + "-refresh-token"
	fmt.Println(userId + "-refresh-token")
	error = redisClient.Set(ctx, userRefreshTokenKey, refreshToken, durationToExpire).Err()
	if error != nil {
		panic(error)
	}

	expirationDate := time.Now().Add(durationToExpire)

	return refreshToken, expirationDate
}
