package middleware

import (
	"fmt"
	"net/http"
	"practice-go/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")

		// Check if the type of the authorization header is a bearer token
		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
			})
			return
		}

		// Get userId from Authorization header
		accessToken := strings.ReplaceAll(authorizationHeader, "Bearer ", "")

		// Verify Token
		payload, error := utils.VerifyJwt(accessToken)
		if error != nil {
			fmt.Println(error)
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
			})
			return
		}

		c.Set("verifiedUserId", payload["userId"])
	}
}
