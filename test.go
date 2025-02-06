package main

import (
	"fmt"
	"net/http"
)

func main() {
	response, error := http.Get("https://reqres.in/api/products")
	if error != nil {
		fmt.Println(error)
	}
//??????????
	fmt.Println(response)
}
