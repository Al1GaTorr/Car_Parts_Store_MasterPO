package models

// ErrorCode represents a Diagnostic Trouble Code (DTC)
type ErrorCode struct {
	Code        string `json:"code"`        // Pxxxx, Bxxxx, Cxxxx, Uxxxx
	Description string `json:"description"` // Human-readable description
	Criticality string `json:"criticality"` // critical, non-critical
	MILStatus   bool   `json:"milStatus"`   // Malfunction Indicator Lamp status
}

