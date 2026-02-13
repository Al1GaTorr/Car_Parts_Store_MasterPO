package handlers

import (
	"car-monitoring/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

// RepairShopsHandler handles repair shop API requests with 2GIS integration
type RepairShopsHandler struct {
	sto2GISService  *services.STO2GISService
	vehicleService  *services.VehicleService
}

// NewRepairShopsHandler creates a new repair shops handler
func NewRepairShopsHandler(sto2GISService *services.STO2GISService, vehicleService *services.VehicleService) *RepairShopsHandler {
	return &RepairShopsHandler{
		sto2GISService: sto2GISService,
		vehicleService: vehicleService,
	}
}

// ServeHTTP handles GET /api/repair-shops
func (h *RepairShopsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	query := r.URL.Query()

	// Get selected vehicle to determine location and brand
	vehicle, err := h.vehicleService.GetSelectedVehicle()
	if err != nil {
		// Default location if no vehicle selected
		vehicle = nil
	}

	// Default location (Almaty, Kazakhstan) or use vehicle location
	lat := 43.2220
	lon := 76.8512
	radius := 10.0 // km

	// Location mapping for Kazakhstan cities
	locationMap := map[string]struct{ lat, lon float64 }{
		"Almaty":    {43.2220, 76.8512},
		"Astana":    {51.1694, 71.4491},
		"Shymkent":  {42.3419, 69.5901},
		"Karaganda": {49.8014, 73.1025},
		"Aktobe":    {50.2833, 57.1667},
		"Taraz":     {42.9000, 71.3667},
		"Kostanay":  {53.2144, 63.6246},
	}

	if vehicle != nil {
		if coords, ok := locationMap[vehicle.Location]; ok {
			lat = coords.lat
			lon = coords.lon
		}
	}

	// Override with query params if provided
	if latStr := query.Get("lat"); latStr != "" {
		if parsed, err := strconv.ParseFloat(latStr, 64); err == nil {
			lat = parsed
		}
	}
	if lonStr := query.Get("lon"); lonStr != "" {
		if parsed, err := strconv.ParseFloat(lonStr, 64); err == nil {
			lon = parsed
		}
	}
	if radiusStr := query.Get("radius"); radiusStr != "" {
		if parsed, err := strconv.ParseFloat(radiusStr, 64); err == nil {
			radius = parsed
		}
	}

	// Build filters
	filters := make(map[string]string)
	if rating := query.Get("rating"); rating != "" {
		filters["rating"] = rating
	}
	if price := query.Get("price"); price != "" {
		filters["price"] = price
	}
	
	// Filter by vehicle brand if available
	if vehicle != nil {
		brand := query.Get("brand")
		if brand == "" {
			brand = vehicle.Brand
		}
		if brand != "" {
			filters["brand"] = brand
		}
	}

	shops, err := h.sto2GISService.GetRepairShopsNearby(lat, lon, radius, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(shops)
}
