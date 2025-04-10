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
	router.HandleFunc("/register", userHandler.CreateUser).Methods("POST")
	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", userHandler.Index).Methods("GET")
	//userRouter.HandleFunc("", userHandler.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{users_id}", userHandler.Updateusers).Methods("PUT")
	userRouter.HandleFunc("/{users_id}", userHandler.Deleteusers).Methods("DELETE")

	todoRouter := userRouter.PathPrefix("/{users_id}/todos").Subrouter()
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
