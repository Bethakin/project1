package main

import (
	"log"
	"net/http"

	"github.com/Bethakin/project1/api/handler"
	"github.com/Bethakin/project1/internal/database"
	"github.com/gorilla/mux"
)

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	todoHandler := handler.NewTodoHandler(db)

	router := mux.NewRouter()

	router.HandleFunc("/todos", todoHandler.Index).Methods("GET")
	router.HandleFunc("/todos/{id}", todoHandler.Show).Methods("GET")
	router.HandleFunc("/todos", todoHandler.Create).Methods("POST")
	router.HandleFunc("/todos/{id}", todoHandler.Delete).Methods("DELETE")
	router.HandleFunc("/todos/{id}", todoHandler.Update).Methods("PUT")

	log.Println("Server starting on port 8081...")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
