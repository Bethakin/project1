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

type UserHandler struct {
	db *database.Database
}

func NewUserHandler(db *database.Database) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.db.GetAlluserss()
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	response := model.Response{
		Message: "Users retrieved successfully",
		Data:    users,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
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
	user := &model.Todousers{
		Email:    req.Email,
		Password: req.Password,
	}
	if err := h.db.CreateUser(user); err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	response := model.Response{
		Message: "User created successfully",
		Data:    user,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) Deleteusers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	user_id, err := strconv.Atoi(params["users_id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	if err := h.db.DeleteUserTodos(user_id); err != nil {
		http.Error(w, "Error deleting users todos", http.StatusInternalServerError)
		result := fmt.Sprintf(err.Error())
		fmt.Println(result)
		return
	}
	if err := h.db.DeleteUser(user_id); err != nil {
		http.Error(w, "Error deleting users", http.StatusInternalServerError)
		return
	}
	response := model.Response{
		Message: "users deleted successfully",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) Updateusers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["users_id"])
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

	if err := h.db.UpdateUser(id, todo); err != nil {
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
