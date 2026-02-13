package handlers

import (
	"car-monitoring/internal/services"
	"encoding/json"
	"net/http"
	"strings"
)

// AdminHandler handles admin panel API requests
type AdminHandler struct {
	vehicleService *services.VehicleService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(vehicleService *services.VehicleService) *AdminHandler {
	return &AdminHandler{
		vehicleService: vehicleService,
	}
}

// ServeHTTP handles admin API requests
func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/admin/")

	switch {
	case path == "vehicles" && r.Method == http.MethodGet:
		// GET /api/admin/vehicles - Get all vehicles
		vehicles := h.vehicleService.GetAllVehicles()
		json.NewEncoder(w).Encode(vehicles)

	case path == "selected" && r.Method == http.MethodGet:
		// GET /api/admin/selected - Get selected vehicle
		vehicle, err := h.vehicleService.GetSelectedVehicle()
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(vehicle)

	case path == "selected" && r.Method == http.MethodPost:
		// POST /api/admin/selected - Set selected vehicle
		var req struct {
			VehicleID string `json:"vehicle_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := h.vehicleService.SetSelectedVehicle(req.VehicleID); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		vehicle, _ := h.vehicleService.GetSelectedVehicle()
		json.NewEncoder(w).Encode(vehicle)

	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

