package models

// ServiceStation represents a car service station
type ServiceStation struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Latitude      float64  `json:"latitude"`
	Longitude     float64  `json:"longitude"`
	Distance      float64  `json:"distance"`      // km from user
	Rating        float64  `json:"rating"`        // 1-5
	PriceLevel    string   `json:"priceLevel"`    // cheap, medium, expensive
	SupportedBrands []string `json:"supportedBrands"` // array of car brands
	Specializations []string `json:"specializations"` // engine, transmission, electrical, etc.
	City          string   `json:"city"`
	Address       string   `json:"address"`
}

