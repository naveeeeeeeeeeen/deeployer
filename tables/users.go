package tables

import (
	"database/sql"
	"deeployer/db"
	"fmt"
)

type User struct {
	ID   int
	Name string
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("error quering db", err)
	}

	var users []User

	for rows.Next() {
		fmt.Println("ASDASDSAD")
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, fmt.Errorf("error scaning row", err)
		}
		users = append(users, u)
	}

	return users, nil
}

type Users []User

func (u Users) insertQuery() error {
	q := "insert into users (`name`) values "

	for i := 0; i < len(u); i += 1 {
		q += fmt.Sprintf("('%s'),", u[i].Name)
	}
	fmt.Println(q[:len(q)-1])
	_, err := db.DB.Query(q[:len(q)-1])
	if err != nil {
		return fmt.Errorf("error inserting users %v", err)
	}
	return nil
}
