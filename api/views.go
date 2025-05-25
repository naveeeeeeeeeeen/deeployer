package api

import (
	"deeployer/tables"
	"encoding/json"
	"fmt"
	"net/http"
)

func getAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := tables.GetAllConfigs()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println("error getting projects", err)
		http.Error(w, "error getting projects", http.StatusBadRequest)
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
