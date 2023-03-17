package main

import (
	"github.com/gin-gonic/gin"

	"practice-go/handlers"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	server := gin.Default()
	server.POST("/register", handlers.RegisterHandler)
	server.Run(":8000")
}
