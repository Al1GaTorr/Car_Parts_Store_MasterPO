package handlers

import (
	"car-monitoring/internal/services"
	"encoding/json"
	"net/http"
	"strings"
)

// DashboardHandler handles dashboard API requests
type DashboardHandler struct {
	dashboardService *services.DashboardService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// ServeHTTP handles GET /api/dashboard/{carId}
func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract car ID from path (optional, uses selected vehicle if not provided)
	path := strings.TrimPrefix(r.URL.Path, "/api/dashboard")
	path = strings.TrimPrefix(path, "/")
	carID := path

	if carID == "" {
		// Use selected vehicle (default)
		carID = "selected"
	}

	dashboard, err := h.dashboardService.GetDashboardData(carID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(dashboard)
}

