package main

import (
	"fmt"
	"net/http"

	"github.com/Bethakin/project1/project1/api/handler"
)

func main() {
	http.HandleFunc("/", handler.Welcome)
	fmt.Println("Server starting on port :8081...")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
		return
	}
}
