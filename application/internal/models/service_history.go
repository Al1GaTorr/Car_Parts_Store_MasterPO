package models

// ServiceHistory represents a service record
type ServiceHistory struct {
	ID          int    `json:"id"`
	Date        string `json:"date"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Mileage     int    `json:"mileage"`
	Cost        int    `json:"cost"`
	Shop        string `json:"shop"`
	Location    string `json:"location"`
	Verified    bool   `json:"verified"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	Parts       []map[string]interface{} `json:"parts,omitempty"`
}

// ServiceHistoryStats represents statistics about service history
type ServiceHistoryStats struct {
	TotalServices int `json:"totalServices"`
	TotalSpent    int `json:"totalSpent"`
}

