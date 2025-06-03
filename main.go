package main

import (
	"context"
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

var ctx = context.Background()

func main() {
	db.MysqlConnection()
	var port string = "3001"
	log.Println("listening on port ", port)
	db.RedisInit()
	api.Serve(port)
}
