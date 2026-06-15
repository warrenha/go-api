package models

// WebsocketResponse represents the JSON response broadcasted by the WebSocket.
type WebsocketResponse struct {
	Payload *MemoryResponse `json:"payload"`
}
