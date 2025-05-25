package api

import (
	"fmt"
	"net/http"
)

func routes() {
	http.HandleFunc("/projects", GetProjects)
}

func Serve(port string) error {
	routes()
	err := http.ListenAndServe(":"+string(port), nil)
	if err != nil {
		fmt.Println("error : ", err)
		return err
	}
	return nil
}
