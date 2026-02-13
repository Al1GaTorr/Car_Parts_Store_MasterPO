package services

import (
	"car-monitoring/internal/models"
	"fmt"
)

// AnalysisService handles maintenance recommendations and alerts
type AnalysisService struct {
	carService *CarService
}

// NewAnalysisService creates a new analysis service instance
func NewAnalysisService(carService *CarService) *AnalysisService {
	return &AnalysisService{
		carService: carService,
	}
}

// GetAnalysis generates maintenance recommendations and alerts for a car
func (s *AnalysisService) GetAnalysis(carID string) (*models.Analysis, error) {
	car, err := s.carService.GetCarByID(carID)
	if err != nil {
		return nil, err
	}

	telemetry, err := s.carService.GetTelemetry(carID)
	if err != nil {
		return nil, err
	}

	errorCodes, err := s.carService.GetErrorCodes(carID)
	if err != nil {
		return nil, err
	}

	analysis := &models.Analysis{
		CarID:         carID,
		OilChangeDue:  s.checkOilChange(car),
		SparkPlugsDue: s.checkSparkPlugs(car),
		CoolantDue:    s.checkCoolant(car),
		BrakeFluidDue: s.checkBrakeFluid(car),
		AirFilterDue:  s.checkAirFilter(car),
		Errors:        errorCodes,
		Alerts:        s.generateAlerts(telemetry, errorCodes),
		Suggestions:   []string{},
	}

	// Generate suggestions based on analysis
	s.generateSuggestions(analysis, car, telemetry)

	return analysis, nil
}

// checkOilChange determines if oil change is due
// Oil change needed if: mileage > 10,000 km since last change OR oilRemainingKm < 500
func (s *AnalysisService) checkOilChange(car *models.Car) bool {
	kmSinceLastChange := car.Mileage - car.LastOilChange
	return kmSinceLastChange > 10000
}

// checkSparkPlugs determines if spark plugs replacement is due (every 30,000 km)
func (s *AnalysisService) checkSparkPlugs(car *models.Car) bool {
	// Assuming last replacement at 0 km, check if 30k interval passed
	return (car.Mileage%30000) < 5000 && car.Mileage > 30000
}

// checkCoolant determines if coolant change is due (every 60,000 km)
func (s *AnalysisService) checkCoolant(car *models.Car) bool {
	return (car.Mileage%60000) < 5000 && car.Mileage > 60000
}

// checkBrakeFluid determines if brake fluid change is due (every 40,000 km)
func (s *AnalysisService) checkBrakeFluid(car *models.Car) bool {
	return (car.Mileage%40000) < 5000 && car.Mileage > 40000
}

// checkAirFilter determines if air filter replacement is due (every 20,000 km)
func (s *AnalysisService) checkAirFilter(car *models.Car) bool {
	return (car.Mileage%20000) < 5000 && car.Mileage > 20000
}

// generateAlerts creates alerts based on telemetry thresholds
func (s *AnalysisService) generateAlerts(telemetry *models.Telemetry, errorCodes []models.ErrorCode) []models.Alert {
	alerts := []models.Alert{}

	// High engine temperature (>105°C)
	if telemetry.CoolantTemp > 105.0 {
		alerts = append(alerts, models.Alert{
			Type:        "critical",
			Message:     "High Engine Temperature",
			Description: fmt.Sprintf("Engine coolant temperature is %.1f°C. This is above the safe operating limit.", telemetry.CoolantTemp),
		})
	}

	// Low battery voltage (<12.0V)
	if telemetry.BatteryVoltage < 12.0 {
		alerts = append(alerts, models.Alert{
			Type:        "critical",
			Message:     "Low Battery Voltage",
			Description: fmt.Sprintf("Battery voltage is %.2fV. Battery may need charging or replacement.", telemetry.BatteryVoltage),
		})
	}

	// ABS errors
	if telemetry.ABSSensorStatus == "error" {
		alerts = append(alerts, models.Alert{
			Type:        "critical",
			Message:     "ABS System Error",
			Description: "ABS sensor system is reporting errors. Immediate inspection recommended.",
		})
	}

	// Airbag errors
	if telemetry.AirbagSystemStatus == "error" {
		alerts = append(alerts, models.Alert{
			Type:        "critical",
			Message:     "Airbag System Error",
			Description: "Airbag system is reporting errors. Safety inspection required immediately.",
		})
	}

	// Emissions system not ready
	if !telemetry.EmissionsMonitorReady {
		alerts = append(alerts, models.Alert{
			Type:        "warning",
			Message:     "Emissions Monitor Not Ready",
			Description: "Emissions monitoring system is not ready. May affect vehicle inspection.",
		})
	}

	// Low fuel level (<10%)
	if telemetry.FuelLevel < 10.0 {
		alerts = append(alerts, models.Alert{
			Type:        "warning",
			Message:     "Low Fuel Level",
			Description: fmt.Sprintf("Fuel level is at %.1f%%. Consider refueling soon.", telemetry.FuelLevel),
		})
	}

	// Critical error codes
	for _, code := range errorCodes {
		if code.Criticality == "critical" {
			alerts = append(alerts, models.Alert{
				Type:        "critical",
				Message:     fmt.Sprintf("DTC: %s", code.Code),
				Description: code.Description,
			})
		}
	}

	return alerts
}

// generateSuggestions creates maintenance suggestions
func (s *AnalysisService) generateSuggestions(analysis *models.Analysis, car *models.Car, telemetry *models.Telemetry) {
	if analysis.OilChangeDue {
		analysis.Suggestions = append(analysis.Suggestions, "Schedule an oil change service. Current mileage exceeds 10,000 km since last change.")
	}

	if analysis.SparkPlugsDue {
		analysis.Suggestions = append(analysis.Suggestions, "Spark plugs replacement is due. Recommended interval: 30,000 km.")
	}

	if analysis.CoolantDue {
		analysis.Suggestions = append(analysis.Suggestions, "Coolant replacement is due. Recommended interval: 60,000 km.")
	}

	if analysis.BrakeFluidDue {
		analysis.Suggestions = append(analysis.Suggestions, "Brake fluid replacement is due. Recommended interval: 40,000 km.")
	}

	if analysis.AirFilterDue {
		analysis.Suggestions = append(analysis.Suggestions, "Air filter replacement is due. Recommended interval: 20,000 km.")
	}

	// Tire pressure suggestions
	for i, pressure := range telemetry.TirePressure {
		if pressure < 28 {
			analysis.Suggestions = append(analysis.Suggestions, fmt.Sprintf("Tire %d pressure is low (%.0f PSI). Consider inflating.", i+1, float64(pressure)))
		} else if pressure > 40 {
			analysis.Suggestions = append(analysis.Suggestions, fmt.Sprintf("Tire %d pressure is high (%.0f PSI). Consider deflating.", i+1, float64(pressure)))
		}
	}

	// Fuel consumption suggestion
	if telemetry.AverageFuelConsumption > 12.0 && car.EngineType != "diesel" {
		analysis.Suggestions = append(analysis.Suggestions, "High fuel consumption detected. Consider checking air filter, spark plugs, and tire pressure.")
	}
}
