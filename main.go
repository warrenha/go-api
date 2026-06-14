package main

import (
	"log"
	"net/http"

	"go-api/handlers"
)

//go:generate go run github.com/swaggo/swag/cmd/swag@latest init

// @title Go System Stats API
// @version 1.0
// @description A simple Go API that provides CPU and Memory statistics, and a WebSocket stream.
// @host localhost:8081
// @BasePath /
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello", handlers.Hello)
    mux.HandleFunc("GET /cpu", handlers.Cpu)
    mux.HandleFunc("GET /memory", handlers.Memory)
    mux.HandleFunc("GET /ws", handlers.Websocket)

	log.Println("Server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
