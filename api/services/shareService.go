package services

import (
	"ToDo/api/db"
	"ToDo/api/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"slices"
)

func ShareWithUser(ctx context.Context, usernameToShareWith string, userID primitive.ObjectID) error {
	usersCollection := db.GetCollection("users")

	user := models.User{}
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return err
	}

	filter := bson.M{"username": usernameToShareWith}
	update := bson.M{"$push": bson.M{"sharedWithMe": user.Username}}

	_, err = usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func UnshareWithUser(ctx context.Context, usernameToUnshareWith string, userID primitive.ObjectID) error {
	usersCollection := db.GetCollection("users")

	user := models.User{}
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return err
	}

	filter := bson.M{"username": usernameToUnshareWith}
	update := bson.M{"$pull": bson.M{"sharedWithMe": user.Username}}

	_, err = usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// checks if the username shared their todos with the authetnicated user
func CheckIfUsernameSharedWithAuthUser(ctx context.Context, username string, userID primitive.ObjectID) (bool, error) {
	usersCollection := db.GetCollection("users")
	user := models.User{}
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return false, err
	}

	if slices.Contains(user.SharedWithMe, username) {
		return true, nil
	}
	return false, nil
}

// checks if the authenticated user shared their todos with the username
func CheckIfAuthUserSharedWithUsername(ctx context.Context, username string, userID primitive.ObjectID) (bool, error) {
	usersCollection := db.GetCollection("users")
	user := models.User{}
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return false, err
	}

	userToShareWith := models.User{}
	err = usersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&userToShareWith)
	if err != nil {
		return false, err
	}

	if slices.Contains(userToShareWith.SharedWithMe, user.Username) {
		return true, nil
	}
	return false, nil
}

func GetUsernamesSharedWithAuth(ctx context.Context, userID primitive.ObjectID) ([]string, error) {
	usersCollection := db.GetCollection("users")
	user := models.User{}
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user.SharedWithMe, nil
}
