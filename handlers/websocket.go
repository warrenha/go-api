package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin allows all connections by default for testing purposes.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Websocket handles incoming WebSocket connections, upgrading the protocol,
// and periodically sending the system's memory info to the client every 5 seconds.
func Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Done channel to signal when the client has disconnected
	done := make(chan struct{})

	// Run read loop in a background goroutine to detect disconnects
	go func() {
		defer close(done)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Helper to send memory stats
	sendMemoryStats := func() error {
		memInfo, err := GetMemoryInfo()
		if err != nil {
			return err
		}

		response := map[string]interface{}{
			"payload": memInfo,
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			return err
		}

		return conn.WriteMessage(websocket.TextMessage, responseBytes)
	}

	// Send initial memory stats immediately on connection
	if err := sendMemoryStats(); err != nil {
		log.Printf("failed to send initial memory stats: %v", err)
		return
	}

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if err := sendMemoryStats(); err != nil {
				log.Printf("failed to send memory stats: %v", err)
				return
			}
		}
	}
}
