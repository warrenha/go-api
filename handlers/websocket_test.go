package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestWebsocket(t *testing.T) {
	// Create a test server with our Websocket handler.
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", Websocket)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// The server should send the initial memory update immediately
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read initial message: %v", err)
	}

	// Unmarshal and verify the response structure
	var response map[string]interface{}
	err = json.Unmarshal(msg, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	payloadVal, ok := response["payload"]
	if !ok {
		t.Fatal("Response does not contain 'payload' key")
	}

	// Convert payloadVal to a map to inspect memory fields
	payloadMap, ok := payloadVal.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected payload to be a JSON object, got %T", payloadVal)
	}

	// Verify required memory keys are present
	requiredKeys := []string{"total", "totalGb", "available", "availableGb", "used", "usedGb", "usedPercent"}
	for _, key := range requiredKeys {
		if _, exists := payloadMap[key]; !exists {
			t.Errorf("Expected memory payload to contain key %q", key)
		}
	}
}
