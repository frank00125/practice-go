package utils

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
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

func GetRefreshTokenKey(userId string) string {
	return userId + "-refresh-token"
}

func GenerateRefreshToken(userId string, ctx context.Context, redisClient *redis.Client) (string, time.Time) {
	token := make([]byte, 32)
	_, error := rand.Read(token)
	if error != nil {
		panic(error)
	}

	hmacFunc := hmac.New(sha256.New, []byte(userId))
	hmacFunc.Write([]byte(userId))

	refreshToken := base64.URLEncoding.EncodeToString(hmacFunc.Sum(nil))

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
