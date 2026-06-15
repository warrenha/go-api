package handlers

import (
	"encoding/json"
	"math"
	"net/http"

	"go-api/models"

	"github.com/shirou/gopsutil/v4/mem"
)

const GB = 1024 * 1024 * 1024

func round(value float64, places int) float64 {
	factor := math.Pow(10, float64(places))
	return math.Round(value*factor) / factor
}

// GetMemoryInfo fetches and formats virtual memory statistics.
func GetMemoryInfo() (*models.MemoryResponse, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	response := &models.MemoryResponse{
		Total:       memInfo.Total,
		TotalGB:     round(float64(memInfo.Total)/GB, 3),
		Available:   memInfo.Available,
		AvailableGB: round(float64(memInfo.Available)/GB, 3),
		Used:        memInfo.Used,
		UsedGB:      round(float64(memInfo.Used)/GB, 3),
		UsedPercent: memInfo.UsedPercent,
	}

	return response, nil
}

// Memory handles the GET /memory request.
// @Summary Get memory statistics
// @Description Returns details about the virtual memory usage of the system (total, used, available, etc.)
// @Produce json
// @Success 200 {object} models.MemoryResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /memory [get]
func Memory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response, err := GetMemoryInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}
