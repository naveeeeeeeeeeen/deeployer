package tables

import (
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Insert interface {
	insertQuery() error
}

func InsertQuery(i Insert) error {
	err := i.insertQuery()
	return err
}

func insert() {
	config := Configs{}
	InsertQuery(config)
}

func GenerateHashedPass(pass string) string {
	password := []byte(pass)

	hashedPassord, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hasheing password", err)
	}
	return string(hashedPassord)
}

func CheckPass(pass string, hash string) bool {
	password := []byte(pass)
	err := bcrypt.CompareHashAndPassword([]byte(hash), password)
	if err != nil {
		return false
	}
	return true
}

func GenerateUserToken() string {
	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Println("error getting a token", err)
	}
	return uuid.String()
}
