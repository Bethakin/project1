package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Bethakin/project1/api/handler"
	"github.com/Bethakin/project1/api/todo"
	"github.com/gorilla/mux"
)

func main() {
	http.HandleFunc("/", handler.Welcome)
	router := mux.NewRouter()
	router.HandleFunc("/todos", todo.Show_todos).Methods("GET")
	router.HandleFunc("/todos", todo.Create_todo).Methods("POST")
	router.HandleFunc("/todos/{title}", todo.Delete_todo).Methods("DELETE")
	router.HandleFunc("/todos/{title}", todo.Update_todo).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", router))
	fmt.Println("Server starting on port :8081...")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
		return
	}
}
