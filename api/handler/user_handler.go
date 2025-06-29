package handler

import (
	"encoding/json"
	_ "fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Bethakin/project1/internal/database"
	"github.com/Bethakin/project1/internal/repository"
	utils "github.com/Bethakin/project1/jwt"
	"github.com/Bethakin/project1/model"
	_ "github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(db *database.Database) *UserHandler {
	return &UserHandler{
		userRepo: repository.NewUserRepository(db),
	}
}

/*
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userInput model.Todousers
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	users, err := h.userRepo.GetAll()
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	var foundUser *model.Todousers
	for _, u := range users {
		if u.Email == userInput.Email && u.Password == userInput.Password {
			foundUser = u
			break
		}
	}
	if foundUser == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		http.Error(w, "JWT secret not set", http.StatusInternalServerError)
		return
	}
	token, err := utils.GenerateJWT(jwtSecret, foundUser.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, token)))
}*/

func (h *UserHandler) LoginUser(c echo.Context) error {
	var user model.Todousers
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	existing, err := h.userRepo.GetByEmail(user.Email)
	if err != nil || existing.Password != user.Password {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	token, err := utils.GenerateJWT(os.Getenv("JWT_SECRET"), existing.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *UserHandler) Index(c echo.Context) error {
	idParam := c.Param("users_id")

	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": user})
}

/*
func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userIDParam := vars["users_id"]
	userIDFromToken := r.Context().Value("userID").(float64)
	if userIDParam != fmt.Sprintf("%.0f", userIDFromToken) {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.Atoi(vars["users_id"])

	users, err := h.userRepo.GetByID(userID)
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
}*/

func (h *UserHandler) IndexAll(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) CreateUser(c echo.Context) error {
	var user model.Todousers
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.userRepo.Create(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User creation failed"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
		"data":    user,
	})
}

/*
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
	userIDParam := params["users_id"]
	userIDFromToken := r.Context().Value("userID").(float64)

	if userIDParam != fmt.Sprintf("%.0f", userIDFromToken) {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}
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
	userIDParam := params["users_id"]
	userIDFromToken := r.Context().Value("userID").(float64)

	if userIDParam != fmt.Sprintf("%.0f", userIDFromToken) {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}
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
}*/

func (h *UserHandler) Updateusers(c echo.Context) error {
	idParam := c.Param("users_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var user model.Todousers
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.userRepo.Update(id, &user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Update failed"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "User updated"})
}

func (h *UserHandler) Deleteusers(c echo.Context) error {
	idParam := c.Param("users_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	if err := h.userRepo.DeleteUserTodos(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error deleting todos"})
	}
	if err := h.userRepo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error deleting user"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted"})
}
