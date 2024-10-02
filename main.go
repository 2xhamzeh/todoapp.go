package main

import (
	"ToDo/api"
	"ToDo/api/config"
)

func main() {

	// connect to database
	config.DBConnect()
	// start the server
	api.Server()

}
