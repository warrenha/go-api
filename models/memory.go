package models

// MemoryResponse represents the JSON response for the memory statistics endpoint.
type MemoryResponse struct {
	Total       uint64  `json:"total"`
	TotalGB     float64 `json:"totalGb"`
	Available   uint64  `json:"available"`
	AvailableGB float64 `json:"availableGb"`
	Used        uint64  `json:"used"`
	UsedGB      float64 `json:"usedGb"`
	UsedPercent float64 `json:"usedPercent"`
}
