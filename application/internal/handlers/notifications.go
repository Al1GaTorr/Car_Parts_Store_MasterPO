package handlers

import (
	"car-monitoring/internal/services"
	"encoding/json"
	"net/http"
	"strings"
)

// NotificationHandler handles notification API requests
type NotificationHandler struct {
	notificationService *services.NotificationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// ServeHTTP handles GET /api/notifications/{carId}
func (h *NotificationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract car ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/notifications/")
	carID := path

	if carID == "" {
		carID = "car1" // Default
	}

	notifications, err := h.notificationService.GetNotifications(carID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(notifications)
}

