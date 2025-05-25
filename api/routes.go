package api

import (
	"fmt"
	"net/http"
)

func routes() {
	http.HandleFunc("/projects", GetProjects)
	http.HandleFunc("/deploy", DeployProject)
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
