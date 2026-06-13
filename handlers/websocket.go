package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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
// and echoing back any received messages wrapped in a JSON payload.
func Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		var payload interface{}
		if json.Valid(message) {
			payload = json.RawMessage(message)
		} else {
			payload = string(message)
		}

		response := map[string]interface{}{
			"payload": payload,
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			log.Printf("failed to marshal response: %v", err)
			continue
		}

		err = conn.WriteMessage(messageType, responseBytes)
		if err != nil {
			log.Printf("failed to write message: %v", err)
			break
		}
	}
}
