package main

import (
	"log"
	"net/http"

	"github.com/Bethakin/project1/api/handler"
	"github.com/Bethakin/project1/internal/database"
	"github.com/gorilla/mux"
)

// HİÇBİRİ KONTROL EDİLMEDİ VE UPDATEİ DÜZENLKEMEDİM
func main() {
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	todoHandler := handler.NewTodoHandler(db)

	router := mux.NewRouter()

	router.HandleFunc("/users", todoHandler.Index).Methods("GET")
	router.HandleFunc("/users/{users_id}", todoHandler.Show).Methods("GET")
	router.HandleFunc("/users/{users_id}/todos", todoHandler.IndexTodo).Methods("GET")
	router.HandleFunc("/users/{users_id}/todos/{id}", todoHandler.ShowTodo).Methods("GET")
	router.HandleFunc("/users", todoHandler.Createusers).Methods("POST")
	router.HandleFunc("/users/{users_id}/todos", todoHandler.CreateTodo).Methods("POST")
	router.HandleFunc("/users/{users_id}", todoHandler.Updateusers).Methods("PUT")
	router.HandleFunc("/users/{users_id}/todos/{id}", todoHandler.Update).Methods("PUT")
	router.HandleFunc("/users/{users_id}", todoHandler.Deleteusers).Methods("DELETE")
	router.HandleFunc("/users/{users_id}/todos/{id}", todoHandler.Delete).Methods("DELETE")

	log.Println("Server starting on port 8081...")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
