package services

import (
	"car-monitoring/internal/data"
	"car-monitoring/internal/models"
	"fmt"
	"time"
)

// NotificationService handles notifications and alerts
type NotificationService struct {
	vehicleService *VehicleService
}

// NewNotificationService creates a new notification service
func NewNotificationService(vehicleService *VehicleService, service *AnalysisService) *NotificationService {
	return &NotificationService{
		vehicleService: vehicleService,
	}
}

// GetNotifications returns all notifications for the selected vehicle
func (s *NotificationService) GetNotifications(carID string) ([]models.Notification, error) {
	// Get selected vehicle
	vehicle, err := s.vehicleService.GetSelectedVehicle()
	if err != nil {
		return nil, err
	}

	notifications := []models.Notification{}

	// Convert vehicle errors to notifications with severity-based messages
	for _, ve := range vehicle.Errors {
		severity, description := data.GetErrorSeverityDescription(ve.Code)

		notificationType := "warning"
		title := fmt.Sprintf("Проблема: %s", ve.Code)

		if severity == "critical" {
			notificationType = "critical"
			title = fmt.Sprintf("КРИТИЧЕСКАЯ ПРОБЛЕМА: %s", ve.Code)
			description = fmt.Sprintf("⚠️ СРОЧНО! %s. Рекомендуется немедленно обратиться в сервис.", description)
		} else if severity == "medium" {
			notificationType = "warning"
			title = fmt.Sprintf("Проблема: %s", ve.Code)
			description = fmt.Sprintf("⚠️ %s. Рекомендуется диагностика.", description)
		} else {
			notificationType = "info"
			title = fmt.Sprintf("Информация: %s", ve.Code)
		}

		notifications = append(notifications, models.Notification{
			ID:         fmt.Sprintf("error-%s-%s", vehicle.VehicleID, ve.Code),
			Type:       notificationType,
			Title:      title,
			Message:    description,
			Timestamp:  time.Now().Unix(),
			Read:       false,
			Actionable: true,
		})
	}

	// Add maintenance notifications
	kmSinceService := vehicle.MileageKm - vehicle.LastServiceKm
	if kmSinceService > 10000 {
		notifications = append(notifications, models.Notification{
			ID:         fmt.Sprintf("maintenance-oil-%s", vehicle.VehicleID),
			Type:       "warning",
			Title:      "Замена масла просрочена",
			Message:    fmt.Sprintf("Замена масла просрочена на %d км. Срочно запланируйте обслуживание.", kmSinceService-10000),
			Timestamp:  time.Now().Unix(),
			Read:       false,
			Actionable: true,
		})
	}

	// Add maintenance alerts from vehicle
	for _, maintAlert := range vehicle.MaintenanceAlerts {
		notifications = append(notifications, models.Notification{
			ID:         fmt.Sprintf("maintenance-%s-%d", vehicle.VehicleID, time.Now().Unix()),
			Type:       "warning",
			Title:      maintAlert,
			Message:    maintAlert,
			Timestamp:  time.Now().Unix(),
			Read:       false,
			Actionable: true,
		})
	}

	return notifications, nil
}
