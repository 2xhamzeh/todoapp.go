package todo

import (
	"ToDo/api/config"
	"ToDo/api/user"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"slices"
)

func GetUserTodos(ctx context.Context, userID primitive.ObjectID) (*[]ToDo, error) {
	// get the collections
	usersCollection := config.GetCollection("users")
	todosCollection := config.GetCollection("todos")

	// get the user
	u := user.User{}
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&u)
	if err != nil {
		return nil, err
	}

	// get the todos of the user
	cursor, err := todosCollection.Find(ctx, bson.M{"_id": bson.M{"$in": u.Todos}}) // query all the todos with an id in the todos array in the user
	if err != nil {
		return nil, err
	}
	// we use this instead of var result [], to make sure we get an empty slice and not nil if the user has no todos
	result := make([]ToDo, 0)
	for cursor.Next(ctx) {
		todo := ToDo{}
		err := cursor.Decode(&todo) // we use a pointer (&) to make sure we decode to our variable and not a copy
		if err != nil {
			return nil, err
		}
		result = append(result, todo)
	}

	// sort the results to make sure the order of the returned results matches the order of the ids in the users todos array
	
	return &result, nil
}

func CreateTodo(ctx context.Context, userID primitive.ObjectID, t createDTO) (*ToDo, error) {
	usersCollection := config.GetCollection("users")
	todosCollection := config.GetCollection("todos")

	// construct the item we want to insert
	todo := ToDo{
		primitive.NewObjectID(),
		t.Title,
		t.Text,
		false,
	}

	// insert the item to the database
	_, err := todosCollection.InsertOne(ctx, todo)
	if err != nil {
		return nil, err
	}

	// insert the id of the item to the users collections array
	_, err = usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$push": bson.M{"todos": todo.ID}})
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func UpdateTodo(ctx context.Context, todoID primitive.ObjectID, update updateDTO) (*ToDo, error) {
	todosCollection := config.GetCollection("todos")
	// update the document
	result := ToDo{}
	err := todosCollection.FindOneAndUpdate(ctx, bson.M{"_id": todoID}, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func DeleteTodo(ctx context.Context, userID primitive.ObjectID, todoID primitive.ObjectID) (*ToDo, error) {
	todosCollection := config.GetCollection("todos")
	usersCollection := config.GetCollection("users")

	// remove The id of the item from the user
	_, err := usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$pull": bson.M{"todos": todoID}})
	if err != nil {
		return nil, err
	}

	// delete the document
	result := ToDo{}
	err = todosCollection.FindOneAndDelete(ctx, bson.M{"_id": todoID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func BelongsToUser(ctx context.Context, userID primitive.ObjectID, todoID primitive.ObjectID) (bool, error) {
	usersCollection := config.GetCollection("users")

	u := user.User{}
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&u)
	if err != nil {
		return false, err
	}

	isOwnedByUser := slices.Contains(u.Todos, todoID)
	if !isOwnedByUser {
		return false, nil
	}
	return true, nil
}
