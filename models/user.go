package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"practice-go/utils"
)

type User struct {
	Id       string
	Email    string
	Password string
}

type UserInsertion struct {
	Email    string
	Password string
}

func getUserCollection(database *mongo.Database) *mongo.Collection {
	return database.Collection("user")
}

func QueryUserByEmail(email string, database *mongo.Database, ctx context.Context) *User {
	userCollection := getUserCollection(database)

	var userDocumentRaw primitive.D
	err := userCollection.FindOne(ctx, bson.D{
		{"email", email},
	}).Decode(&userDocumentRaw)
	if err != nil {
		return nil
	}

	userDocument := userDocumentRaw.Map()

	return &User{
		Id:       userDocument["_id"].(primitive.ObjectID).Hex(),
		Email:    userDocument["email"].(string),
		Password: userDocument["password"].(string),
	}
}

func InsertUser(registrationInfo UserInsertion, database *mongo.Database, ctx context.Context) User {
	userCollection := getUserCollection(database)

	insertResult, err := userCollection.InsertOne(ctx, bson.D{
		{"email", registrationInfo.Email},
		{"password", utils.PasswordHashing(registrationInfo.Password)},
	})
	if err != nil {
		panic(err)
	}

	return User{
		Id:    insertResult.InsertedID.(primitive.ObjectID).Hex(),
		Email: registrationInfo.Email,
	}
}
