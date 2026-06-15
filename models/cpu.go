package models

// CpuResponse represents the JSON response for the CPU usage endpoint.
type CpuResponse struct {
	Usage float64 `json:"usage"`
}
