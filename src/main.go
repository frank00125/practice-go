package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RegistrationInfo struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type User struct {
	id    string
	email string
}

func getDatabaseConnection(databaseName string) (context.Context, *mongo.Database, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@localhost:27017/?directConnection=true"))
	if err != nil {
		panic(err)
	}
	userDatabase := client.Database(databaseName)

	return ctx, userDatabase, cancel
}

func passwordHashing(password string) string {
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}

	return string(passwordHashed)
}

func queryUserByEmail(email string, database *mongo.Database, ctx context.Context) *User {
	userCollection := database.Collection("user")
	var userDocumentRaw primitive.D
	err := userCollection.FindOne(ctx, bson.D{
		{"email", email},
	}).Decode(&userDocumentRaw)
	if err != nil {
		return nil
	}

	userDocument := userDocumentRaw.Map()

	return &User{
		id:    userDocument["_id"].(primitive.ObjectID).Hex(),
		email: userDocument["email"].(string),
	}
}

func insertUser(registrationInfo RegistrationInfo, database *mongo.Database, ctx context.Context) User {
	passwordHashed := passwordHashing(registrationInfo.Password)

	userCollection := database.Collection("user")
	insertResult, err := userCollection.InsertOne(ctx, bson.D{
		{"email", registrationInfo.Email},
		{"password", passwordHashed},
	})
	if err != nil {
		panic(err)
	}

	return User{
		id:    insertResult.InsertedID.(primitive.ObjectID).Hex(),
		email: registrationInfo.Email,
	}
}

func registerHandler(c *gin.Context) {
	var registrationInfo RegistrationInfo
	if err := c.ShouldBindJSON(&registrationInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get database connection
	ctx, practiceGoDatabase, cancel := getDatabaseConnection("practiceGoDatabase")
	defer cancel()

	userInCollection := queryUserByEmail(registrationInfo.Email, practiceGoDatabase, ctx)
	if userInCollection != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User exists"})
		return
	}

	newUser := insertUser(registrationInfo, practiceGoDatabase, ctx)

	c.JSON(http.StatusOK, gin.H{
		"userId": newUser.id,
	})
}

func main() {
	server := gin.Default()
	server.POST("/register", registerHandler)
	server.Run(":8000")
}
