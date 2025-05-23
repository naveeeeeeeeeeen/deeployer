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

	users := getAllUsers()

	fmt.Println(users)

	_, err := funcs.CreateUser("Naveen")
	if err != nil {
		panic(err)
	}

	c := funcs.CreateConfig(1, "SSH KEY", "GITHUB KEY", "PROJECT NAME", "REPO URL", "jump-aws-staging.playo.io", "naveen")
	funcs.Connect(c)

}
