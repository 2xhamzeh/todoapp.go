package handlers

import (
	"ToDo/api/services"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// Shares the authenticated users todos with a user
func HandleShareWithUser(w http.ResponseWriter, r *http.Request) {
	// get user id
	userID, err := primitive.ObjectIDFromHex(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usernameToShareWith := r.PathValue("username")

	// check if username exists
	exists, err := services.UsernameExists(r.Context(), usernameToShareWith)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Username does not exist", http.StatusNotFound)
		return
	}

	// check if username has already been shared with
	shared, err := services.CheckIfAuthUserSharedWithUsername(r.Context(), usernameToShareWith, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if shared {
		http.Error(w, "Already shared with provided username", http.StatusNotFound)
		return
	}

	// shared my todos with user
	err = services.ShareWithUser(r.Context(), usernameToShareWith, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func HandleUnshareWithUser(w http.ResponseWriter, r *http.Request) {
	// get user id
	userID, err := primitive.ObjectIDFromHex(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usernameToShareWith := r.PathValue("username")

	// check if username exists
	exists, err := services.UsernameExists(r.Context(), usernameToShareWith)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Username does not exist", http.StatusNotFound)
		return
	}

	// check if username has already been shared with
	shared, err := services.CheckIfAuthUserSharedWithUsername(r.Context(), usernameToShareWith, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !shared {
		http.Error(w, "Your todos are not shared with provided username", http.StatusNotFound)
		return
	}

	// unshared my todos with user
	err = services.UnshareWithUser(r.Context(), usernameToShareWith, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func HandleGetSharedTodosFromUser(w http.ResponseWriter, r *http.Request) {
	// get user id
	userID, err := primitive.ObjectIDFromHex(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username := r.PathValue("username")

	// check if username exists
	exists, err := services.UsernameExists(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Username does not exist", http.StatusNotFound)
		return
	}

	// check if username shared his todos with us
	shared, err := services.CheckIfUsernameSharedWithAuthUser(r.Context(), username, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !shared {
		http.Error(w, "The user hasn't shared their todos with you", http.StatusForbidden)
		return
	}

	id, err := services.GetIDFromUsername(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := services.GetTodos(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleGetUsersSharedWithMe(w http.ResponseWriter, r *http.Request) {
	userID, err := primitive.ObjectIDFromHex(r.Context().Value("userID").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	listOfUsernames, err := services.GetUsernamesSharedWithAuth(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(listOfUsernames); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func HandleChangeDoneOfSharedTodo(w http.ResponseWriter, r *http.Request) {
	
}
