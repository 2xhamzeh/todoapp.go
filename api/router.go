package api

import (
	"ToDo/api/handlers"
	"ToDo/api/middleware"
	"log"
	"net/http"
)

func router() *http.ServeMux {
	myRouter := http.NewServeMux()

	myRouter.HandleFunc("GET /todo", middleware.AuthMiddleware(handlers.HandleGetTodos))
	myRouter.HandleFunc("POST /todo", middleware.AuthMiddleware(handlers.HandleCreateTodo))
	myRouter.HandleFunc("PUT /todo/{id}", middleware.AuthMiddleware(handlers.HandleUpdateTodo))
	myRouter.HandleFunc("DELETE /todo/{id}", middleware.AuthMiddleware(handlers.HandleRemoveTodo))
	myRouter.HandleFunc("PUT /todo/reorder", middleware.AuthMiddleware(handlers.HandleReorderTodos))

	myRouter.HandleFunc("POST /register", handlers.HandleRegister)
	myRouter.HandleFunc("POST /login", handlers.HandleLogin)

	myRouter.HandleFunc("GET /category", middleware.AuthMiddleware(handlers.HandleGetCategories))           // gets all categories
	myRouter.HandleFunc("GET /category/{name}", middleware.AuthMiddleware(handlers.HandleGetCategoryTodos)) // gets specific category's content
	myRouter.HandleFunc("POST /category/{name}", middleware.AuthMiddleware(handlers.HandleCreateCategory))
	myRouter.HandleFunc("DELETE  /category/{name}", middleware.AuthMiddleware(handlers.HandleDeleteCategory))
	myRouter.HandleFunc("POST  /category/{name}/{id}", middleware.AuthMiddleware(handlers.HandleAddTodoToCategory))
	myRouter.HandleFunc("DELETE   /category/{name}/{id}", middleware.AuthMiddleware(handlers.HandleDeleteTodoFromCategory))

	return myRouter
}

func Server() {
	log.Fatal(http.ListenAndServe(":8080", router()))
}
