package ws

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func StartClient(name string) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8081/ws", nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	// Receive messages
	go func() {
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			fmt.Printf("[%s]: %s\n", msg.From, msg.Body)
		}
	}()

	// Send messages
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Message: ")
		text, _ := reader.ReadString('\n')
		conn.WriteJSON(Message{From: name, Body: text})
	}
}
