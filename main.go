package main

import (
	"ToDo/api"
	"ToDo/api/db"
)

func main() {

	// connect to database
	db.DBConnect()
	// start the server
	api.Server()

}
