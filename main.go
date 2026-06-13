package main

import (
	"log"
	"net/http"

	"go-api/handlers"
)

/*
 *
 */
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello", handlers.Hello)
    mux.HandleFunc("GET /cpu", handlers.Cpu)
    mux.HandleFunc("GET /memory", handlers.Memory)
    mux.HandleFunc("GET /ws", handlers.Websocket)

	log.Println("Server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
