package data

import (
    "car-monitoring/internal/models"
    "encoding/json"
    "os"
    "sync"
)

var (
    // MockServiceHistory maps vehicle_id -> list of service history entries
    MockServiceHistory = make(map[string][]models.ServiceHistory)
    historyMu          sync.Mutex
)

// LoadServiceHistory loads service_history.json if present
func LoadServiceHistory() error {
    historyMu.Lock()
    defer historyMu.Unlock()

    data, err := os.ReadFile("service_history.json")
    if err != nil {
        if os.IsNotExist(err) {
            MockServiceHistory = make(map[string][]models.ServiceHistory)
            return nil
        }
        return err
    }
    if err := json.Unmarshal(data, &MockServiceHistory); err != nil {
        return err
    }
    return nil
}

// SaveServiceHistory persists MockServiceHistory to disk
func SaveServiceHistory() error {
    historyMu.Lock()
    defer historyMu.Unlock()

    data, err := json.MarshalIndent(MockServiceHistory, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile("service_history.json", data, 0644)
}

// AppendServiceHistory adds a record for vehicleID and persists to disk
func AppendServiceHistory(vehicleID string, rec models.ServiceHistory) error {
    historyMu.Lock()
    defer historyMu.Unlock()
    MockServiceHistory[vehicleID] = append(MockServiceHistory[vehicleID], rec)
    data, err := json.MarshalIndent(MockServiceHistory, "", "  ")
    if err != nil {
        return err
    }
    if err := os.WriteFile("service_history.json", data, 0644); err != nil {
        return err
    }

    // Also mirror into STOChanges history for quick access by frontend
    // Do not import data.UpdateSTOChange here to avoid shadowing; call package-level function
    // Build a simple map entry for this vehicle's history
    // Note: UpdateSTOChange is in same package so we can call it directly
    change := map[string]interface{}{
        "history": MockServiceHistory[vehicleID],
    }
    _ = UpdateSTOChange(vehicleID, change)
    return nil
}

// GetServiceHistory returns history for vehicleID
func GetServiceHistory(vehicleID string) []models.ServiceHistory {
    historyMu.Lock()
    defer historyMu.Unlock()
    return MockServiceHistory[vehicleID]
}
