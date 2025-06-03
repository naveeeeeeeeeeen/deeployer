package tables

import (
	"database/sql"
	"deeployer/db"
	"fmt"
	"log"
)

type User struct {
	ID       int
	Name     string
	Password string
	UserName string
	Token    string
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

func (u User) InsertUser() error {
	q := "insert into users (`name`, `username`, `password`) values (?,?,?);"
	_, err := db.DB.Query(q, u.Name, u.UserName, u.Password)
	if err != nil {
		log.Println("error inserting user", err)
		return err
	}
	log.Println("Inserted user in db")
	return nil
}

func GetUserByUsername(usernamme string) (User, error) {
	var user User
	query := "select id, name, password, username from users where username = ?;"

	rows, err := db.DB.Query(query, usernamme)

	defer rows.Close()
	if err != nil {
		log.Println("error reading users", err)
		return user, err
	}
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.Password, &user.UserName)
		if err != nil {
			log.Println("error getting values ", err)
			return user, err
		}
	}
	return user, nil
}

func (user User) CheckPassword(pass string) bool {
	matched := CheckPass(pass, user.Password)
	return matched
}

func (user User) Json() map[string]any {
	return map[string]any{
		"id":       user.ID,
		"username": user.UserName,
		"name":     user.Name,
		"token":    user.Token,
	}
}

func (user User) CreateUserToken() string {
	uuid := GenerateUserToken()
	return uuid
}
