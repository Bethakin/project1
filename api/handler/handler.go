package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Welcome(w http.ResponseWriter, req *http.Request) {
	value := req.URL.Query().Get("word")
	object := map[string]string{"message": fmt.Sprintf("Welcome to the API, %s", value)}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(object); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
