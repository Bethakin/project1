package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Bethakin/project1/api/handler"
	"github.com/Bethakin/project1/internal/database"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	todoHandler := handler.NewTodoHandler(db)

	router := mux.NewRouter()
	godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbPassword := os.Getenv("DB_PASSWORD")
	sslMode := os.Getenv("DB_SSLMODE")
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("DB_USER is not set in .env file")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME is not set in .env file")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, sslMode)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
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
	defer conn.Close(context.Background())

}
