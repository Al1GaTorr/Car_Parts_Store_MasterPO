package data

import (
	"car-monitoring/internal/models"
	"math/rand"
	"time"
)

var (
	// MockCars contains simulated vehicle data
	MockCars []models.Car
	// MockTelemetry contains current telemetry for each car
	MockTelemetry map[string]models.Telemetry
	// MockErrorCodes contains DTC codes for each car
	MockErrorCodes map[string][]models.ErrorCode
	// MockServiceStations contains service station data
	MockServiceStations []models.ServiceStation
)

// InitializeMockData populates all mock data structures
func InitializeMockData() {
	rand.Seed(time.Now().UnixNano())

	// Initialize cars
	MockCars = []models.Car{
		{
			ID:            "car1",
			VIN:           "1HGBH41JXMN109186",
			Brand:         "Toyota",
			Model:         "Camry",
			Year:          2020,
			EngineType:    "gasoline",
			LicensePlate:  "ABC-1234",
			Mileage:       45000,
			LastOilChange: 40000,
		},
		{
			ID:            "car2",
			VIN:           "WBA3A5C58ED123456",
			Brand:         "BMW",
			Model:         "320i",
			Year:          2019,
			EngineType:    "gasoline",
			LicensePlate:  "XYZ-5678",
			Mileage:       78000,
			LastOilChange: 68000,
		},
		{
			ID:            "car3",
			VIN:           "KMHDN45D17U123456",
			Brand:         "Hyundai",
			Model:         "Elantra",
			Year:          2021,
			EngineType:    "gasoline",
			LicensePlate:  "DEF-9012",
			Mileage:       25000,
			LastOilChange: 20000,
		},
		{
			ID:            "car4",
			VIN:           "1FTFW1ET5DFC12345",
			Brand:         "Ford",
			Model:         "F-150",
			Year:          2018,
			EngineType:    "diesel",
			LicensePlate:  "GHI-3456",
			Mileage:       95000,
			LastOilChange: 85000,
		},
		{
			ID:            "car5",
			VIN:           "5YJ3E1EB0KF123456",
			Brand:         "Tesla",
			Model:         "Model 3",
			Year:          2022,
			EngineType:    "hybrid",
			LicensePlate:  "JKL-7890",
			Mileage:       15000,
			LastOilChange: 10000,
		},
	}

	// Initialize telemetry map
	MockTelemetry = make(map[string]models.Telemetry)
	MockTelemetry["car1"] = models.Telemetry{
		CarID:                  "car1",
		Speed:                  65.5,
		RPM:                    2200,
		EngineLoad:             45.2,
		ThrottlePosition:       35.0,
		IntakeAirTemp:          28.5,
		CoolantTemp:            92.0,
		FuelLevel:              45.0,
		AverageFuelConsumption: 8.5,
		InstantFuelConsumption: 7.8,
		AverageSpeed:           58.3,
		TripDistance:           125.5,
		FuelSystemStatus:       "closed loop",
		O2SensorVoltages:       []float64{0.45, 0.48, 0.44, 0.46},
		ShortTermFuelTrim:      2.5,
		LongTermFuelTrim:       1.2,
		IntakeManifoldPressure: 45.2,
		BatteryVoltage:         12.6,
		ABSSensorStatus:        "ok",
		TirePressure:           []float64{32, 33, 31, 32},
		BrakeSystemStatus:      "ok",
		AirbagSystemStatus:     "ok",
		EmissionsMonitorReady:  true,
		FreezeFrameData:        "P0123-Engine Coolant Temperature Sensor Circuit High",
		Timestamp:              time.Now().Unix(),
	}

	MockTelemetry["car2"] = models.Telemetry{
		CarID:                  "car2",
		Speed:                  120.0,
		RPM:                    3500,
		EngineLoad:             78.5,
		ThrottlePosition:       85.0,
		IntakeAirTemp:          35.2,
		CoolantTemp:            108.5, // High temperature alert
		FuelLevel:              8.0,   // Low fuel alert
		AverageFuelConsumption: 10.2,
		InstantFuelConsumption: 12.5,
		AverageSpeed:           95.0,
		TripDistance:           245.8,
		FuelSystemStatus:       "open loop",
		O2SensorVoltages:       []float64{0.52, 0.55, 0.48, 0.50},
		ShortTermFuelTrim:      -5.2,
		LongTermFuelTrim:       -3.8,
		IntakeManifoldPressure: 78.5,
		BatteryVoltage:         11.8, // Low battery alert
		ABSSensorStatus:        "error",
		TirePressure:           []float64{28, 29, 30, 28},
		BrakeSystemStatus:      "warning",
		AirbagSystemStatus:     "ok",
		EmissionsMonitorReady:  false,
		FreezeFrameData:        "P0300-Random/Multiple Cylinder Misfire Detected",
		Timestamp:              time.Now().Unix(),
	}

	MockTelemetry["car3"] = models.Telemetry{
		CarID:                  "car3",
		Speed:                  45.0,
		RPM:                    1800,
		EngineLoad:             32.0,
		ThrottlePosition:       25.0,
		IntakeAirTemp:          22.0,
		CoolantTemp:            88.0,
		FuelLevel:              75.0,
		AverageFuelConsumption: 7.2,
		InstantFuelConsumption: 6.5,
		AverageSpeed:           42.0,
		TripDistance:           89.3,
		FuelSystemStatus:       "closed loop",
		O2SensorVoltages:       []float64{0.44, 0.46, 0.45, 0.47},
		ShortTermFuelTrim:      1.5,
		LongTermFuelTrim:       0.8,
		IntakeManifoldPressure: 38.5,
		BatteryVoltage:         12.8,
		ABSSensorStatus:        "ok",
		TirePressure:           []float64{33, 34, 33, 34},
		BrakeSystemStatus:      "ok",
		AirbagSystemStatus:     "ok",
		EmissionsMonitorReady:  true,
		FreezeFrameData:        "",
		Timestamp:              time.Now().Unix(),
	}

	MockTelemetry["car4"] = models.Telemetry{
		CarID:                  "car4",
		Speed:                  85.0,
		RPM:                    2800,
		EngineLoad:             65.0,
		ThrottlePosition:       60.0,
		IntakeAirTemp:          30.0,
		CoolantTemp:            95.0,
		FuelLevel:              55.0,
		AverageFuelConsumption: 12.5,
		InstantFuelConsumption: 13.2,
		AverageSpeed:           78.0,
		TripDistance:           312.5,
		FuelSystemStatus:       "closed loop",
		O2SensorVoltages:       []float64{0.46, 0.48, 0.47, 0.49},
		ShortTermFuelTrim:      3.2,
		LongTermFuelTrim:       2.1,
		IntakeManifoldPressure: 65.0,
		BatteryVoltage:         12.4,
		ABSSensorStatus:        "ok",
		TirePressure:           []float64{35, 36, 35, 36},
		BrakeSystemStatus:      "ok",
		AirbagSystemStatus:     "error",
		EmissionsMonitorReady:  true,
		FreezeFrameData:        "B0001-Airbag System Malfunction",
		Timestamp:              time.Now().Unix(),
	}

	MockTelemetry["car5"] = models.Telemetry{
		CarID:                  "car5",
		Speed:                  95.0,
		RPM:                    0, // Electric motor
		EngineLoad:             55.0,
		ThrottlePosition:       70.0,
		IntakeAirTemp:          25.0,
		CoolantTemp:            75.0,
		FuelLevel:              80.0, // Battery level for hybrid
		AverageFuelConsumption: 5.5,
		InstantFuelConsumption: 4.8,
		AverageSpeed:           88.0,
		TripDistance:           198.2,
		FuelSystemStatus:       "closed loop",
		O2SensorVoltages:       []float64{0.45, 0.47},
		ShortTermFuelTrim:      0.5,
		LongTermFuelTrim:       0.2,
		IntakeManifoldPressure: 30.0,
		BatteryVoltage:         13.2,
		ABSSensorStatus:        "ok",
		TirePressure:           []float64{42, 43, 42, 43},
		BrakeSystemStatus:      "ok",
		AirbagSystemStatus:     "ok",
		EmissionsMonitorReady:  true,
		FreezeFrameData:        "",
		Timestamp:              time.Now().Unix(),
	}

	// Initialize error codes
	MockErrorCodes = make(map[string][]models.ErrorCode)
	MockErrorCodes["car1"] = []models.ErrorCode{
		{
			Code:        "P0123",
			Description: "Engine Coolant Temperature Sensor Circuit High",
			Criticality: "non-critical",
			MILStatus:   false,
		},
	}
	MockErrorCodes["car2"] = []models.ErrorCode{
		{
			Code:        "P0300",
			Description: "Random/Multiple Cylinder Misfire Detected",
			Criticality: "critical",
			MILStatus:   true,
		},
		{
			Code:        "C1201",
			Description: "ABS Control Module Communication Error",
			Criticality: "critical",
			MILStatus:   true,
		},
	}
	MockErrorCodes["car3"] = []models.ErrorCode{} // No errors
	MockErrorCodes["car4"] = []models.ErrorCode{
		{
			Code:        "B0001",
			Description: "Airbag System Malfunction",
			Criticality: "critical",
			MILStatus:   true,
		},
	}
	MockErrorCodes["car5"] = []models.ErrorCode{} // No errors

	// Initialize service stations
	MockServiceStations = []models.ServiceStation{
		{
			ID:              "sto1",
			Name:            "AutoService Center",
			Latitude:        55.7558,
			Longitude:       37.6173,
			Distance:        2.5,
			Rating:          4.5,
			PriceLevel:      "medium",
			SupportedBrands: []string{"Toyota", "Hyundai", "Ford"},
			Specializations: []string{"engine", "transmission", "electrical"},
			City:            "Moscow",
			Address:         "Tverskaya St, 10",
		},
		{
			ID:              "sto2",
			Name:            "BMW Service Premium",
			Latitude:        55.7520,
			Longitude:       37.6156,
			Distance:        3.2,
			Rating:          4.8,
			PriceLevel:      "expensive",
			SupportedBrands: []string{"BMW", "Mercedes", "Audi"},
			Specializations: []string{"engine", "electrical", "diagnostics"},
			City:            "Moscow",
			Address:         "Kutuzovsky Ave, 12",
		},
		{
			ID:              "sto3",
			Name:            "Quick Fix Garage",
			Latitude:        55.7580,
			Longitude:       37.6200,
			Distance:        1.8,
			Rating:          3.5,
			PriceLevel:      "cheap",
			SupportedBrands: []string{"Toyota", "Hyundai", "Ford", "BMW"},
			Specializations: []string{"engine", "brakes", "tires"},
			City:            "Moscow",
			Address:         "Arbat St, 25",
		},
		{
			ID:              "sto4",
			Name:            "Elite Auto Repair",
			Latitude:        55.7500,
			Longitude:       37.6100,
			Distance:        4.5,
			Rating:          4.9,
			PriceLevel:      "expensive",
			SupportedBrands: []string{"Tesla", "BMW", "Mercedes"},
			Specializations: []string{"electrical", "hybrid systems", "diagnostics"},
			City:            "Moscow",
			Address:         "Leninsky Ave, 50",
		},
		{
			ID:              "sto5",
			Name:            "Budget Auto Service",
			Latitude:        55.7600,
			Longitude:       37.6250,
			Distance:        2.1,
			Rating:          3.8,
			PriceLevel:      "cheap",
			SupportedBrands: []string{"Toyota", "Hyundai", "Ford"},
			Specializations: []string{"oil change", "filters", "brakes"},
			City:            "Moscow",
			Address:         "Spartakovskaya St, 15",
		},
		{
			ID:              "sto6",
			Name:            "Universal Car Service",
			Latitude:        55.7450,
			Longitude:       37.6300,
			Distance:        5.2,
			Rating:          4.2,
			PriceLevel:      "medium",
			SupportedBrands: []string{"Toyota", "BMW", "Hyundai", "Ford", "Tesla"},
			Specializations: []string{"engine", "transmission", "electrical", "brakes", "airbags"},
			City:            "Moscow",
			Address:         "Prospekt Mira, 100",
		},
	}
}

