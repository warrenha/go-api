package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go-api/models"

	"github.com/shirou/gopsutil/v4/cpu"
)

// Cpu handles the GET /cpu request.
// @Summary Get CPU usage
// @Description Returns the current CPU usage percentage measured over a 1-second interval
// @Produce json
// @Success 200 {object} models.CpuResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /cpu [get]
func Cpu(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.CpuResponse{
		Usage: percent[0],
	}
	json.NewEncoder(w).Encode(response)
}
