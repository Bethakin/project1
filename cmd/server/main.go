package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Bethakin/project1/api/handler"
	"github.com/Bethakin/project1/internal/database"
	utils "github.com/Bethakin/project1/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebSocketHub struct {
	clients map[*websocket.Conn]bool
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (hub *WebSocketHub) broadcastMessage(sender *websocket.Conn, message []byte) {
	for client := range hub.clients {
		if client != sender {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Write error:", err)
				client.Close()
				delete(hub.clients, client)
			}
		}
	}
}

func (hub *WebSocketHub) wsHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return err
	}
	defer conn.Close()

	hub.clients[conn] = true
	log.Println("New WebSocket client connected")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			delete(hub.clients, conn)
			break
		}

		log.Printf("Received message: %s\n", string(msg))
		hub.broadcastMessage(conn, msg)
	}

	return nil
}

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if db == nil {
		log.Fatal("Database is nil after initialization!")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in the environment")
	}

	e := echo.New()

	userHandler := handler.NewUserHandler(db)
	todoHandler := handler.NewTodoHandler(db)

	e.POST("/register", userHandler.CreateUser)
	e.POST("/login", userHandler.LoginUser)

	protected := e.Group("/users/:users_id", utils.AuthMiddleware(jwtSecret))

	protected.GET("", userHandler.Index)
	protected.PUT("", userHandler.Updateusers)
	protected.DELETE("", userHandler.Deleteusers)

	todo := protected.Group("/todos")
	todo.GET("", todoHandler.IndexTodo)
	todo.POST("", todoHandler.CreateTodo)
	todo.GET("/:id", todoHandler.ShowTodo)
	todo.PUT("/:id", todoHandler.Update)
	todo.DELETE("/:id", todoHandler.Delete)

	//WS
	hub := NewWebSocketHub()

	e.GET("/ws", hub.wsHandler)

	log.Println("Server starting on port 8081...")
	e.Logger.Fatal(e.Start(":8081"))
}
