package todo

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	// get the users ID
	userID, err := primitive.ObjectIDFromHex(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get the items of the user
	result, err := GetUserTodos(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the results
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	// get the user id from context
	userID, err := primitive.ObjectIDFromHex(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// parse the body for the input
	t := createDTO{}
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input, only title is required, text is optional
	if t.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// create the item by calling the service layer
	todo, err := CreateTodo(r.Context(), userID, t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return a success response
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	// get user id
	userID, err := primitive.ObjectIDFromHex(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get todos id
	todoID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the fields we need to update
	update := updateDTO{}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if user is allowed to do the update
	allowed, err := BelongsToUser(r.Context(), userID, todoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// call the method to update the item
	result, err := UpdateTodo(r.Context(), todoID, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// on success
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Remove(w http.ResponseWriter, r *http.Request) {
	// get the user id
	userID, err := primitive.ObjectIDFromHex(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get the item id
	todoID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if user is allowed to do the update
	allowed, err := BelongsToUser(r.Context(), userID, todoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// call the method to delete the item
	result, err := DeleteTodo(r.Context(), userID, todoID)
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

func ReOrder(w http.ResponseWriter, r *http.Request) {
	userID, err := primitive.ObjectIDFromHex(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get the input array
	ids := ReOrderDTO{}
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var newOrder []primitive.ObjectID
	for _, id := range ids.IDs {
		idObj, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newOrder = append(newOrder, idObj)
	}

	sorted, err := SortTodos(r.Context(), userID, newOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !sorted {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
