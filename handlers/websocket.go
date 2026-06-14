package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// pingInterval is how often we send a ping to the client.
	pingInterval = 30 * time.Second
	// pongTimeout is how long we wait for a pong before considering the client dead.
	// Must be greater than pingInterval.
	pongTimeout = 60 * time.Second
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

// listenForDisconnect starts a background goroutine to read from the connection.
// It sets up ping/pong handlers to detect dead connections and returns a channel
// that closes when the client disconnects or stops responding to pings.
func listenForDisconnect(conn *websocket.Conn) <-chan struct{} {
	done := make(chan struct{})

	// Set an initial read deadline. This will be extended each time a pong is received.
	conn.SetReadDeadline(time.Now().Add(pongTimeout))

	// When a pong is received, reset the read deadline, keeping the connection alive.
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongTimeout))
		return nil
	})

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

// startPingLoop sends periodic pings to the client to detect dead connections.
// It stops when the done channel is closed.
func startPingLoop(conn *websocket.Conn, done <-chan struct{}) {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("failed to send ping: %v", err)
				return
			}
		}
	}
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

	// Send periodic pings to detect dead connections
	go startPingLoop(conn, done)

	// Broadcast memory stats to the client every 5 seconds until they disconnect
	broadcastMemoryStats(conn, done, 5*time.Second)
}
