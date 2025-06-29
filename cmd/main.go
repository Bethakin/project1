package main

import (
	_ "context"
	_ "fmt"
	"log"
	_ "net/http"
	"os"

	"github.com/Bethakin/project1/api/handler"
	"github.com/Bethakin/project1/internal/database"
	_ "github.com/Bethakin/project1/internal/repository"
	utils "github.com/Bethakin/project1/jwt"
	_ "github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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

	e := echo.New()

	userHandler := handler.NewUserHandler(db)
	todoHandler := handler.NewTodoHandler(db)

	e.POST("/register", userHandler.CreateUser)
	e.POST("/login", userHandler.LoginUser)

	protected := e.Group("/users/:users_id", utils.AuthMiddleware(jwtSecret))

	protected.GET("", userHandler.Index)
	protected.PUT("", userHandler.Updateusers)
	protected.DELETE("", userHandler.Deleteusers)

	todo := protected.Group("/todos")
	todo.GET("", todoHandler.IndexTodo)
	todo.POST("", todoHandler.CreateTodo)
	todo.GET("/:id", todoHandler.ShowTodo)
	todo.PUT("/:id", todoHandler.Update)
	todo.DELETE("/:id", todoHandler.Delete)

	e.Logger.Fatal(e.Start(":8081"))
	/*
		todoHandler := handler.NewTodoHandler(db)
		userHandler := handler.NewUserHandler(db)

		router := mux.NewRouter()
		router.HandleFunc("/register", userHandler.CreateUser).Methods("POST")
		router.HandleFunc("/login", userHandler.LoginUser).Methods("POST")

		protected := router.PathPrefix("/users").Subrouter()
		//protected.Use(utils.ValidateToken(jwtSecret))
		protected.Use(utils.AuthMiddleware(jwtSecret))

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
		}*/
}
