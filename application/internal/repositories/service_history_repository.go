package repositories

import (
	"car-monitoring/internal/data"
	"car-monitoring/internal/models"
)

type ServiceHistoryRepository struct{}

func NewServiceHistoryRepository() *ServiceHistoryRepository {
	return &ServiceHistoryRepository{}
}

func (r *ServiceHistoryRepository) Get(vehicleID string) []models.ServiceHistory {
	return data.GetServiceHistory(vehicleID)
}

func (r *ServiceHistoryRepository) Append(vehicleID string, rec models.ServiceHistory) error {
	return data.AppendServiceHistory(vehicleID, rec)
}
