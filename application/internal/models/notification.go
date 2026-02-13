package models

// Notification represents a user notification
type Notification struct {
	ID         string `json:"id"`
	Type       string `json:"type"` // warning, critical, info
	Title      string `json:"title"`
	Message    string `json:"message"`
	Timestamp  int64  `json:"timestamp"`
	Read       bool   `json:"read"`
	Actionable bool   `json:"actionable"`
}

