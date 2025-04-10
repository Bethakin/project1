package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Bethakin/project1/internal/database"
	"github.com/Bethakin/project1/internal/repository"
	"github.com/Bethakin/project1/model"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(db *database.Database) *UserHandler {
	return &UserHandler{
		userRepo: repository.NewUserRepository(db),
	}
}

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.userRepo.GetAll()
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
	//hata kontrolü eksik???
	users, _ := h.userRepo.GetAll()
	for _, user := range users {
		if user.Email == req.Email {
			http.Error(w, "This email is already registered", http.StatusBadRequest)
			return
		}
	}
	user := &model.Todousers{
		Email:    req.Email,
		Password: req.Password,
	}
	if err := h.userRepo.Create(user); err != nil {
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
	userID, err := strconv.Atoi(params["users_id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := h.userRepo.DeleteUserTodos(userID); err != nil {
		http.Error(w, "Error deleting user todos", http.StatusInternalServerError)
		result := fmt.Sprintf(err.Error())
		fmt.Println(result)
		return
	}

	if err := h.userRepo.Delete(userID); err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	response := model.Response{
		Message: "User deleted successfully",
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
	user := &model.Todousers{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.userRepo.Update(id, user); err != nil {
		result := fmt.Sprintf(err.Error())
		fmt.Println(result)
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}
	response := model.Response{
		Message: "User updated successfully",
		Data:    user,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
