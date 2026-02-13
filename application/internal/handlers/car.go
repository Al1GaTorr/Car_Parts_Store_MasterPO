package handlers

import (
	"car-monitoring/internal/services"
	"encoding/json"
	"net/http"
	"strings"
)

// CarHandler handles car-related API requests
type CarHandler struct {
	carService *services.CarService
}

// NewCarHandler creates a new car handler
func NewCarHandler(carService *services.CarService) *CarHandler {
	return &CarHandler{
		carService: carService,
	}
}

// ServeHTTP handles car API requests
func (h *CarHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	path := strings.TrimPrefix(r.URL.Path, "/api/cars/")
	
	if path == "" || path == "/api/cars" {
		// GET /api/cars - list all cars
		if r.Method == http.MethodGet {
			cars := h.carService.GetAllCars()
			json.NewEncoder(w).Encode(cars)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if it's a specific car endpoint
	parts := strings.Split(path, "/")
	carID := parts[0]

	if len(parts) == 1 {
		// GET /api/cars/{id} - get specific car
		if r.Method == http.MethodGet {
			car, err := h.carService.GetCarByID(carID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(car)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// GET /api/cars/{id}/telemetry
	if len(parts) == 2 && parts[1] == "telemetry" {
		if r.Method == http.MethodGet {
			telemetry, err := h.carService.GetTelemetry(carID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(telemetry)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// POST /api/cars/{id}/records - add a service record and broadcast
	if len(parts) == 2 && parts[1] == "records" {
		if r.Method == http.MethodPost {
			var rec struct {
				Date        string `json:"date"`
				Mileage     int    `json:"mileage"`
				Description string `json:"description"`
				ServiceName string `json:"serviceName"`
				Type        string `json:"type"`
			}
			if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
				http.Error(w, "invalid payload", http.StatusBadRequest)
				return
			}

			// Acknowledge save (no persistent change to DB per constraints)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

			// Broadcast event to subscribers for this car ID
			payload := map[string]interface{}{
				"date":        rec.Date,
				"mileage":     rec.Mileage,
				"description": rec.Description,
				"serviceName": rec.ServiceName,
			}
			BroadcastToVIN(carID, "SERVICE_RECORD_ADDED", payload)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

