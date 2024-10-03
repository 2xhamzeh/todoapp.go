package services

import (
	"ToDo/api/db"
	"ToDo/api/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCategory(ctx context.Context, userID primitive.ObjectID, name string) (*models.Category, error) {
	// get the users collection
	usersCollection := db.GetCollection("users")
	cat := models.Category{
		name,
		[]primitive.ObjectID{},
	}
	_, err := usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$push": bson.M{"categories": cat}})
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// CategoryExists : if err is nil, category exists.
// if err is ErrNoDocument, category doesn't exist.
// if err is different, internal error
func CategoryExists(ctx context.Context, userID primitive.ObjectID, name string) error {
	usersCollection := db.GetCollection("users")
	// check if the category already exists
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID, "categories.name": name}).Err()
	return err
}

func DeleteCategory(ctx context.Context, userID primitive.ObjectID, name string) error {
	usersCollection := db.GetCollection("users")
	_, err := usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$pull": bson.M{"categories": bson.M{"name": name}}})
	if err != nil {
		return err
	}
	return nil
}

func AddTodoToCategory(ctx context.Context, userID primitive.ObjectID, name string, todo primitive.ObjectID) error {
	usersCollection := db.GetCollection("users")

	_, err := usersCollection.UpdateOne(ctx, bson.M{"_id": userID, "categories.name": name}, bson.M{"$push": bson.M{"categories.$.todos": todo}})
	if err != nil {
		return err
	}
	return nil
}

func RemoveTodoFromCategory(ctx context.Context, userID primitive.ObjectID, name string, todo primitive.ObjectID) error {
	usersCollection := db.GetCollection("users")

	_, err := usersCollection.UpdateOne(ctx, bson.M{"_id": userID, "categories.name": name}, bson.M{"$pull": bson.M{"categories.$.todos": todo}})
	if err != nil {
		return err
	}
	return nil
}

func GetUserCategories(ctx context.Context, userID primitive.ObjectID) ([]models.Category, error) {
	usersCollection := db.GetCollection("users")

	u := models.User{}
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return u.Categories, nil
}

func GetCategory(ctx context.Context, userID primitive.ObjectID, name string) ([]primitive.ObjectID, error) {
	usersCollection := db.GetCollection("users")

	// get the user
	user := models.User{}
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	// find the category
	for _, category := range user.Categories {
		if category.Name == name {
			return category.Todos, nil
		}
	}
	return nil, errors.New("Category not found") // this error probably won't be returned, because we check that the category exists before calling this method

}
