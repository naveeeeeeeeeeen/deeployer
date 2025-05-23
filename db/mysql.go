package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func MysqlConnection() {
	dsn := "root:@tcp(127.0.0.1:3306)/pipes"
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB connection error", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("DB is not available for connections", err)
	}

	row, err := DB.Query("SELECT version()")
	if err != nil {
		log.Printf("ERROR IN QUERY", err)
	}
	defer row.Close()

	var version string
	for row.Next() {
		err := row.Scan(&version)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Mysql Version: ", version)
}
