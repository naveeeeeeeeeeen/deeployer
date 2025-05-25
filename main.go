package main

import (
	"deeployer/api"
	"deeployer/db"
	"deeployer/tables"
	"log"
)

func getAllUsers() tables.Users {
	users, err := tables.GetAllUsers(db.DB)
	if err != nil {
		log.Fatal("Error fetching users", err)
	}
	return users
}

func main() {
	db.MysqlConnection()
	var port string = "3001"
	log.Println("listening on port ", port)
	api.Serve(port)
}
