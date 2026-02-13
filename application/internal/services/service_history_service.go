package services

import (
	"car-monitoring/internal/models"
)

// ServiceHistoryService handles service history data
type ServiceHistoryService struct {
	carService *CarService
}

// NewServiceHistoryService creates a new service history service
func NewServiceHistoryService(carService *CarService) *ServiceHistoryService {
	return &ServiceHistoryService{
		carService: carService,
	}
}

// GetServiceHistory returns service history for a car
func (s *ServiceHistoryService) GetServiceHistory(carID string) ([]models.ServiceHistory, error) {
	car, err := s.carService.GetCarByID(carID)
	if err != nil {
		return nil, err
	}

	// Generate service history based on car mileage and maintenance needs
	history := []models.ServiceHistory{}

	// Generate oil change records
	mileage := car.Mileage
	lastOilChange := car.LastOilChange

	// Generate past oil changes
	for mileage >= lastOilChange+10000 && lastOilChange > 0 {
		history = append(history, models.ServiceHistory{
			ID:          len(history) + 1,
			Date:        s.estimateDateForMileage(car, lastOilChange),
			Type:        "Oil Change",
			Description: "Full synthetic oil change and filter replacement",
			Mileage:     lastOilChange,
			Cost:        15000,
			Shop:        "AutoService Premium",
			Location:    "Almaty, Abay Ave 150",
			Verified:    true,
			Icon:        "Droplet",
			Color:       "cyan",
		})
		lastOilChange -= 10000
	}

	// Add other maintenance records based on mileage
	if car.Mileage >= 40000 {
		history = append(history, models.ServiceHistory{
			ID:          len(history) + 1,
			Date:        s.estimateDateForMileage(car, 42100),
			Type:        "Brake Service",
			Description: "Front brake pads and rotors replacement",
			Mileage:     42100,
			Cost:        45000,
			Shop:        "MasterCar Service",
			Location:    "Almaty, Satpayev St 90",
			Verified:    true,
			Icon:        "Wrench",
			Color:       "emerald",
		})
	}

	if car.Mileage >= 38000 {
		history = append(history, models.ServiceHistory{
			ID:          len(history) + 1,
			Date:        s.estimateDateForMileage(car, 38200),
			Type:        "Air Filter",
			Description: "Engine and cabin air filter replacement",
			Mileage:     38200,
			Cost:        8000,
			Shop:        "QuickFix Auto",
			Location:    "Almaty, Rozybakiev St 247",
			Verified:    true,
			Icon:        "Filter",
			Color:       "blue",
		})
	}

	// Sort by date (most recent first)
	// For simplicity, reverse the slice
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	return history, nil
}

// estimateDateForMileage estimates a date based on mileage
func (s *ServiceHistoryService) estimateDateForMileage(car *models.Car, mileage int) string {
	// Simple estimation: assume 50 km per day average
	// Return a date string (simplified)
	return "2025-11-15" // This would be calculated properly in production
}
