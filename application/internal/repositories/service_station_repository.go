package repositories

import (
	"car-monitoring/internal/data"
	"car-monitoring/internal/models"
)

type ServiceStationRepository struct{}

func NewServiceStationRepository() *ServiceStationRepository {
	return &ServiceStationRepository{}
}

func (r *ServiceStationRepository) GetAll() []models.ServiceStation {
	return data.MockServiceStations
}
