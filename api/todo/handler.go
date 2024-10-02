package todo

import (
	"ToDo/api/config"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	todosCollection := config.GetCollection("todos")

	userIdph := r.Context().Value("userId")
	userId, err := primitive.ObjectIDFromHex(userIdph.(string))

	cursor, err := todosCollection.Find(r.Context(), bson.M{"user_id": userId})
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
	todosCollection := config.GetCollection("todos")

	userIdph := r.Context().Value("userId")
	userId, err := primitive.ObjectIDFromHex(userIdph.(string))

	// parse the body
	t := new(createDTO)
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input
	if t.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// insert the todo item to the database
	result, err := todosCollection.InsertOne(r.Context(), bson.M{
		"user_id": userId,
		"title":   t.Title,
		"text":    t.Text,
		"done":    false,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// return a success response
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	todosCollection := config.GetCollection("todos")

	userIdph := r.Context().Value("userId")
	userId, err := primitive.ObjectIDFromHex(userIdph.(string))

	// get todos id
	idPlaceHolder := r.PathValue("id")
	id, err := primitive.ObjectIDFromHex(idPlaceHolder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	update := new(updateDTO)
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// update the document
	result, err := todosCollection.UpdateOne(r.Context(), bson.M{"_id": id, "user_id": userId}, bson.M{"$set": update})
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

	userIdph := r.Context().Value("userId")
	userId, err := primitive.ObjectIDFromHex(userIdph.(string))

	// get todos id
	idPlaceHolder := r.PathValue("id")
	id, err := primitive.ObjectIDFromHex(idPlaceHolder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// delete the todo
	result, err := todosCollection.DeleteOne(r.Context(), bson.M{"_id": id, "user_id": userId})
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
