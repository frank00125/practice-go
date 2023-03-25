package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"practice-go/models"
	"practice-go/utils"
)

type RegistrationInfo struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginInfo struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func getRefreshToken(userId string) (string, time.Time) {
	ctx, redisClient := models.GetRedisConnection()
	defer redisClient.Conn().Close()

	refreshTokenInCache, expiredAt, error := models.GetRefreshTokenFromInMemoryDB(ctx, redisClient, userId)
	if error != nil {
		return utils.GenerateRefreshToken(userId, ctx, redisClient)
	}

	fmt.Println("Successfully get refresh token from redis.")
	return refreshTokenInCache, expiredAt
}

func RegisterHandler(c *gin.Context) {
	var registrationInfo RegistrationInfo
	if err := c.ShouldBindJSON(&registrationInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get database connection
	ctx, practiceGoDatabase, cancel := models.GetDatabaseConnection("practiceGoDatabase")
	defer cancel()

	// query user from database
	userInCollection := models.QueryUserByEmail(registrationInfo.Email, practiceGoDatabase, ctx)
	if userInCollection != nil {
		// user exists in database => 400
		c.JSON(http.StatusBadRequest, gin.H{"error": "User exists"})
		return
	}

	var userDocument models.UserInsertion = models.UserInsertion{
		Email:    registrationInfo.Email,
		Password: registrationInfo.Password,
	}
	newUser := models.InsertUser(userDocument, practiceGoDatabase, ctx)

	// generate jwt
	token := utils.GenerateJWT(newUser.Id)

	// generate refresh token
	refreshToken, expiredAt := getRefreshToken(newUser.Id)

	c.JSON(http.StatusOK, gin.H{
		"token":        token,
		"refreshToken": refreshToken,
		"expiredAt":    expiredAt,
	})
}

func LoginHandler(c *gin.Context) {
	var loginInfo LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get database connection
	ctx, practiceGoDatabase, cancel := models.GetDatabaseConnection("practiceGoDatabase")
	defer cancel()

	userInCollection := models.QueryUserByEmail(loginInfo.Email, practiceGoDatabase, ctx)
	isPasswordVerified := utils.PasswordVerify(userInCollection.Password, loginInfo.Password)

	if !isPasswordVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// generate jwt
	token := utils.GenerateJWT(userInCollection.Id)

	// generate refresh token
	refreshToken, expiredAt := getRefreshToken(userInCollection.Id)

	c.JSON(http.StatusOK, gin.H{
		"token":        token,
		"refreshToken": refreshToken,
		"expiredAt":    expiredAt,
	})
}
