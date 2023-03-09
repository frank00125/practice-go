package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"practice-go/utils"
)

type User struct {
	Id    string
	Email string
}

type UserInsertion struct {
	Email    string
	Password string
}

func QueryUserByEmail(email string, database *mongo.Database, ctx context.Context) *User {
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
		Id:    userDocument["_id"].(primitive.ObjectID).Hex(),
		Email: userDocument["email"].(string),
	}
}

func InsertUser(registrationInfo UserInsertion, database *mongo.Database, ctx context.Context) User {
	passwordHashed := utils.PasswordHashing(registrationInfo.Password)

	userCollection := database.Collection("user")
	insertResult, err := userCollection.InsertOne(ctx, bson.D{
		{"email", registrationInfo.Email},
		{"password", passwordHashed},
	})
	if err != nil {
		panic(err)
	}

	return User{
		Id:    insertResult.InsertedID.(primitive.ObjectID).Hex(),
		Email: registrationInfo.Email,
	}
}
