package handlers

import (
	"car-monitoring/internal/services"
	"encoding/json"
	"net/http"
	"strings"
)

// AnalysisHandler handles analysis API requests
type AnalysisHandler struct {
	analysisService *services.AnalysisService
}

// NewAnalysisHandler creates a new analysis handler
func NewAnalysisHandler(analysisService *services.AnalysisService) *AnalysisHandler {
	return &AnalysisHandler{
		analysisService: analysisService,
	}
}

// ServeHTTP handles GET /api/cars/{id}/analysis
func (h *AnalysisHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract car ID from path: /api/cars/{id}/analysis
	path := strings.TrimPrefix(r.URL.Path, "/api/cars/")
	parts := strings.Split(path, "/")
	
	if len(parts) < 2 || parts[1] != "analysis" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	carID := parts[0]
	analysis, err := h.analysisService.GetAnalysis(carID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(analysis)
}

