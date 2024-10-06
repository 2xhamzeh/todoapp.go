package handlers

import (
	"ToDo/api/db"
	"ToDo/api/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// Shares the authenticated users todos with a user
func HandleShareWithUser(w http.ResponseWriter, r *http.Request) {
	usersCollection := db.GetCollection("users")

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
	if exists {
		http.Error(w, "Username taken", http.StatusBadRequest)
		return
	}

}
func HandleUnshareWithUser(w http.ResponseWriter, r *http.Request)        {}
func HandleGetSharedTodosFromUser(w http.ResponseWriter, r *http.Request) {}
func HandleGetUsersSharedWithMe(w http.ResponseWriter, r *http.Request)   {}
func HandleChangeDoneOfSharedTodo(w http.ResponseWriter, r *http.Request) {}
