package data

import (
	"car-monitoring/internal/models"
	"encoding/json"
	"errors"
	"os"
)

var (
	// MockVehicles contains vehicles from mock_vehicles.json
	MockVehicles []models.Vehicle
	// ErrorSeverityMap contains error code severity mappings
	ErrorSeverityMap map[string]map[string]string
	// CurrentlySelectedVehicle is the vehicle ID selected in admin panel
	CurrentlySelectedVehicle = "V001" // Default to first vehicle
)

// LoadMockVehicles loads vehicles from mock_vehicles.json
func LoadMockVehicles() error {
	data, err := readFileWithFallback("mock_vehicles.json")
	if err != nil {
		return err
	}

	var vehiclesData struct {
		Vehicles []models.Vehicle `json:"vehicles"`
	}

	if err := json.Unmarshal(data, &vehiclesData); err != nil {
		return err
	}

	MockVehicles = vehiclesData.Vehicles
	return nil
}

// SaveMockVehicles writes current MockVehicles back to mock_vehicles.json
func SaveMockVehicles() error {
	var out struct {
		GeneratedAt string           `json:"generated_at"`
		Source      string           `json:"source"`
		Vehicles    []models.Vehicle `json:"vehicles"`
	}
	out.GeneratedAt = "2026-01-24T12:00:00Z"
	out.Source = "mock_obd_provider"
	out.Vehicles = MockVehicles

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	path := resolveDataPath("mock_vehicles.json")
	return os.WriteFile(path, data, 0644)
}

// LoadErrorSeverityMap loads error severity mappings
func LoadErrorSeverityMap() error {
	data, err := readFileWithFallback("error_severity_map.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &ErrorSeverityMap); err != nil {
		return err
	}

	return nil
}

func resolveDataPath(name string) string {
	candidates := []string{
		name,
		"../../" + name,
		"../" + name,
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return name
}

func readFileWithFallback(name string) ([]byte, error) {
	candidates := []string{
		name,
		"../../" + name,
		"../" + name,
	}
	var lastErr error
	for _, p := range candidates {
		b, err := os.ReadFile(p)
		if err == nil {
			return b, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, errors.New("file not found")
}

// GetSelectedVehicle returns the currently selected vehicle
func GetSelectedVehicle() *models.Vehicle {
	for i := range MockVehicles {
		if MockVehicles[i].VehicleID == CurrentlySelectedVehicle {
			return &MockVehicles[i]
		}
	}
	// Return first vehicle if selected not found
	if len(MockVehicles) > 0 {
		return &MockVehicles[0]
	}
	return nil
}

// SetSelectedVehicle sets the currently selected vehicle
func SetSelectedVehicle(vehicleID string) bool {
	for i := range MockVehicles {
		if MockVehicles[i].VehicleID == vehicleID {
			CurrentlySelectedVehicle = vehicleID
			return true
		}
	}
	return false
}

// GetErrorSeverityDescription returns the severity description for an error code
func GetErrorSeverityDescription(code string) (string, string) {
	// Check critical
	if desc, ok := ErrorSeverityMap["critical"][code]; ok {
		return "critical", desc
	}
	// Check medium
	if desc, ok := ErrorSeverityMap["medium"][code]; ok {
		return "medium", desc
	}
	// Check info
	if desc, ok := ErrorSeverityMap["info"][code]; ok {
		return "info", desc
	}
	return "unknown", "Неизвестная ошибка"
}
