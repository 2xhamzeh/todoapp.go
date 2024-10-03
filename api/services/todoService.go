package services

import (
	"ToDo/api/db"
	"ToDo/api/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"slices"
)

func GetTodos(ctx context.Context, userID primitive.ObjectID) (*[]models.ToDo, error) {
	// array of todo ids
	todoArray, err := getUserTodosIDArray(ctx, userID)
	if err != nil {
		return nil, err
	}

	result, err := GetTodosFromIDArray(ctx, todoArray)
	if err != nil {
		return nil, err
	}

	sortedResult := sortTodosBasedOnIDArray(ctx, todoArray, result)

	return sortedResult, nil
}

func CreateTodo(ctx context.Context, userID primitive.ObjectID, t models.CreateTodoDTO) (*models.ToDo, error) {
	usersCollection := db.GetCollection("users")
	todosCollection := db.GetCollection("todos")

	// construct the item we want to insert
	todo := models.ToDo{
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

func UpdateTodo(ctx context.Context, todoID primitive.ObjectID, update models.UpdateTodoDTO) (*models.ToDo, error) {
	todosCollection := db.GetCollection("todos")
	// update the document
	result := models.ToDo{}
	err := todosCollection.FindOneAndUpdate(ctx, bson.M{"_id": todoID}, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func DeleteTodo(ctx context.Context, userID primitive.ObjectID, todoID primitive.ObjectID) (*models.ToDo, error) {
	todosCollection := db.GetCollection("todos")
	usersCollection := db.GetCollection("users")

	// remove The id of the item from the user
	_, err := usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$pull": bson.M{"todos": todoID}})
	if err != nil {
		return nil, err
	}

	// delete the document
	result := models.ToDo{}
	err = todosCollection.FindOneAndDelete(ctx, bson.M{"_id": todoID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func TodoBelongsToUser(ctx context.Context, userID primitive.ObjectID, todoID primitive.ObjectID) (bool, error) {
	usersCollection := db.GetCollection("users")

	u := models.User{}
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

func SortTodos(ctx context.Context, userID primitive.ObjectID, newOrder []primitive.ObjectID) (bool, error) {
	//fmt.Println(newOrder)
	currentOrder, err := getUserTodosIDArray(ctx, userID)
	if err != nil {
		return false, err
	}
	//fmt.Println(currentOrder)
	if !idSliceHasSameContent(currentOrder, newOrder) {
		return false, nil
	}

	usersCollection := db.GetCollection("users")
	_, err = usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"todos": newOrder}})
	if err != nil {
		return false, err
	}
	return true, nil
}

// this method checks if two slices have the same elements, it ignores the order
func idSliceHasSameContent(x []primitive.ObjectID, y []primitive.ObjectID) bool {
	if len(x) != len(y) {
		return false
	}

	diff := map[primitive.ObjectID]int{}
	for _, id := range x {
		diff[id]++
	}
	for _, id := range y {
		diff[id]--
	}
	for _, v := range diff {
		if v != 0 {
			return false
		}
	}
	return true
}

func getUserTodosIDArray(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	usersCollection := db.GetCollection("users")
	// get the user
	u := models.User{}
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return u.Todos, nil
}

func GetTodosFromIDArray(ctx context.Context, todoIDArr []primitive.ObjectID) ([]models.ToDo, error) {
	todosCollection := db.GetCollection("todos")

	// get the todos of the user
	cursor, err := todosCollection.Find(ctx, bson.M{"_id": bson.M{"$in": todoIDArr}}) // query all the todos with an id in the todos array in the user
	if err != nil {
		return nil, err
	}
	// we use this instead of var result [], to make sure we get an empty slice and not nil if the user has no todos
	result := make([]models.ToDo, 0)
	for cursor.Next(ctx) {
		todo := models.ToDo{}
		err = cursor.Decode(&todo) // we use a pointer (&) to make sure we decode to our variable and not a copy
		if err != nil {
			return nil, err
		}
		result = append(result, todo)
	}
	return result, nil
}

func sortTodosBasedOnIDArray(ctx context.Context, todoIDArr []primitive.ObjectID, Todos []models.ToDo) *[]models.ToDo {
	// sort the results to make sure the order of the returned results matches the order of the ids in the users todos array
	// an empty slice for the sorted results
	sortedResult := make([]models.ToDo, 0)

	// creating a map of our todos for easy lookup
	mapOfTodos := map[primitive.ObjectID]models.ToDo{}
	for _, todo := range Todos {
		mapOfTodos[todo.ID] = todo
	}

	// going through the array of the user
	for _, id := range todoIDArr {
		// checking for each item to make sure it is contained in the map
		todo, ok := mapOfTodos[id]
		if ok {
			sortedResult = append(sortedResult, todo)
		}
	}

	return &sortedResult
}
