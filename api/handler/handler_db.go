package handler

import (
	_ "context"
	_ "log"

	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/Bethakin/project1/model"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

type TodoHandlerDB struct {
	DB     *sql.DB
	todos  map[int]*model.TodoUser
	nextID int
	mutex  sync.RWMutex
}

func NewTodoHandlerDB() *TodoHandlerDB {
	connStr := "user=postgres password=27Berfin72 dbname=users sslmode=disable"
	db, _ := sql.Open("postgres", connStr)
	return &TodoHandlerDB{
		todos:  make(map[int]*model.TodoUser),
		nextID: 1,
		mutex:  sync.RWMutex{},
		DB:     db,
	}
}

func (h *TodoHandlerDB) Index(w http.ResponseWriter, r *http.Request) {
	rows, _ := h.DB.Query("SELECT id, email, password FROM todos")

	defer rows.Close()
	w.Header().Set("Content-Type", "application/json")

	h.mutex.RLock()
	todoList := make([]*model.TodoUser, 0, len(h.todos))

	for rows.Next() {
		var todo model.TodoUser
		if err := rows.Scan(&todo.ID, &todo.Email, &todo.Password); err != nil {
			http.Error(w, "Error scanning todos", http.StatusInternalServerError)
			return
		}
		todoList = append(todoList, &todo)

	}
	h.mutex.RUnlock()

	response := model.Response{
		Message: "Todos retrieved successfully",
		Data:    todoList,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandlerDB) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req model.TodoRequestUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	h.mutex.Lock()
	todo := &model.TodoUser{
		ID:       h.nextID,
		Email:    req.Email,
		Password: req.Password,
	}
	h.todos[h.nextID] = todo
	h.nextID++

	var id int
	h.DB.QueryRow("INSERT INTO todos (email, password) VALUES ($1, $2) RETURNING id", todo.Email, todo.Password).Scan(&id)

	todo.ID = id
	h.mutex.Unlock()

	response := model.Response{
		Message: "Todo created successfully",
		Data:    todo,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *TodoHandlerDB) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	h.DB.Exec("DELETE FROM todos WHERE id = $1", id)

	response := model.Response{
		Message: "Todo deleted successfully",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *TodoHandlerDB) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var req model.TodoRequestUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	h.DB.Exec("UPDATE todos SET email = $1, password = $2 WHERE id = $3", req.Email, req.Password, id)

	response := model.Response{
		Message: "Todo updated successfully",
		Data:    req,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
