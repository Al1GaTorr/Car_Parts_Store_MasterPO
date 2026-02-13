package handlers

import (
    "car-monitoring/internal/data"
    "car-monitoring/internal/models"
    "encoding/json"
    "net/http"
    "strings"
)

// STOPanelHandler handles /api/sto-panel/... requests
func STOPanelHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    // Trim prefix /api/sto-panel/
    path := strings.TrimPrefix(r.URL.Path, "/api/sto-panel/")

    // GET /api/sto-panel/vehicles -> list
    if path == "vehicles" && r.Method == http.MethodGet {
        json.NewEncoder(w).Encode(data.MockVehicles)
        return
    }

    // GET /api/sto-panel/vehicles/{id}/history -> return persisted history for vehicle
    if strings.HasPrefix(path, "vehicles/") && strings.HasSuffix(path, "/history") && r.Method == http.MethodGet {
        id := strings.TrimSuffix(strings.TrimPrefix(path, "vehicles/"), "/history")
        hist := data.GetServiceHistory(id)
        json.NewEncoder(w).Encode(hist)
        return
    }

    // GET /api/sto-panel/vehicles/{id}/changes -> return STO changes summary for vehicle (includes history, parts, last_update)
    if strings.HasPrefix(path, "vehicles/") && strings.HasSuffix(path, "/changes") && r.Method == http.MethodGet {
        id := strings.TrimSuffix(strings.TrimPrefix(path, "vehicles/"), "/changes")
        ch := data.GetSTOChange(id)
        if ch == nil {
            ch = map[string]interface{}{}
        }
        json.NewEncoder(w).Encode(ch)
        return
    }

    // GET /api/sto-panel/vehicles/{id}
    if strings.HasPrefix(path, "vehicles/") {
        id := strings.TrimPrefix(path, "vehicles/")
        // POST /api/sto-panel/vehicles/{id} -> update vehicle
        if r.Method == http.MethodGet {
            for _, v := range data.MockVehicles {
                if strings.EqualFold(v.VehicleID, id) || strings.EqualFold(v.Brand+v.Model, id) || strings.EqualFold(v.VehicleID, strings.ToUpper(id)) {
                    json.NewEncoder(w).Encode(v)
                    return
                }
            }
            http.Error(w, "not found", http.StatusNotFound)
            return
        }

        if r.Method == http.MethodPost {
            var upd map[string]interface{}
            if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
                http.Error(w, "invalid payload", http.StatusBadRequest)
                return
            }

            // Find and update vehicle in memory
            for i := range data.MockVehicles {
                v := &data.MockVehicles[i]
                if strings.EqualFold(v.VehicleID, id) {
                    // Apply some simple updates: mileage_km, maintenance_alerts, errors
                    if mm, ok := upd["mileage_km"].(float64); ok {
                        v.MileageKm = int(mm)
                    }
                    if alerts, ok := upd["maintenance_alerts"].([]interface{}); ok {
                        newAlerts := []string{}
                        for _, a := range alerts {
                            if s, ok := a.(string); ok {
                                newAlerts = append(newAlerts, s)
                            }
                        }
                        v.MaintenanceAlerts = newAlerts
                    }
                    if errs, ok := upd["errors"].([]interface{}); ok {
                        newErrs := []models.VehicleError{}
                        for _, e := range errs {
                            if m, ok := e.(map[string]interface{}); ok {
                                newErrs = append(newErrs, models.VehicleError{
                                    Code:        toStr(m["code"]),
                                    Severity:    toStr(m["severity"]),
                                    Description: toStr(m["description"]),
                                })
                            }
                        }
                        v.Errors = newErrs
                    }

                    // persist vehicles file
                    if err := data.SaveMockVehicles(); err != nil {
                        http.Error(w, "failed to persist", http.StatusInternalServerError)
                        return
                    }

                    // update STO changes summary (maintenance alerts, mileage)
                    _ = data.UpdateSTOChange(v.VehicleID, map[string]interface{}{
                        "maintenance_alerts": v.MaintenanceAlerts,
                        "mileage_km": v.MileageKm,
                    })

                    // Broadcast updated vehicle state to subscribers
                    payload := map[string]interface{}{
                        "maintenance_alerts": v.MaintenanceAlerts,
                        "mileage": v.MileageKm,
                        "vehicle_id": v.VehicleID,
                    }
                    BroadcastToVIN(v.VehicleID, "VEHICLE_STATE_UPDATED", payload)

                    json.NewEncoder(w).Encode(v)
                    return
                }
            }
            http.Error(w, "not found", http.StatusNotFound)
            return
        }
    }

    // POST /api/sto-panel/vehicles/{id}/records -> add service history record
    if strings.HasPrefix(path, "vehicles/") && strings.HasSuffix(path, "/records") && r.Method == http.MethodPost {
        // extract id
        id := strings.TrimSuffix(strings.TrimPrefix(path, "vehicles/"), "/records")

        // decode into generic map to accept extra fields like parts, serviceName
        var raw map[string]interface{}
        if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
            http.Error(w, "invalid payload", http.StatusBadRequest)
            return
        }

        // build models.ServiceHistory from raw
        var rec models.ServiceHistory
        if d, ok := raw["date"].(string); ok {
            rec.Date = d
        }
        if s, ok := raw["serviceName"].(string); ok {
            rec.Type = s
        } else if s, ok := raw["type"].(string); ok {
            rec.Type = s
        }
        if desc, ok := raw["description"].(string); ok {
            rec.Description = desc
        }
        if mm, ok := raw["mileage"].(float64); ok {
            rec.Mileage = int(mm)
        } else if mm, ok := raw["mileage"].(int); ok {
            rec.Mileage = mm
        }
        if cost, ok := raw["cost"].(float64); ok {
            rec.Cost = int(cost)
        }
        if shop, ok := raw["shop"].(string); ok {
            rec.Shop = shop
        }
        if loc, ok := raw["location"].(string); ok {
            rec.Location = loc
        }
        if parts, ok := raw["parts"].([]interface{}); ok {
            rec.Parts = []map[string]interface{}{}
            for _, p := range parts {
                if pm, ok := p.(map[string]interface{}); ok {
                    rec.Parts = append(rec.Parts, pm)
                }
            }
        }

        // append to in-memory map and persist
        if err := data.AppendServiceHistory(id, rec); err != nil {
            http.Error(w, "failed to save history", http.StatusInternalServerError)
            return
        }

        // Also store a quick STO change summary so SPA/others can read parts/last update
        change := map[string]interface{}{
            "last_update": rec.Date,
            "last_mileage": rec.Mileage,
            "last_service": rec.Type,
            "parts": rec.Parts,
        }
        _ = data.UpdateSTOChange(id, change)

        // Broadcast to VIN subscribers (use id as VIN here)
        BroadcastToVIN(id, "SERVICE_RECORD_ADDED", rec)

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
        return
    }

    http.Error(w, "not found", http.StatusNotFound)
}

func toStr(v interface{}) string {
    if v == nil {
        return ""
    }
    if s, ok := v.(string); ok {
        return s
    }
    return ""
}
