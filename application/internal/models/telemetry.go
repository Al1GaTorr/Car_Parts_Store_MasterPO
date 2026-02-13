package models

// Telemetry contains all real-time vehicle sensor data
type Telemetry struct {
	CarID                    string   `json:"carId"`
	Speed                    float64  `json:"speed"`                    // km/h
	RPM                      float64  `json:"rpm"`                      // revolutions per minute
	EngineLoad               float64  `json:"engineLoad"`               // percentage
	ThrottlePosition         float64  `json:"throttlePosition"`         // percentage
	IntakeAirTemp            float64  `json:"intakeAirTemp"`            // Celsius
	CoolantTemp              float64  `json:"coolantTemp"`              // Celsius
	FuelLevel                float64  `json:"fuelLevel"`                // percentage
	AverageFuelConsumption   float64  `json:"averageFuelConsumption"`   // l/100km
	InstantFuelConsumption   float64  `json:"instantFuelConsumption"`   // l/100km
	AverageSpeed             float64  `json:"averageSpeed"`             // km/h
	TripDistance             float64  `json:"tripDistance"`             // km since reset
	FuelSystemStatus         string   `json:"fuelSystemStatus"`         // open/closed loop
	O2SensorVoltages         []float64 `json:"o2SensorVoltages"`        // array of voltages
	ShortTermFuelTrim        float64  `json:"shortTermFuelTrim"`        // percentage
	LongTermFuelTrim         float64  `json:"longTermFuelTrim"`          // percentage
	IntakeManifoldPressure    float64  `json:"intakeManifoldPressure"`    // kPa
	BatteryVoltage           float64  `json:"batteryVoltage"`           // volts
	ABSSensorStatus          string   `json:"absSensorStatus"`          // ok/error
	TirePressure           []float64  `json:"tirePressure"`              // PSI for each tire
	BrakeSystemStatus        string   `json:"brakeSystemStatus"`        // ok/warning/error
	AirbagSystemStatus       string   `json:"airbagSystemStatus"`        // ok/error
	EmissionsMonitorReady    bool     `json:"emissionsMonitorReady"`    // true if ready
	FreezeFrameData          string   `json:"freezeFrameData"`           // encoded freeze frame
	Timestamp                int64    `json:"timestamp"`                 // Unix timestamp
}

