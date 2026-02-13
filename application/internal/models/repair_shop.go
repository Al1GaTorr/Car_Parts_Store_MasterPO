package models

// RepairShop represents a repair shop/service station from 2GIS
type RepairShop struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Rating     float64  `json:"rating"`
	Reviews    int      `json:"reviews"`
	Distance   float64  `json:"distance"` // km
	Address    string   `json:"address"`
	Phone      string   `json:"phone"`
	Hours      string   `json:"hours"`
	Services   []string `json:"services"`
	Verified   bool     `json:"verified"`
	PriceLevel int      `json:"priceLevel"` // 1-3 (cheap, medium, expensive)
	Latitude   float64  `json:"latitude,omitempty"`
	Longitude  float64  `json:"longitude,omitempty"`
}

