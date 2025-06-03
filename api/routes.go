package api

import (
	"fmt"
	"log"
	"net/http"
)

func routes() {
	HandleWithCors("/projects", GetProjects)
	HandleWithCors("/deploy", DeployProject)
	HandleWithCors("/login", Login)
}

func HandleWithCors(patter string, handlerFunc http.HandlerFunc) {
	newHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handlerFunc(w, r)
	})
	http.HandleFunc(patter, newHandler)
}

func Serve(port string) error {
	routes()
	log.Println("registered routes")
	err := http.ListenAndServe(":"+string(port), nil)
	if err != nil {
		fmt.Println("error : ", err)
		return err
	}
	return nil
}
