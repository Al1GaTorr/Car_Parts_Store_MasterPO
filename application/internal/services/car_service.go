package services

import (
	"car-monitoring/internal/models"
	"car-monitoring/internal/repositories"
)

// CarService handles business logic for cars
type CarService struct {
	carRepo *repositories.CarRepository
}

// NewCarService creates a new car service instance
func NewCarService() *CarService {
	return &CarService{
		carRepo: repositories.NewCarRepository(),
	}
}

// GetAllCars returns all available cars
func (s *CarService) GetAllCars() []models.Car {
	return s.carRepo.GetAllCars()
}

// GetCarByID returns a car by its ID
func (s *CarService) GetCarByID(id string) (*models.Car, error) {
	return s.carRepo.GetCarByID(id)
}

// GetTelemetry returns current telemetry for a car
func (s *CarService) GetTelemetry(carID string) (*models.Telemetry, error) {
	return s.carRepo.GetTelemetry(carID)
}

// SetTelemetry updates telemetry for a car
func (s *CarService) SetTelemetry(carID string, telemetry models.Telemetry) {
	s.carRepo.SetTelemetry(carID, telemetry)
}

// GetErrorCodes returns all DTC codes for a car
func (s *CarService) GetErrorCodes(carID string) ([]models.ErrorCode, error) {
	return s.carRepo.GetErrorCodes(carID)
}

// GetMileageByCarID returns mileage for car ID
func (s *CarService) GetMileageByCarID(carID string) int {
	return s.carRepo.GetMileageByCarID(carID)
}
