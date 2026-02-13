package data

import (
    "encoding/json"
    "os"
    "sync"
)

var (
    // STOChanges maps vehicle_id -> arbitrary change object
    STOChanges = make(map[string]map[string]interface{})
    stoMu      sync.Mutex
)

// LoadSTOChanges loads sto_changes.json if present
func LoadSTOChanges() error {
    stoMu.Lock()
    defer stoMu.Unlock()

    b, err := os.ReadFile("sto_changes.json")
    if err != nil {
        if os.IsNotExist(err) {
            STOChanges = make(map[string]map[string]interface{})
            return nil
        }
        return err
    }
    if err := json.Unmarshal(b, &STOChanges); err != nil {
        return err
    }
    return nil
}

// SaveSTOChanges persists STOChanges to disk
func SaveSTOChanges() error {
    stoMu.Lock()
    defer stoMu.Unlock()
    data, err := json.MarshalIndent(STOChanges, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile("sto_changes.json", data, 0644)
}

// UpdateSTOChange merges change map into existing entry for vehicleID and persists
func UpdateSTOChange(vehicleID string, change map[string]interface{}) error {
    stoMu.Lock()
    defer stoMu.Unlock()
    cur, ok := STOChanges[vehicleID]
    if !ok {
        cur = make(map[string]interface{})
    }
    for k, v := range change {
        cur[k] = v
    }
    STOChanges[vehicleID] = cur
    data, err := json.MarshalIndent(STOChanges, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile("sto_changes.json", data, 0644)
}

// GetSTOChange returns stored change object for vehicleID
func GetSTOChange(vehicleID string) map[string]interface{} {
    stoMu.Lock()
    defer stoMu.Unlock()
    return STOChanges[vehicleID]
}
