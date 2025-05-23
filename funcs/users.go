package funcs

import (
	"deeployer/tables"
	"fmt"
)

// users functionalities

func CreateUser(name string) (tables.User, error) {
	user := tables.User{
		Name: name,
	}
	users := tables.Users{
		user,
	}

	err := tables.InsertQuery(users)
	if err != nil {
		return user, fmt.Errorf("error inserting users %v", err)
	}
	return user, nil
}
