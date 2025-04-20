package main

import (
	_ "context"
	_ "fmt"
	"log"
	"net/http"
	"os"

	"github.com/Bethakin/project1/api/handler"
	"github.com/Bethakin/project1/internal/database"
	_ "github.com/Bethakin/project1/internal/repository"
	utils "github.com/Bethakin/project1/jwt"
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in the environment")
	}

	todoHandler := handler.NewTodoHandler(db)
	userHandler := handler.NewUserHandler(db)

	router := mux.NewRouter()
	router.HandleFunc("/register", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/login", userHandler.LoginUser).Methods("POST")

	protected := router.PathPrefix("/users").Subrouter()
	protected.Use(utils.ValidateToken(jwtSecret))

	protected.HandleFunc("", userHandler.IndexAll).Methods("GET")
	protected.HandleFunc("/{users_id}", userHandler.Index).Methods("GET")
	protected.HandleFunc("/{users_id}", userHandler.Updateusers).Methods("PUT")
	protected.HandleFunc("/{users_id}", userHandler.Deleteusers).Methods("DELETE")

	todoRouter := protected.PathPrefix("/{users_id}/todos").Subrouter()
	todoRouter.HandleFunc("", todoHandler.IndexTodo).Methods("GET")
	todoRouter.HandleFunc("", todoHandler.CreateTodo).Methods("POST")
	todoRouter.HandleFunc("/{id}", todoHandler.ShowTodo).Methods("GET")
	todoRouter.HandleFunc("/{id}", todoHandler.Update).Methods("PUT")
	todoRouter.HandleFunc("/{id}", todoHandler.Delete).Methods("DELETE")

	log.Println("Server starting on port 8081...")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
