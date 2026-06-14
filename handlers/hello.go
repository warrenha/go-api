package handlers

import (
	"encoding/json"
	"net/http"
)

// HelloResponse represents the JSON response for the hello endpoint.
type HelloResponse struct {
	Message string `json:"message"`
}

// Hello handles the GET /hello request.
// @Summary Get a hello message
// @Description Returns a simple greeting message
// @Produce json
// @Success 200 {object} HelloResponse
// @Router /hello [get]
func Hello(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	response := HelloResponse{
		Message: "Hello World",
	}
	json.NewEncoder(w).Encode(response)
}
