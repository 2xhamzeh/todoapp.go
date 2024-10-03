package services

import (
	"ToDo/api/auth"
	"ToDo/api/db"
	"ToDo/api/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(ctx context.Context, u models.AuthDTO) error {
	usersCollection := db.GetCollection("users")

	// hash password
	hashedPassword, err := auth.HashPassword(u.Password)
	if err != nil {
		return err
	}

	user := models.User{
		ID:         primitive.NewObjectID(),
		Username:   u.Username,
		Password:   hashedPassword,
		Todos:      []primitive.ObjectID{},
		Categories: []models.Category{},
	}

	_, err = usersCollection.InsertOne(ctx, user)
	return err
}

func UsernameExists(ctx context.Context, username string) error {
	usersCollection := db.GetCollection("users")
	err := usersCollection.FindOne(ctx, bson.M{"username": username}).Err()
	return err
}

func LoginUser(ctx context.Context, u models.AuthDTO) (*string, error) {
	usersCollection := db.GetCollection("users")

	// check if user exists
	user := models.User{}
	err := usersCollection.FindOne(ctx, bson.M{"username": u.Username}).Decode(&user)
	if err != nil {
		return nil, err
	}

	// check password
	isAuth := auth.CheckPasswordHash(u.Password, user.Password)
	if !isAuth {
		return nil, errors.New("invalid credentials")
	}

	// generate token
	token, err := auth.GenerateJWT(user.ID.Hex())
	if err != nil {
		return nil, err
	}
	return &token, nil
}
