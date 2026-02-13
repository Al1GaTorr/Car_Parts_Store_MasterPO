package services

import (
	"car-monitoring/internal/data"
	"car-monitoring/internal/models"
	"car-monitoring/internal/repositories"
	"errors"
)

// VehicleService handles vehicle operations
type VehicleService struct {
	vehicleRepo *repositories.VehicleRepository
}

// NewVehicleService creates a new vehicle service
func NewVehicleService() *VehicleService {
	return &VehicleService{
		vehicleRepo: repositories.NewVehicleRepository(),
	}
}

// GetAllVehicles returns all available vehicles
func (s *VehicleService) GetAllVehicles() []models.Vehicle {
	return s.vehicleRepo.GetAllVehicles()
}

// GetVehicleByID returns a vehicle by its ID
func (s *VehicleService) GetVehicleByID(vehicleID string) (*models.Vehicle, error) {
	return s.vehicleRepo.GetVehicleByID(vehicleID)
}

// GetVehicleByFlexibleID returns vehicle by multiple matching rules
func (s *VehicleService) GetVehicleByFlexibleID(vehicleID string) (*models.Vehicle, error) {
	return s.vehicleRepo.GetVehicleByFlexibleID(vehicleID)
}

// GetSelectedVehicle returns the currently selected vehicle
func (s *VehicleService) GetSelectedVehicle() (*models.Vehicle, error) {
	vehicle := s.vehicleRepo.GetSelectedVehicle()
	if vehicle == nil {
		return nil, errors.New("no vehicle selected")
	}
	return vehicle, nil
}

// SetSelectedVehicle sets the currently selected vehicle
func (s *VehicleService) SetSelectedVehicle(vehicleID string) error {
	if !s.vehicleRepo.SetSelectedVehicle(vehicleID) {
		return errors.New("vehicle not found")
	}
	return nil
}

// SaveVehicles persists vehicle list
func (s *VehicleService) SaveVehicles() error {
	return s.vehicleRepo.SaveVehicles()
}

// ConvertVehicleToCar converts a Vehicle to Car model for compatibility
func (s *VehicleService) ConvertVehicleToCar(vehicle *models.Vehicle) *models.Car {
	// Calculate last oil change (assume every 10,000 km)
	lastOilChange := vehicle.LastServiceKm
	if vehicle.MileageKm-vehicle.LastServiceKm > 10000 {
		// Oil change overdue
		lastOilChange = vehicle.MileageKm - 10000
	}

	return &models.Car{
		ID:            vehicle.VehicleID,
		VIN:           vehicle.VehicleID,
		Brand:         vehicle.Brand,
		Model:         vehicle.Model,
		Year:          vehicle.Year,
		EngineType:    vehicle.EngineType,
		LicensePlate:  "E 777 KZ", // Default, can be customized
		Mileage:       vehicle.MileageKm,
		LastOilChange: lastOilChange,
	}
}

// GetVehicleTelemetry generates telemetry for a vehicle
func (s *VehicleService) GetVehicleTelemetry(vehicle *models.Vehicle) *models.Telemetry {
	// Generate realistic telemetry based on vehicle state
	telemetry := &models.Telemetry{
		CarID:                  vehicle.VehicleID,
		Speed:                  65.5,
		RPM:                    2200,
		EngineLoad:             45.2,
		ThrottlePosition:       35.0,
		IntakeAirTemp:          28.5,
		CoolantTemp:            92.0,
		FuelLevel:              45.0,
		AverageFuelConsumption: 8.5,
		InstantFuelConsumption: 7.8,
		AverageSpeed:           58.3,
		TripDistance:           125.5,
		FuelSystemStatus:       "closed loop",
		O2SensorVoltages:       []float64{0.45, 0.48, 0.44, 0.46},
		ShortTermFuelTrim:      2.5,
		LongTermFuelTrim:       1.2,
		IntakeManifoldPressure: 45.2,
		BatteryVoltage:         12.6,
		ABSSensorStatus:        "ok",
		TirePressure:           []float64{32, 33, 31, 32},
		BrakeSystemStatus:      "ok",
		AirbagSystemStatus:     "ok",
		EmissionsMonitorReady:  true,
		FreezeFrameData:        "",
		Timestamp:              0,
	}

	// Adjust based on errors
	for _, err := range vehicle.Errors {
		switch err.Code {
		case "P0300", "P0200":
			// Engine problems
			telemetry.CoolantTemp = 108.5
			telemetry.EngineLoad = 95.0
		case "P0562":
			// Low voltage
			telemetry.BatteryVoltage = 11.5
		case "C1234", "C0035":
			// ABS problems
			telemetry.ABSSensorStatus = "error"
			telemetry.BrakeSystemStatus = "error"
		case "B1000", "B0020":
			// Airbag problems
			telemetry.AirbagSystemStatus = "error"
		}
	}

	return telemetry
}

// GetVehicleErrorCodes converts vehicle errors to ErrorCode models
func (s *VehicleService) GetVehicleErrorCodes(vehicle *models.Vehicle) []models.ErrorCode {
	errorCodes := []models.ErrorCode{}

	for _, ve := range vehicle.Errors {
		severity, description := data.GetErrorSeverityDescription(ve.Code)

		// Use description from severity map if available, otherwise use vehicle description
		if description == "Неизвестная ошибка" {
			description = ve.Description
		}

		errorCodes = append(errorCodes, models.ErrorCode{
			Code:        ve.Code,
			Description: description,
			Criticality: severity,
			MILStatus:   severity == "critical" || severity == "medium",
		})
	}

	return errorCodes
}

// GetErrorSeverity returns severity and description for an error code
func (s *VehicleService) GetErrorSeverity(code string) (string, string) {
	return data.GetErrorSeverityDescription(code)
}
