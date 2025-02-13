package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Todo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

var todos []Todo

func Show_todos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func Create_todo(w http.ResponseWriter, r *http.Request) {
	var newtodo Todo
	_ = json.NewDecoder(r.Body).Decode(&newtodo)
	todos = append(todos, newtodo)

}

func Delete_todo(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	title := parameters["title"]
	for i, todo := range todos {
		if todo.Title == title {
			todos = append(todos[:i], todos[i+1:]...)
			return
		}
	}
}

func Update_todo(w http.ResponseWriter, r *http.Request) {
	var uptodo Todo
	_ = json.NewDecoder(r.Body).Decode(&uptodo)
	parameters := mux.Vars(r)
	title := parameters["title"]
	for i, todo := range todos {
		if todo.Title == title {
			todos[i].Description = uptodo.Description
		}
	}
}

/*
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/todos", show_todos).Methods("GET")
	router.HandleFunc("/todos", create_todo).Methods("POST")
	router.HandleFunc("/todos/{title}", delete_todo).Methods("DELETE")
	router.HandleFunc("/todos/{title}", update_todo).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", router))
}*/
