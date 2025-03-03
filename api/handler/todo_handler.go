package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bethakin/project1/internal/database"
	"github.com/Bethakin/project1/internal/model"
	"github.com/gorilla/mux"
)

type TodoHandler struct {
	db *database.Database
}

func NewTodoHandler(db *database.Database) *TodoHandler {
	return &TodoHandler{
		db: db,
	}
}

func (h *TodoHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	todos, err := h.db.GetAllTodos()
	if err != nil {
		http.Error(w, "Error fetching todos", http.StatusInternalServerError)
		return
	}

	response := model.Response{
		Message: "Todos retrieved successfully",
		Data:    todos,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) Show(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	todo, err := h.db.GetTodoByID(id)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	response := model.Response{
		Message: "Todo retrieved successfully",
		Data:    todo,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	todo := &model.TodoUser{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.db.CreateTodo(todo); err != nil {
		http.Error(w, "Error creating todo", http.StatusInternalServerError)
		return
	}

	response := model.Response{
		Message: "Todo created successfully",
		Data:    todo,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteTodo(id); err != nil {
		http.Error(w, "Error deleting todo", http.StatusInternalServerError)
		return
	}

	response := model.Response{
		Message: "Todo deleted successfully",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	todo := &model.TodoUser{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.db.UpdateTodo(id, todo); err != nil {
		http.Error(w, "Error updating todo", http.StatusInternalServerError)
		return
	}

	response := model.Response{
		Message: "Todo updated successfully",
		Data:    todo,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
