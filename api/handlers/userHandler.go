package handlers

import (
	"ToDo/api/models"
	"ToDo/api/services"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	// get the username and password
	user := models.AuthDTO{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the input
	// note len() returns byte length not string length
	if user.Password == "" || user.Username == "" || len(user.Password) < 6 || len(user.Username) < 3 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// check if username exists
	exists, err := services.UsernameExists(r.Context(), user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Username taken", http.StatusBadRequest)
		return
	}

	// save user
	err = services.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// get the input
	user := models.AuthDTO{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate input
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	token, err := services.LoginUser(r.Context(), user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+*token)
	w.WriteHeader(http.StatusOK)
}
