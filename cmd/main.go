package main

import (
	"log"
	"net/http"

	"github.com/Bethakin/project1/api/handler"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", handler.Welcome).Methods("GET")
	router.HandleFunc("/todos", handler.Show_todos).Methods("GET")
	router.HandleFunc("/todos", handler.Create_todo).Methods("POST")
	router.HandleFunc("/todos/{id}", handler.Delete_todo).Methods("DELETE")
	router.HandleFunc("/todos/{id}", handler.Update_todo).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", router))
}
