package models

// Car represents a vehicle with identification and basic information
type Car struct {
	ID            string `json:"id"`
	VIN           string `json:"vin"`
	Brand         string `json:"brand"`
	Model         string `json:"model"`
	Year          int    `json:"year"`
	EngineType  string  `json:"engineType"` // gasoline, diesel, hybrid
	LicensePlate  string `json:"licensePlate"`
	Mileage       int    `json:"mileage"` // in km
	LastOilChange int    `json:"lastOilChange"` // mileage at last oil change
}

