package main

import (
	_ "context"
	_ "fmt"
	"log"
	"net/http"
	_ "os"

	"github.com/Bethakin/project1/api/handler"
	"github.com/Bethakin/project1/internal/database"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/joho/godotenv"
)

func main() {

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if db == nil {
		log.Fatal("Database is nil after initialization!")
	}
	todoHandler := handler.NewTodoHandler(db)
	userHandler := handler.NewUserHandler(db)
	router := mux.NewRouter()
	router.HandleFunc("/users", userHandler.Index).Methods("GET")
	router.HandleFunc("/users/{users_id}", todoHandler.Show).Methods("GET")
	router.HandleFunc("/users/{users_id}/todos", todoHandler.IndexTodo).Methods("GET")
	router.HandleFunc("/users/{users_id}/todos/{id}", todoHandler.ShowTodo).Methods("GET")
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{users_id}/todos", todoHandler.CreateTodo).Methods("POST")
	router.HandleFunc("/users/{users_id}", userHandler.Updateusers).Methods("PUT")
	router.HandleFunc("/users/{users_id}/todos/{id}", todoHandler.Update).Methods("PUT")
	router.HandleFunc("/users/{users_id}", userHandler.Deleteusers).Methods("DELETE")
	router.HandleFunc("/users/{users_id}/todos/{id}", todoHandler.Delete).Methods("DELETE")

	log.Println("Server starting on port 8081...")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
