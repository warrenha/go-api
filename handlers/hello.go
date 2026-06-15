package handlers

import (
	"encoding/json"
	"net/http"

	"go-api/models"
)

// Hello handles the GET /hello request.
// @Summary Get a hello message
// @Description Returns a simple greeting message
// @Produce json
// @Success 200 {object} models.HelloResponse
// @Router /hello [get]
func Hello(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	response := models.HelloResponse{
		Message: "Hello World",
	}
	json.NewEncoder(w).Encode(response)
}
