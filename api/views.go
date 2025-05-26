package api

import (
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
