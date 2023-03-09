package main

import (
	"github.com/gin-gonic/gin"

	"practice-go/handlers"
)

func main() {
	server := gin.Default()
	server.POST("/register", handlers.RegisterHandler)
	server.Run(":8000")
}
