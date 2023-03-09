package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"practice-go/repository"
)

type RegistrationInfo struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func RegisterHandler(c *gin.Context) {
	var registrationInfo RegistrationInfo
	if err := c.ShouldBindJSON(&registrationInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get database connection
	ctx, practiceGoDatabase, cancel := repository.GetDatabaseConnection("practiceGoDatabase")
	defer cancel()

	// query user from database
	userInCollection := repository.QueryUserByEmail(registrationInfo.Email, practiceGoDatabase, ctx)
	if userInCollection != nil {
		// user exists in database => 400
		c.JSON(http.StatusBadRequest, gin.H{"error": "User exists"})
		return
	}

	var userDocument repository.UserInsertion = repository.UserInsertion{
		Email:    registrationInfo.Email,
		Password: registrationInfo.Password,
	}
	newUser := repository.InsertUser(userDocument, practiceGoDatabase, ctx)

	c.JSON(http.StatusOK, gin.H{
		"userId": newUser.Id,
	})
}

func LoginHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, "Not Found")
}
