package models

// Vehicle represents a vehicle from mock_vehicles.json
type Vehicle struct {
	VehicleID         string              `json:"vehicle_id"`
	Brand             string              `json:"brand"`
	Model             string              `json:"model"`
	Year              int                 `json:"year"`
	MileageKm         int                 `json:"mileage_km"`
	EngineType        string              `json:"engine_type"`
	LastServiceKm     int                 `json:"last_service_km"`
	Errors            []VehicleError      `json:"errors"`
	MaintenanceAlerts []string            `json:"maintenance_alerts"`
	RiskLevel         string              `json:"risk_level"`
	Location          string              `json:"location"`
	DTPHistory        bool                `json:"dtp_history"`
}

// VehicleError represents an error code with severity
type VehicleError struct {
	Code              string `json:"code"`
	Severity          string `json:"severity"` // critical, medium, info
	Description       string `json:"description"`
	RecommendedAction string `json:"recommended_action"`
}

// SelectedVehicle tracks which vehicle is currently selected
type SelectedVehicle struct {
	VehicleID string `json:"vehicle_id"`
}

