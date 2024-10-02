package api

import (
	"ToDo/api/middleware"
	"ToDo/api/todo"
	"ToDo/api/user"
	"log"
	"net/http"
)

func router() *http.ServeMux {
	myRouter := http.NewServeMux()

	myRouter.HandleFunc("GET /todo", middleware.AuthMiddleware(todo.GetAll))
	myRouter.HandleFunc("POST /todo", middleware.AuthMiddleware(todo.Create))
	myRouter.HandleFunc("PUT /todo/{id}", middleware.AuthMiddleware(todo.Update))
	myRouter.HandleFunc("DELETE /todo/{id}", middleware.AuthMiddleware(todo.Remove))

	myRouter.HandleFunc("POST /register", user.Register)
	myRouter.HandleFunc("POST /login", user.Login)
	myRouter.HandleFunc("POST /logout", user.Logout)

	return myRouter
}

func Server() {
	log.Fatal(http.ListenAndServe(":8080", router()))
}
