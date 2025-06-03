package api

import (
	"deeployer/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func routes() {
	HandleWithCors("/projects", GetProjects, true)
	HandleWithCors("/deploy", DeployProject, true)
	HandleWithCors("/login", Login, true)
	HandleWithCors("/home/", Home, true)
}

func HandleWithCors(patter string, handlerFunc http.HandlerFunc, auth bool) {
	newHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		var response AppResponse

		authHeader := r.Header.Get("authToken")
		val, err := db.RedisGet(authHeader)

		if err != nil {
			log.Println("error getting redis key ", err)
			response.status = 0
			response.message = "soemthing went wrong"
			json.NewEncoder(w).Encode(response.Json())
			return
		}
		if len(val) == 0 {
			response.status = -1
			response.message = "Not Authorized"
			json.NewEncoder(w).Encode(response.Json())
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
