package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"crypto/md5"
)

type RegistrationInfo struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func registerHandler(c *gin.Context) {
	var registrationInfo RegistrationInfo
	if err := c.ShouldBindJSON(&registrationInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passwordHashed := md5.Sum([]byte(registrationInfo.Password))

	c.JSON(http.StatusOK, gin.H{
		"email":    registrationInfo.Email,
		"password": fmt.Sprintf("%x", passwordHashed),
	})
}

func main() {
	server := gin.Default()
	server.POST("/register", registerHandler)
	server.Run(":8000")
}
