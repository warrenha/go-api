package handlers

import (
	"encoding/json"
	"net/http"
    "math"

	"github.com/shirou/gopsutil/v4/mem"
)

const GB = 1024 * 1024 * 1024

func round(value float64, places int) float64 {
	factor := math.Pow(10, float64(places))
	return math.Round(value*factor) / factor
}

type MemoryResponse struct {
	Total       uint64  `json:"total"`
	TotalGB     float64 `json:"totalGb"`
	Available   uint64  `json:"available"`
	AvailableGB float64 `json:"availableGb"`
	Used        uint64  `json:"used"`
	UsedGB      float64 `json:"usedGb"`
	UsedPercent float64 `json:"usedPercent"`
}

func Memory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := MemoryResponse{
		Total:       memInfo.Total,
        TotalGB:     round(float64(memInfo.Total)/GB, 3),
		Available:   memInfo.Available,
        AvailableGB: round(float64(memInfo.Available)/GB, 3),
		Used:        memInfo.Used,
        UsedGB:      round(float64(memInfo.Used)/GB, 3),
		UsedPercent: memInfo.UsedPercent,
	}

	json.NewEncoder(w).Encode(response)
}
