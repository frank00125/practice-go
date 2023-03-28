package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Article struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func CreateArticle(c *gin.Context) {
	verifiedUserId := c.GetString("verifiedUserId")
	fmt.Println("verifiedUserId", verifiedUserId)

	var newArticle Article
	if err := c.ShouldBindJSON(&newArticle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Internal Server Error",
	})
}
