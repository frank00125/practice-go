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

var jwtSecret string = os.Getenv("JWT_SECRET")

func GenerateJWT(userId string) string {
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

	hmacFunc := hmac.New(sha256.New, []byte(userId))
	hmacFunc.Write([]byte(userId))

	refreshToken := base64.URLEncoding.EncodeToString(hmacFunc.Sum(nil))

	durationToExpire := time.Duration(24 * 30 * int(time.Hour))
	error = redisClient.Set(ctx, refreshToken, userId, durationToExpire).Err()
	if error != nil {
		panic(error)
	}

	expirationDate := time.Now().Add(durationToExpire)

	return refreshToken, expirationDate
}

func GetJwtPayload(tokenStr string) (jwt.MapClaims, error) {
	token, error := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if error != nil {
		return nil, error
	}

	return token.Claims.(jwt.MapClaims), nil
}

func VerifyJwt(tokenStr string) (jwt.MapClaims, error) {
	token, error := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(jwtSecret), nil
	})
	if error != nil {
		return nil, error
	}
	return token.Claims.(jwt.MapClaims), nil
}
