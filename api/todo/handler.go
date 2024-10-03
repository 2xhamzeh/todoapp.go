package todo

import (
	"ToDo/api/config"
	"ToDo/api/user"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"slices"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	usersCollection := config.GetCollection("users")
	todosCollection := config.GetCollection("todos")

	userIdph := r.Context().Value("userId")
	userId, err := primitive.ObjectIDFromHex(userIdph.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u := new(user.User)
	err = usersCollection.FindOne(r.Context(), bson.M{"_id": userId}).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cursor, err := todosCollection.Find(r.Context(), bson.M{"_id": bson.M{"$in": u.Todos}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// we use this instead of var result [], to make sure we get an empty slice and not nil if the user has no todos
	result := make([]ToDo, 0)
	for cursor.Next(r.Context()) {
		todo := ToDo{}
		err := cursor.Decode(&todo) // we use a pointer (&) to make sure we decode to our variable and not a copy
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, todo)
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	// get the collections
	todosCollection := config.GetCollection("todos")
	usersCollection := config.GetCollection("users")

	// get the user id from context
	userIdph := r.Context().Value("userId")
	userId, err := primitive.ObjectIDFromHex(userIdph.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// parse the body for the input
	t := new(createDTO)
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input, only title is required, text is optional
	if t.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// create the item we want to insert
	todo := ToDo{
		primitive.NewObjectID(),
		t.Title,
		t.Text,
		false,
	}

	// insert the item to the database
	_, err = todosCollection.InsertOne(r.Context(), todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// insert the id of the item to the users collections array
	_, err = usersCollection.UpdateOne(r.Context(), bson.M{"_id": userId}, bson.M{"$push": bson.M{"todos": todo.ID}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// return a success response
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	todosCollection := config.GetCollection("todos")
	usersCollection := config.GetCollection("users")

	userIdph := r.Context().Value("userId")
	userId, err := primitive.ObjectIDFromHex(userIdph.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// get todos id
	idPlaceHolder := r.PathValue("id")
	id, err := primitive.ObjectIDFromHex(idPlaceHolder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if the item belongs to the user
	u := new(user.User)
	err = usersCollection.FindOne(r.Context(), bson.M{"_id": userId}).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	isOwnedByUser := slices.Contains(u.Todos, id)
	if !isOwnedByUser {
		http.Error(w, "No access", http.StatusBadRequest)
	}

	update := new(updateDTO)
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// update the document
	result, err := todosCollection.UpdateOne(r.Context(), bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Remove(w http.ResponseWriter, r *http.Request) {
	todosCollection := config.GetCollection("todos")
	usersCollection := config.GetCollection("users")

	userIdph := r.Context().Value("userId")
	userId, err := primitive.ObjectIDFromHex(userIdph.(string))

	// get todos id
	idPlaceHolder := r.PathValue("id")
	id, err := primitive.ObjectIDFromHex(idPlaceHolder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if the item belongs to the user
	u := new(user.User)
	err = usersCollection.FindOne(r.Context(), bson.M{"_id": userId}).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	isOwnedByUser := slices.Contains(u.Todos, id)
	if !isOwnedByUser {
		http.Error(w, "No access", http.StatusBadRequest)
	}

	// delete the item
	result, err := todosCollection.DeleteOne(r.Context(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
