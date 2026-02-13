package models

// Analysis contains maintenance recommendations and alerts
type Analysis struct {
	CarID           string      `json:"carId"`
	OilChangeDue    bool        `json:"oilChangeDue"`
	SparkPlugsDue  bool        `json:"sparkPlugsDue"`
	CoolantDue      bool        `json:"coolantDue"`
	BrakeFluidDue   bool        `json:"brakeFluidDue"`
	AirFilterDue    bool        `json:"airFilterDue"`
	Errors          []ErrorCode `json:"errors"`
	Alerts          []Alert     `json:"alerts"`
	Suggestions     []string    `json:"suggestions"`
}

// Alert represents a warning or critical issue
type Alert struct {
	Type        string `json:"type"`        // warning, critical
	Message     string `json:"message"`
	Description string `json:"description"`
	IssueCode   string `json:"issueCode,omitempty"`
}

