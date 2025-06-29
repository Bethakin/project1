package handler

import (
	_ "encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Bethakin/project1/internal/database"
	"github.com/Bethakin/project1/internal/repository"
	"github.com/Bethakin/project1/model"
	_ "github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
)

type TodoHandler struct {
	todoRepo *repository.TodoRepository
}

func NewTodoHandler(db *database.Database) *TodoHandler {
	return &TodoHandler{
		todoRepo: repository.NewTodoRepository(db),
	}
}

/*
func (h *TodoHandler) IndexTodo(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	todos, err := h.todoRepo.GetByUserID(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("No todos found: %s", err.Error()), http.StatusNotFound)
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
}*/

func (h *TodoHandler) IndexTodo(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("users_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}
	todos, err := h.todoRepo.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": todos})
}

/*
func (h *TodoHandler) ShowTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	userIDParam := params["users_id"]
	userIDFromToken := r.Context().Value("userID").(float64)
	if userIDParam != fmt.Sprintf("%.0f", userIDFromToken) {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	todo, err := h.todoRepo.GetByID(id)
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
}*/

func (h *TodoHandler) ShowTodo(c echo.Context) error {
	todoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid todo ID"})
	}
	todo, err := h.todoRepo.GetByID(todoID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": todo})
}

/*
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	var req model.TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	todo := &model.Todo{
		Title:       req.Title,
		Description: req.Description,
		UsersID:     userID,
	}

	if err := h.todoRepo.Create(todo); err != nil {
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
}*/

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("users_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var todo model.Todo
	if err := c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	todo.UsersID = userID

	if err := h.todoRepo.Create(&todo); err != nil {
		fmt.Println("DB error during todo creation:", err)

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Todo creation failed"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Todo created",
		"data":    todo,
	})
}

/*
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	userIDParam := params["users_id"]
	userIDFromToken := r.Context().Value("userID").(float64)
	if userIDParam != fmt.Sprintf("%.0f", userIDFromToken) {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := h.todoRepo.Delete(id); err != nil {
		http.Error(w, "Error deleting todo", http.StatusInternalServerError)
		return
	}

	response := model.Response{
		Message: "Todo deleted successfully",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}*/

func (h *TodoHandler) Delete(c echo.Context) error {
	todoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid todo ID"})
	}
	if err := h.todoRepo.Delete(todoID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Todo deleted"})
}

func (h *TodoHandler) Update(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("users_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}
	todoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid todo ID"})
	}

	var todo model.Todo
	if err := c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.todoRepo.Update(todoID, userID, &todo); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Todo updated"})
}

/*
func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userIDParam := params["users_id"]
	userIDFromToken := r.Context().Value("userID").(float64)
	if userIDParam != fmt.Sprintf("%.0f", userIDFromToken) {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(params["users_id"])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
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

	if err := h.todoRepo.Update(id, userID, todo); err != nil {
		http.Error(w, "Error updating todo", http.StatusInternalServerError)
		return
	}

	response := model.Response{
		Message: "Todo updated successfully",
		Data:    todo,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}*/
