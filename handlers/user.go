package handlers

import (
	"net/http"
	"strings"

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

type RefreshTokenInfo struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func RegisterHandler(c *gin.Context) {
	var registrationInfo RegistrationInfo
	if err := c.ShouldBindJSON(&registrationInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get database connection
	dbCtx, practiceGoDatabase, cancel := models.GetDatabaseConnection("practiceGoDatabase")
	defer cancel()

	// get redis connection
	redisCtx, redisClient := models.GetRedisConnection()
	defer redisClient.Conn().Close()

	// query user from database
	userInCollection := models.QueryUserByEmail(registrationInfo.Email, practiceGoDatabase, dbCtx)
	if userInCollection != nil {
		// user exists in database => 400
		c.JSON(http.StatusBadRequest, gin.H{"error": "User exists"})
		return
	}

	var userDocument models.UserInsertion = models.UserInsertion{
		Email:    registrationInfo.Email,
		Password: registrationInfo.Password,
	}
	newUser := models.InsertUser(userDocument, practiceGoDatabase, dbCtx)

	// generate jwt
	token := utils.GenerateJWT(newUser.Id)

	// generate refresh token
	refreshToken, expiredAt := utils.GenerateRefreshToken(newUser.Id, redisCtx, redisClient)

	c.JSON(http.StatusOK, gin.H{
		"tokenType":    "Bearer",
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
	dbCtx, practiceGoDatabase, cancel := models.GetDatabaseConnection("practiceGoDatabase")
	defer cancel()

	// get redis connection
	redisCtx, redisClient := models.GetRedisConnection()
	defer redisClient.Conn().Close()

	// Get user in database
	userInCollection := models.QueryUserByEmail(loginInfo.Email, practiceGoDatabase, dbCtx)
	// verify password
	isPasswordVerified := utils.PasswordVerify(userInCollection.Password, loginInfo.Password)

	if !isPasswordVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// generate jwt
	token := utils.GenerateJWT(userInCollection.Id)

	// generate refresh token
	refreshToken, expiredAt := utils.GenerateRefreshToken(userInCollection.Id, redisCtx, redisClient)

	c.JSON(http.StatusOK, gin.H{
		"tokenType":    "Bearer",
		"token":        token,
		"refreshToken": refreshToken,
		"expiredAt":    expiredAt,
	})
}

func RefreshTokenHandler(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	redisCtx, redisClient := models.GetRedisConnection()

	var refreshTokenInfo RefreshTokenInfo
	if err := c.ShouldBindJSON(&refreshTokenInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the type of the authorization header is a bearer token
	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden",
		})
		return
	}

	// Get userId from Authorization header
	accessToken := strings.ReplaceAll(authorizationHeader, "Bearer ", "")
	jwtPayload, error := utils.GetJwtPayload(accessToken)
	if error != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden",
		})
		return
	}

	userIdInPayload := jwtPayload["userId"].(string)

	// Get value for refresh token from redis
	userIdInRedis, _, error := models.GetRefreshTokenFromInMemoryDB(redisCtx, redisClient, refreshTokenInfo.RefreshToken)

	if error != nil || userIdInPayload != userIdInRedis {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden",
		})
		return
	}

	// Delete old refresh token from in-memory database
	models.RemoveRefreshTokenFromInMemoryDB(redisCtx, redisClient, refreshTokenInfo.RefreshToken)

	// Generate new jwt
	accessToken = utils.GenerateJWT(userIdInRedis)

	// Generate new refresh token
	refreshToken, expiredAt := utils.GenerateRefreshToken(userIdInRedis, redisCtx, redisClient)

	c.JSON(http.StatusOK, gin.H{
		"tokenType":    "Bearer",
		"token":        accessToken,
		"refreshToken": refreshToken,
		"expiredAt":    expiredAt,
	})
}
