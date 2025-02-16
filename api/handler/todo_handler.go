package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/Bethakin/project1/model"
	"github.com/gorilla/mux"
)

type TodoHandler struct {
	todos  map[int]*model.Todo
	nextID int
	mutex  sync.RWMutex
}

func NewTodoHandler() *TodoHandler {
	return &TodoHandler{
		todos:  make(map[int]*model.Todo),
		nextID: 1,
		mutex:  sync.RWMutex{},
	}
}

func (h *TodoHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	h.mutex.RLock()
	todoList := make([]*model.Todo, 0, len(h.todos))
	for _, todo := range h.todos {
		todoList = append(todoList, todo)
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

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req model.TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	h.mutex.Lock()
	todo := &model.Todo{
		ID:          h.nextID,
		Title:       req.Title,
		Description: req.Description,
	}
	h.todos[h.nextID] = todo
	h.nextID++
	h.mutex.Unlock()

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

	h.mutex.Lock()
	if _, exists := h.todos[id]; !exists {
		h.mutex.Unlock()
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	delete(h.todos, id)
	h.mutex.Unlock()

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

	var req model.TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.mutex.Lock()
	todo, exists := h.todos[id]
	if !exists {
		h.mutex.Unlock()
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	if req.Title != "" {
		todo.Title = req.Title
	}
	if req.Description != "" {
		todo.Description = req.Description
	}
	h.mutex.Unlock()

	response := model.Response{
		Message: "Todo updated successfully",
		Data:    todo,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
