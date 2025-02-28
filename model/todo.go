package model

type TodoRequest struct {
	Title       string `json:"title"` //??
	Description string `json:"description"`
}

type TodoRequestUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoUser struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`  //??
	Error   string      `json:"error,omitempty"` //??
}
