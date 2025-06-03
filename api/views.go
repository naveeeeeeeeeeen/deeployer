package api

import (
	"deeployer/db"
	"deeployer/funcs"
	"deeployer/tables"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func getAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := tables.GetAllConfigs()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println("error getting projects", err)
		http.Error(w, "error getting projects", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(projects)
}

func GetProjects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllProjects(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func buildByProjectId(w http.ResponseWriter, id string) {
	w.Header().Set("Content-Type", "application/json")
	err := funcs.Build(id)
	if err != nil {
		log.Println("Error building, err", err)
		http.Error(w, "eror building", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "success"})
	return
}

func DeployProject(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		projectId := r.URL.Query().Get("projectId")
		buildByProjectId(w, projectId)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	var body map[string]string
	json.NewDecoder(r.Body).Decode(&body)
	username := body["username"]
	pass := body["password"]
	var response AppResponse
	user, err := tables.GetUserByUsername(username)
	if err != nil {
		response.status = 0
		response.message = "Something went wrong"
	} else {
		if user.CheckPassword(pass) {
			user.CreateUserToken()
			err := db.RedisSet(user.Token, "")
			if err != nil {
				response.message = "error generating token"
				response.status = 0
			}
			response.data = user.Json()
			json.NewEncoder(w).Encode(response)
			return
		}
		response.status = 0
		if len(user.UserName) > 0 {
			response.message = "No user found"
		} else {
			response.message = "Wrong password"
		}
		json.NewEncoder(w).Encode(response)
		return
	}

}

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		loginUser(w, r)
		return
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
