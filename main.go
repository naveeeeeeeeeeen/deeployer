package main

import (
	"deeployer/db"
	"deeployer/funcs"
	"deeployer/tables"
	"fmt"
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
	projectId := 2
	fmt.Println("Building dist files for a react project")
	err := funcs.Build(projectId)
	if err != nil {
		fmt.Println("Error buildling, err: ", err)
	}
}
