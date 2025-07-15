package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	url := "ws://localhost:8081/ws"

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			fmt.Printf("Gelen mesaj: %s\n", string(msg))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Mesaj: ")
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()

		err := conn.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
