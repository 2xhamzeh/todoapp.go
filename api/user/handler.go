package user

import (
	"ToDo/api/auth"
	"ToDo/api/config"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	userCollection := config.GetCollection("users")

	// get the username and password
	user := new(AuthDTO)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the input
	// note len() returns byte length not string length
	if user.Password == "" || user.Username == "" || len(user.Password) < 5 || len(user.Username) < 5 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// check if username is taken
	err := userCollection.FindOne(r.Context(), bson.M{"username": user.Username}).Err()
	if err == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	// save user
	result, err := userCollection.InsertOne(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	userCollection := config.GetCollection("users")

	// get the input
	input := new(AuthDTO)
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate input
	if input.Username == "" || input.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// check if user exists
	user := new(User)
	err := userCollection.FindOne(r.Context(), bson.M{"username": input.Username}).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid username", http.StatusNotFound)
		return
	}

	// check password
	isAuth := auth.CheckPasswordHash(input.Password, user.Password)
	if !isAuth {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// generate token
	token, err := auth.GenerateJWT(user.ID.Hex())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
func Logout(w http.ResponseWriter, r *http.Request) {}
