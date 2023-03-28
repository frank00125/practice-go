package main

import (
	"practice-go/handlers"
	"practice-go/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	go authServer()

	articleServer()
}

func authServer() {
	authServer := gin.Default()

	authServer.POST("/register", handlers.RegisterHandler)
	authServer.POST("/login", handlers.LoginHandler)
	authServer.POST("/refresh-token", handlers.RefreshTokenHandler)

	authServer.Run(":8000")
}

func articleServer() {
	articleServer := gin.Default()
	articleServer.Use(middleware.VerifyToken())

	articleServer.POST("/article/create", handlers.CreateArticle)

	articleServer.Run(":8001")
}
