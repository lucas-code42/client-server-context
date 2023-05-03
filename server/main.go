package main

import (
	"fmt"
	"net/http"

	"github.com/desafio/sever/api"
)

var (
	PORT = 8080
)

func main() {
	fmt.Println("Server is running...")
	http.HandleFunc("/cotacao", api.Handler)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
