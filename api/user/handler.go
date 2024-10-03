package user

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	// get the username and password
	user := AuthDTO{}
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
	err := UsernameExists(r.Context(), user.Username)
	if err == nil {
		http.Error(w, "Username taken", http.StatusBadRequest)
		return
	}

	// save user
	err = CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// get the input
	user := AuthDTO{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate input
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	token, err := LoginUser(r.Context(), user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+*token)
	w.WriteHeader(http.StatusOK)
}
