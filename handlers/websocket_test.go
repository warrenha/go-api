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

	// Test case 1: Valid JSON payload
	testPayload := `{"hello":"world","count":42}`
	err = conn.WriteMessage(websocket.TextMessage, []byte(testPayload))
	if err != nil {
		t.Fatalf("Failed to write message: %v", err)
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	// Verify the response contains the original JSON as a child 'payload' object
	var response map[string]interface{}
	err = json.Unmarshal(msg, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	payloadVal, ok := response["payload"]
	if !ok {
		t.Fatal("Response does not contain 'payload' key")
	}

	// Convert payloadVal back to JSON to check its structure
	payloadJSON, err := json.Marshal(payloadVal)
	if err != nil {
		t.Fatalf("Failed to marshal payload back to JSON: %v", err)
	}

	// Compare raw JSON values (ignoring space differences)
	var originalMap, returnedMap map[string]interface{}
	if err := json.Unmarshal([]byte(testPayload), &originalMap); err != nil {
		t.Fatalf("Failed to unmarshal original: %v", err)
	}
	if err := json.Unmarshal(payloadJSON, &returnedMap); err != nil {
		t.Fatalf("Failed to unmarshal returned payload: %v", err)
	}

	if len(originalMap) != len(returnedMap) || returnedMap["hello"] != "world" || returnedMap["count"].(float64) != 42 {
		t.Errorf("Expected payload to match original map, got %v", returnedMap)
	}

	// Test case 2: Non-JSON payload (should be handled gracefully as string)
	testNonJSON := `hello world`
	err = conn.WriteMessage(websocket.TextMessage, []byte(testNonJSON))
	if err != nil {
		t.Fatalf("Failed to write message: %v", err)
	}

	_, msg, err = conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	var response2 map[string]interface{}
	err = json.Unmarshal(msg, &response2)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	payloadVal2, ok := response2["payload"]
	if !ok {
		t.Fatal("Response does not contain 'payload' key")
	}

	if payloadVal2 != "hello world" {
		t.Errorf("Expected payload to be string 'hello world', got %v", payloadVal2)
	}
}
