package model

type TodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoRequestusers struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UsersID     int    `json:"users_id"`
}

type Todousers struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
