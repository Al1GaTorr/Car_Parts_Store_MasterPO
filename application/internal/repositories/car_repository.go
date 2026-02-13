package repositories

import (
	"car-monitoring/internal/data"
	"car-monitoring/internal/models"
	"errors"
)

type CarRepository struct{}

func NewCarRepository() *CarRepository {
	return &CarRepository{}
}

func (r *CarRepository) GetAllCars() []models.Car {
	return data.MockCars
}

func (r *CarRepository) GetCarByID(id string) (*models.Car, error) {
	for i := range data.MockCars {
		if data.MockCars[i].ID == id {
			return &data.MockCars[i], nil
		}
	}
	return nil, errors.New("car not found")
}

func (r *CarRepository) GetTelemetry(carID string) (*models.Telemetry, error) {
	telemetry, exists := data.MockTelemetry[carID]
	if !exists {
		return nil, errors.New("telemetry not found for car")
	}
	return &telemetry, nil
}

func (r *CarRepository) SetTelemetry(carID string, telemetry models.Telemetry) {
	if data.MockTelemetry == nil {
		data.MockTelemetry = make(map[string]models.Telemetry)
	}
	data.MockTelemetry[carID] = telemetry
}

func (r *CarRepository) GetErrorCodes(carID string) ([]models.ErrorCode, error) {
	codes, exists := data.MockErrorCodes[carID]
	if !exists {
		return []models.ErrorCode{}, nil
	}
	return codes, nil
}

func (r *CarRepository) GetMileageByCarID(carID string) int {
	for _, c := range data.MockCars {
		if c.ID == carID {
			return c.Mileage
		}
	}
	return 0
}
