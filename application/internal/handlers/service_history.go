package handlers

import (
	"car-monitoring/internal/data"
	"car-monitoring/internal/services"
	"encoding/json"
	"net/http"
	"strings"
)

// ServiceHistoryHandler handles service history API requests
type ServiceHistoryHandler struct {
	historyService *services.ServiceHistoryService
}

// NewServiceHistoryHandler creates a new service history handler
func NewServiceHistoryHandler(historyService *services.ServiceHistoryService) *ServiceHistoryHandler {
	return &ServiceHistoryHandler{
		historyService: historyService,
	}
}

// ServeHTTP handles GET /api/service-history/{carId}
func (h *ServiceHistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract car ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/service-history/")
	carID := path

	if carID == "" {
		carID = "car1" // Default
	}

	// First check persisted service history (created by STO panel)
	persisted := data.GetServiceHistory(carID)
	if len(persisted) > 0 {
		json.NewEncoder(w).Encode(persisted)
		return
	}

	// Fallback to generated history
	history, err := h.historyService.GetServiceHistory(carID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(history)
}

