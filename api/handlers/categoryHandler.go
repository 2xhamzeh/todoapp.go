package handlers

import (
	"ToDo/api/services"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	// get user id
	userID, err := primitive.ObjectIDFromHex(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categoryName := r.PathValue("name")

	// check if it already exists
	err = services.CategoryExists(r.Context(), userID, categoryName)
	if err == nil {
		http.Error(w, "Category already exists", http.StatusConflict)
		return
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create  new category
	result, err := services.CreateCategory(r.Context(), userID, categoryName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleDeleteCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := primitive.ObjectIDFromHex(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categoryName := r.PathValue("name")

	err = services.DeleteCategory(r.Context(), userID, categoryName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func HandleAddTodoToCategory(w http.ResponseWriter, r *http.Request) {
	// get user id
	userID, err := primitive.ObjectIDFromHex(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categoryName := r.PathValue("name")
	idString := r.PathValue("id")

	// convert id to objectID
	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if category doesn't exist
	err = services.CategoryExists(r.Context(), userID, categoryName)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "Category doesn't exist", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check if user owns the todo
	allowed, err := services.TodoBelongsToUser(r.Context(), userID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, "Todo is not yours", http.StatusUnauthorized)
		return
	}

	// add item to category
	err = services.AddTodoToCategory(r.Context(), userID, categoryName, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func HandleDeleteTodoFromCategory(w http.ResponseWriter, r *http.Request) {
	// get user id
	userID, err := primitive.ObjectIDFromHex(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categoryName := r.PathValue("name")
	idString := r.PathValue("id")

	// convert id to objectID
	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if category doesn't exist
	err = services.CategoryExists(r.Context(), userID, categoryName)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "Category doesn't exist", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check if user owns the todo
	allowed, err := services.TodoBelongsToUser(r.Context(), userID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, "Todo is not yours", http.StatusUnauthorized)
		return
	}

	// remove item from category
	err = services.RemoveTodoFromCategory(r.Context(), userID, categoryName, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	userID, err := primitive.ObjectIDFromHex(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := services.GetUserCategories(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func HandleGetCategoryTodos(w http.ResponseWriter, r *http.Request) {
	userID, err := primitive.ObjectIDFromHex(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// get category name
	categoryName := r.PathValue("name")
	// check if category doesn't exist
	err = services.CategoryExists(r.Context(), userID, categoryName)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "Category doesn't exist", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// get array with category todo ids
	cat, err := services.GetCategory(r.Context(), userID, categoryName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result, err := services.GetTodosFromIDArray(r.Context(), cat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
