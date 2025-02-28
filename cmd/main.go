package main

import (
	"log"
	"net/http"

	"github.com/Bethakin/project1/api/handler"
	"github.com/gorilla/mux"
)

func main() {
	todoHandler := handler.NewTodoHandlerDB()
	router := mux.NewRouter()

	router.HandleFunc("/todos", todoHandler.Index).Methods("GET")
	router.HandleFunc("/todos", todoHandler.Create).Methods("POST")
	router.HandleFunc("/todos/{id}", todoHandler.Delete).Methods("DELETE")
	router.HandleFunc("/todos/{id}", todoHandler.Update).Methods("PUT")

	log.Println("Server starting on port 8081...")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
