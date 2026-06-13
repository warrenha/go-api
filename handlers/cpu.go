package handlers

import (
	"encoding/json"
	"net/http"
	"time"

    "github.com/shirou/gopsutil/v4/cpu"
)

type CpuResponse struct {
	Usage float64 `json:"usage"`
}

func Cpu(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

    percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := CpuResponse{
		Usage: percent[0],
	}
	json.NewEncoder(w).Encode(response)
}
