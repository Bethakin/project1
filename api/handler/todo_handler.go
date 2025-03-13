package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Bethakin/project1/internal/database"
	"github.com/Bethakin/project1/model"
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

	todos, err := h.db.GetAlluserss()
	if err != nil {
		http.Error(w, "Error fetching todos", http.StatusInternalServerError)
		return
	}

	response := model.Response{
		Message: "userss retrieved successfully",
		Data:    todos,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) IndexTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	todos, err := h.db.GetAllTodos()
	if err != nil {
		result := fmt.Sprintf(err.Error())
		fmt.Println(result)
		http.Error(w, "Error fetching todos", http.StatusInternalServerError)
		return
	}

	response := model.Response{
		Message: "userss retrieved successfully",
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
	id, err := strconv.Atoi(params["users_id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	todo, err := h.db.GetusersByID(id)
	if err != nil {
		http.Error(w, "users not found", http.StatusNotFound)
		return
	}

	response := model.Response{
		Message: "users retrieved successfully",
		Data:    todo,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) ShowTodo(w http.ResponseWriter, r *http.Request) {
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

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	UsersID, _ := strconv.Atoi(params["users_id"])

	var req model.TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	todo := &model.Todo{
		Title:       req.Title,
		Description: req.Description,
		UsersID:     UsersID,
	}

	if err := h.db.CreateTodo(todo); err != nil {
		http.Error(w, "Error creating todo", http.StatusInternalServerError)
		result := fmt.Sprintf(err.Error())
		fmt.Println(result)

		return
	}

	response := model.Response{
		Message: "Todo created successfully",
		Data:    todo,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *TodoHandler) Createusers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req model.TodoRequestusers
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	todo := &model.Todousers{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.db.Createusers(todo); err != nil {
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

func (h *TodoHandler) Deleteusers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	user_id, err := strconv.Atoi(params["users_id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteuserssTodo(user_id); err != nil {
		http.Error(w, "Error deleting users todos", http.StatusInternalServerError)
		result := fmt.Sprintf(err.Error())
		fmt.Println(result)

		return
	}

	if err := h.db.DeleteTheusers(user_id); err != nil {
		http.Error(w, "Error deleting users", http.StatusInternalServerError)
		return
	}

	response := model.Response{
		Message: "users deleted successfully",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *TodoHandler) Updateusers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["users_id"])
	//users_id, _ := strconv.Atoi(params["ıser_id"])

	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var req model.TodoRequestusers
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	todo := &model.Todousers{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.db.Updateusers(id, todo); err != nil {
		result := fmt.Sprintf(err.Error())
		fmt.Println(result)
		http.Error(w, "Error updating users", http.StatusInternalServerError)
		return
	}

	response := model.Response{
		Message: "users updated successfully",
		Data:    todo,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	users_id, _ := strconv.Atoi(params["users_id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var req model.TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	todo := &model.Todo{
		Title:       req.Title,
		Description: req.Description,
	}

	if err := h.db.UpdateTodo(id, users_id, todo); err != nil {
		http.Error(w, "Error updating users", http.StatusInternalServerError)
		return
	}

	response := model.Response{
		Message: "users updated successfully",
		Data:    todo,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
