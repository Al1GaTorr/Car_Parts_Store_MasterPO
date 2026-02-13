package services

import (
	"car-monitoring/internal/models"
	"fmt"
)

// DashboardService handles dashboard data aggregation
type DashboardService struct {
	vehicleService  *VehicleService
	carService      *CarService
	analysisService *AnalysisService
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(vehicleService *VehicleService, carService *CarService, analysisService *AnalysisService) *DashboardService {
	return &DashboardService{
		vehicleService:  vehicleService,
		carService:      carService,
		analysisService: analysisService,
	}
}

// GetDashboardData returns complete dashboard data for the selected vehicle
func (s *DashboardService) GetDashboardData(carID string) (*models.DashboardData, error) {
	// Get selected vehicle
	vehicle, err := s.vehicleService.GetSelectedVehicle()
	if err != nil {
		return nil, err
	}

	// Convert vehicle to car for compatibility
	car := s.vehicleService.ConvertVehicleToCar(vehicle)

	// Get telemetry for vehicle
	telemetry := s.vehicleService.GetVehicleTelemetry(vehicle)

	// Store telemetry for analysis service compatibility
	s.carService.SetTelemetry(vehicle.VehicleID, *telemetry)

	// Try to get analysis, fallback to creating from vehicle data
	analysis, err := s.analysisService.GetAnalysis(vehicle.VehicleID)
	if err != nil {
		// If analysis fails, create basic analysis from vehicle data
		analysis = s.createAnalysisFromVehicle(vehicle, car)
	}

	// Build car info from vehicle
	carInfo := models.CarInfo{
		VIN:    car.VIN,
		Model:   vehicle.Brand + " " + vehicle.Model,
		Year:    fmt.Sprintf("%d", vehicle.Year),
		Plate:   "E 777 KZ", // Default plate, can be customized
		Mileage: vehicle.MileageKm,
	}

	// Build health metrics
	healthMetrics := []models.HealthMetric{
		{
			Label:  "Engine",
			Value:  s.calculateEngineHealth(telemetry),
			Status: s.getEngineStatus(telemetry),
			Icon:   "Gauge",
			Color:  s.getStatusColor(s.getEngineStatus(telemetry)),
		},
		{
			Label:  "Battery",
			Value:  s.calculateBatteryHealth(telemetry),
			Status: s.getBatteryStatus(telemetry),
			Icon:   "Battery",
			Color:  s.getStatusColor(s.getBatteryStatus(telemetry)),
		},
		{
			Label:  "Oil Level",
			Value:  s.calculateOilLevel(analysis),
			Status: s.getOilStatus(analysis),
			Icon:   "Droplet",
			Color:  s.getStatusColor(s.getOilStatus(analysis)),
		},
	}

	// Build oil change data from vehicle
	kmSinceLastChange := vehicle.MileageKm - vehicle.LastServiceKm
	nextChangeKm := vehicle.LastServiceKm + 10000
	if nextChangeKm < vehicle.MileageKm {
		nextChangeKm = vehicle.MileageKm + 1000 // Overdue
	}
	daysRemaining := s.calculateDaysRemaining(kmSinceLastChange)

	oilChangeData := models.OilChangeData{
		CurrentKm:     vehicle.MileageKm,
		NextChangeKm:  nextChangeKm,
		DaysRemaining: daysRemaining,
	}

	// Build recent alerts (limit to 3 most recent)
	recentAlerts := []models.Alert{}
	if len(analysis.Alerts) > 0 {
		maxAlerts := 3
		if len(analysis.Alerts) < maxAlerts {
			maxAlerts = len(analysis.Alerts)
		}
		recentAlerts = analysis.Alerts[:maxAlerts]
	}

	return &models.DashboardData{
		CarInfo:       carInfo,
		HealthMetrics: healthMetrics,
		OilChangeData: oilChangeData,
		RecentAlerts:  recentAlerts,
	}, nil
}

// Helper methods
func (s *DashboardService) calculateEngineHealth(telemetry *models.Telemetry) int {
	// Base health on coolant temp, engine load, and RPM
	health := 100

	// Coolant temp penalty
	if telemetry.CoolantTemp > 105 {
		health -= 30
	} else if telemetry.CoolantTemp > 95 {
		health -= 10
	}

	// Engine load penalty (very high load)
	if telemetry.EngineLoad > 90 {
		health -= 5
	}

	if health < 0 {
		health = 0
	}
	return health
}

func (s *DashboardService) getEngineStatus(telemetry *models.Telemetry) string {
	if telemetry.CoolantTemp > 105 {
		return "critical"
	}
	if telemetry.CoolantTemp > 95 || telemetry.EngineLoad > 90 {
		return "warning"
	}
	return "good"
}

func (s *DashboardService) calculateBatteryHealth(telemetry *models.Telemetry) int {
	// Battery health based on voltage
	if telemetry.BatteryVoltage >= 12.6 {
		return 95
	} else if telemetry.BatteryVoltage >= 12.4 {
		return 87
	} else if telemetry.BatteryVoltage >= 12.0 {
		return 75
	} else if telemetry.BatteryVoltage >= 11.8 {
		return 60
	}
	return 40
}

func (s *DashboardService) getBatteryStatus(telemetry *models.Telemetry) string {
	if telemetry.BatteryVoltage < 12.0 {
		return "critical"
	}
	if telemetry.BatteryVoltage < 12.4 {
		return "warning"
	}
	return "good"
}

func (s *DashboardService) calculateOilLevel(analysis *models.Analysis) int {
	// Oil level based on maintenance status
	if analysis.OilChangeDue {
		return 62 // Low
	}
	// Calculate based on mileage since last change
	return 85 // Good
}

func (s *DashboardService) getOilStatus(analysis *models.Analysis) string {
	if analysis.OilChangeDue {
		return "warning"
	}
	return "good"
}

func (s *DashboardService) getStatusColor(status string) string {
	switch status {
	case "critical":
		return "text-red-400"
	case "warning":
		return "text-yellow-400"
	default:
		return "text-emerald-400"
	}
}

func (s *DashboardService) calculateDaysRemaining(kmSinceLastChange int) int {
	// Estimate days based on average daily driving (50 km/day)
	kmRemaining := 10000 - kmSinceLastChange
	if kmRemaining < 0 {
		return 0
	}
	days := kmRemaining / 50
	if days < 1 {
		return 1
	}
	return days
}

// createAnalysisFromVehicle creates analysis directly from vehicle data
func (s *DashboardService) createAnalysisFromVehicle(vehicle *models.Vehicle, car *models.Car) *models.Analysis {
	// Check maintenance needs
	kmSinceService := vehicle.MileageKm - vehicle.LastServiceKm
	oilChangeDue := kmSinceService > 10000

	// Convert vehicle errors to alerts. For Mongo-backed vehicles, the code is an issue-code
	// like "air_filter_dirty" and severity is already mapped to critical/medium/info.
	alerts := []models.Alert{}
	for _, ve := range vehicle.Errors {
		alertType := "warning"
		switch ve.Severity {
		case "critical":
			alertType = "critical"
		case "medium":
			alertType = "warning"
		default:
			alertType = "info"
		}

		msg := ve.Description
		if msg == "" {
			msg = ve.Code
		}

		alerts = append(alerts, models.Alert{
			Type:        alertType,
			Message:     msg,
			Description: ve.RecommendedAction,
			IssueCode:   ve.Code,
		})
	}

	// Add maintenance alerts
	for _, maintAlert := range vehicle.MaintenanceAlerts {
		alerts = append(alerts, models.Alert{
			Type:        "warning",
			Message:     maintAlert,
			Description: maintAlert,
		})
	}

	// Convert vehicle errors to error codes
	errorCodes := s.vehicleService.GetVehicleErrorCodes(vehicle)

	return &models.Analysis{
		CarID:         vehicle.VehicleID,
		OilChangeDue:  oilChangeDue,
		SparkPlugsDue: (vehicle.MileageKm%30000) < 5000 && vehicle.MileageKm > 30000,
		CoolantDue:    (vehicle.MileageKm%60000) < 5000 && vehicle.MileageKm > 60000,
		BrakeFluidDue: (vehicle.MileageKm%40000) < 5000 && vehicle.MileageKm > 40000,
		AirFilterDue:  (vehicle.MileageKm%20000) < 5000 && vehicle.MileageKm > 20000,
		Errors:        errorCodes,
		Alerts:        alerts,
		Suggestions:   vehicle.MaintenanceAlerts,
	}
}
