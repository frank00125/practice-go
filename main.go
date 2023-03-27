package main

import (
	"practice-go/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	server := gin.Default()
	server.POST("/register", handlers.RegisterHandler)
	server.POST("/login", handlers.LoginHandler)
	server.POST("/refresh-token", handlers.RefreshTokenHandler)
	server.Run(":8000")
}
