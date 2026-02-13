package repositories

import (
	"car-monitoring/internal/data"
	"car-monitoring/internal/models"
	"errors"
	"strings"
)

type VehicleRepository struct{}

func NewVehicleRepository() *VehicleRepository {
	return &VehicleRepository{}
}

func (r *VehicleRepository) GetAllVehicles() []models.Vehicle {
	return data.MockVehicles
}

func (r *VehicleRepository) GetVehicleByID(vehicleID string) (*models.Vehicle, error) {
	for i := range data.MockVehicles {
		if data.MockVehicles[i].VehicleID == vehicleID {
			return &data.MockVehicles[i], nil
		}
	}
	return nil, errors.New("vehicle not found")
}

func (r *VehicleRepository) GetVehicleByFlexibleID(id string) (*models.Vehicle, error) {
	for i := range data.MockVehicles {
		v := &data.MockVehicles[i]
		if strings.EqualFold(v.VehicleID, id) ||
			strings.EqualFold(v.Brand+v.Model, id) ||
			strings.EqualFold(v.VehicleID, strings.ToUpper(id)) {
			return v, nil
		}
	}
	return nil, errors.New("vehicle not found")
}

func (r *VehicleRepository) GetSelectedVehicle() *models.Vehicle {
	return data.GetSelectedVehicle()
}

func (r *VehicleRepository) SetSelectedVehicle(vehicleID string) bool {
	return data.SetSelectedVehicle(vehicleID)
}

func (r *VehicleRepository) SaveVehicles() error {
	return data.SaveMockVehicles()
}
