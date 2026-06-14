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

// sendMemoryStats fetches, formats, and writes the current memory stats to the WebSocket connection.
func sendMemoryStats(conn *websocket.Conn) error {
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

// listenForDisconnect starts a background goroutine to read from the connection
// and returns a channel that will close when the client disconnects or an error occurs.
func listenForDisconnect(conn *websocket.Conn) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
	return done
}

// broadcastMemoryStats starts a ticker and sends memory updates to the client
// at the specified interval until the done channel is closed or a write error occurs.
func broadcastMemoryStats(conn *websocket.Conn, done <-chan struct{}, interval time.Duration) {
	// Send initial memory stats immediately on connection
	if err := sendMemoryStats(conn); err != nil {
		log.Printf("failed to send initial memory stats: %v", err)
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if err := sendMemoryStats(conn); err != nil {
				log.Printf("failed to send memory stats: %v", err)
				return
			}
		}
	}
}

// Websocket handles incoming WebSocket connections, upgrading the protocol,
// and periodically sending the system's memory info to the client every 5 seconds.
func Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade connection: %v", err)
		return
	}
	defer func() {
		conn.Close()
		log.Printf("WebSocket connection closed: %s", r.RemoteAddr)
	}()

	log.Printf("WebSocket connection established: %s", r.RemoteAddr)

	// Listen for client disconnection in the background
	done := listenForDisconnect(conn)

	// Broadcast memory stats to the client every 5 seconds until they disconnect
	broadcastMemoryStats(conn, done, 5*time.Second)
}
