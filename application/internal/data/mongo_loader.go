package data

import (
  "car-monitoring/internal/models"
  "encoding/json"
  "errors"
  "os/exec"
  "strings"
)

type mongoCarDoc struct {
  VIN       string `json:"vin"`
  Make      string `json:"make"`
  Model     string `json:"model"`
  Year      int    `json:"year"`
  MileageKm int    `json:"mileage_km"`
  Issues    []struct {
    Code     string `json:"code"`
    Title    string `json:"title"`
    Severity string `json:"severity"` // high|medium|low
    DTC      *struct {
      Code string `json:"code"`
      Desc string `json:"desc"`
    } `json:"dtc"`
  } `json:"issues"`
}

// LoadVehiclesFromMongo loads vehicles from MongoDB using mongosh and fills MockVehicles.
// mongoURI example: mongodb://localhost:27017/bazarPO
func LoadVehiclesFromMongo(mongoURI string) error {
  // Use --quiet to avoid extra text. We output strict JSON with JSON.stringify.
  eval := "JSON.stringify(db.cars.find({}, { _id: 0 }).toArray())"
  cmd := exec.Command("mongosh", "--quiet", mongoURI, "--eval", eval)
  out, err := cmd.Output()
  if err != nil {
    return err
  }

  raw := strings.TrimSpace(string(out))
  if raw == "" {
    return errors.New("empty response from mongosh")
  }

  var docs []mongoCarDoc
  if err := json.Unmarshal([]byte(raw), &docs); err != nil {
    return err
  }

  vehicles := make([]models.Vehicle, 0, len(docs))
  for _, d := range docs {
    if strings.TrimSpace(d.VIN) == "" {
      continue
    }

    // Map issues -> VehicleError
    errs := make([]models.VehicleError, 0, len(d.Issues))
    maintenanceAlerts := make([]string, 0)
    maxSeverity := "info"

    for _, iss := range d.Issues {
      sev := mapIssueSeverity(iss.Severity)
      if severityRank(sev) > severityRank(maxSeverity) {
        maxSeverity = sev
      }

      desc := iss.Title
      if desc == "" && iss.DTC != nil {
        desc = iss.DTC.Desc
      }

      rec := defaultRecommendedAction(iss.Code)

      errs = append(errs, models.VehicleError{
        Code:              iss.Code,
        Severity:          sev,
        Description:       desc,
        RecommendedAction: rec,
      })

      // Simple maintenance list for UI
      if strings.HasSuffix(iss.Code, "_due") || strings.Contains(iss.Code, "_worn") || strings.Contains(iss.Code, "_dirty") {
        maintenanceAlerts = append(maintenanceAlerts, desc)
      }
    }

    // Approximate last service based on mileage (every 10k)
    lastService := d.MileageKm - (d.MileageKm % 10000)
    if lastService < 0 {
      lastService = 0
    }

    vehicles = append(vehicles, models.Vehicle{
      VehicleID:         d.VIN,
      Brand:             d.Make,
      Model:             d.Model,
      Year:              d.Year,
      MileageKm:         d.MileageKm,
      EngineType:        "",
      LastServiceKm:     lastService,
      Errors:            errs,
      MaintenanceAlerts: maintenanceAlerts,
      RiskLevel:         maxSeverity,
      Location:          "",
      DTPHistory:        false,
    })
  }

  if len(vehicles) == 0 {
    return errors.New("no vehicles found in MongoDB")
  }

  MockVehicles = vehicles
  // Set default selected vehicle to first VIN
  CurrentlySelectedVehicle = vehicles[0].VehicleID
  return nil
}

func mapIssueSeverity(s string) string {
  switch strings.ToLower(strings.TrimSpace(s)) {
  case "high", "critical":
    return "critical"
  case "medium", "warning":
    return "medium"
  case "low", "info":
    return "info"
  default:
    return "info"
  }
}

func severityRank(s string) int {
  switch s {
  case "critical":
    return 3
  case "medium":
    return 2
  case "info":
    return 1
  default:
    return 0
  }
}

func defaultRecommendedAction(issueCode string) string {
  // Keep short: this is a demo.
  switch issueCode {
  case "oil_service_due":
    return "Replace engine oil and oil filter"
  case "spark_plugs_due":
    return "Replace spark plugs"
  case "brake_pads_worn":
    return "Inspect and replace brake pads"
  case "battery_weak":
    return "Check battery and charging system"
  case "air_filter_dirty":
    return "Replace air filter"
  case "coolant_low", "overheating_risk":
    return "Check coolant level and cooling system"
  case "abs_fault":
    return "Run ABS diagnostics (sensor/wiring)"
  case "check_engine":
    return "Run OBD diagnostics and inspect engine systems"
  default:
    return "Inspect and diagnose the issue"
  }
}
