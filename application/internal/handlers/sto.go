package handlers

import (
	"car-monitoring/internal/services"
	"encoding/json"
	"net/http"
	"strings"
)

// STOHandler handles service station API requests
type STOHandler struct {
	stoService *services.STOService
}

// NewSTOHandler creates a new STO handler
func NewSTOHandler(stoService *services.STOService) *STOHandler {
	return &STOHandler{
		stoService: stoService,
	}
}

// ServeHTTP handles GET /api/sto with optional query parameters
func (h *STOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	filters := make(map[string]string)
	query := r.URL.Query()

	for key, values := range query {
		if len(values) > 0 {
			// Handle special cases like "rating>=4"
			value := values[0]
			if strings.Contains(value, ">=") || strings.Contains(value, ">") {
				filters[key] = value
			} else {
				filters[key] = value
			}
		}
	}

	stations := h.stoService.GetServiceStations(filters)
	json.NewEncoder(w).Encode(stations)
}

