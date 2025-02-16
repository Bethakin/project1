package handler

import (
	"encoding/json"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type New_todo struct {
	//ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var (
	todos  = make(map[int]*Todo)
	nextID = 1
)

func Show_todos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func Create_todo(w http.ResponseWriter, r *http.Request) {
	var newtodo *New_todo
	_ = json.NewDecoder(r.Body).Decode(&newtodo)
	nextID += 1
	n_todo := &Todo{
		ID:          nextID,
		Title:       newtodo.Title,
		Description: newtodo.Description,
	}
	todos[nextID] = n_todo
}

func Delete_todo(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	id_str := parameters["id"]
	id, _ := strconv.Atoi(id_str)
	delete(todos, id)
}

func Update_todo(w http.ResponseWriter, r *http.Request) {
	var uptodo *Todo
	_ = json.NewDecoder(r.Body).Decode(&uptodo)
	parameters := mux.Vars(r)
	id_str := parameters["id"]
	id, _ := strconv.Atoi(id_str)
	todos[id].Description = uptodo.Description
	todos[id].Title = uptodo.Title
}
