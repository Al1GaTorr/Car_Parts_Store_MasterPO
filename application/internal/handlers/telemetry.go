package handlers

import (
    "car-monitoring/internal/data"
    "car-monitoring/internal/models"
    "encoding/json"
    "net/http"
)

// TelemetryHandler handles POST /api/telemetry
func TelemetryHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var t models.Telemetry
    if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
        http.Error(w, "invalid payload", http.StatusBadRequest)
        return
    }

    // Update in-memory telemetry store
    if data.MockTelemetry == nil {
        data.MockTelemetry = make(map[string]models.Telemetry)
    }
    data.MockTelemetry[t.CarID] = t

    // Apply simple rules and build vehicle state payload
    status := "OK"
    if t.CoolantTemp > 100 {
        status = "WARNING"
    }
    if t.CoolantTemp > 120 {
        status = "CRITICAL"
    }

    // Try to get mileage from MockCars
    mileage := 0
    for _, c := range data.MockCars {
        if c.ID == t.CarID {
            mileage = c.Mileage
            break
        }
    }

    payload := map[string]interface{}{
        "mileage":   mileage,
        "engineTemp": t.CoolantTemp,
        "status":    status,
    }

    // Broadcast vehicle state update
    BroadcastToVIN(t.CarID, "VEHICLE_STATE_UPDATED", payload)

    // If critical, create a simple alert and broadcast ALERT_CREATED
    if status == "CRITICAL" {
        alert := map[string]interface{}{
            "title":   "Engine Overheat",
            "message": "Engine coolant temperature critical",
            "temp":    t.CoolantTemp,
        }
        BroadcastToVIN(t.CarID, "ALERT_CREATED", alert)
    }

    w.WriteHeader(http.StatusAccepted)
    json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}
