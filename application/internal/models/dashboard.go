package models

// DashboardData represents the dashboard view data
type DashboardData struct {
	CarInfo        CarInfo        `json:"carInfo"`
	HealthMetrics  []HealthMetric `json:"healthMetrics"`
	OilChangeData  OilChangeData  `json:"oilChangeData"`
	RecentAlerts   []Alert        `json:"recentAlerts"`
}

// CarInfo represents basic car information for dashboard
type CarInfo struct {
	VIN  string `json:"vin"`
	Model string `json:"model"`
	Year  string `json:"year"`
	Plate string `json:"plate"`
	Mileage int  `json:"mileage"`
}

// HealthMetric represents a vehicle health metric
type HealthMetric struct {
	Label  string `json:"label"`
	Value  int    `json:"value"`
	Status string `json:"status"` // good, warning, critical
	Icon   string `json:"icon"`
	Color  string `json:"color"`
}

// OilChangeData represents oil change information
type OilChangeData struct {
	CurrentKm    int `json:"currentKm"`
	NextChangeKm int `json:"nextChangeKm"`
	DaysRemaining int `json:"daysRemaining"`
}

