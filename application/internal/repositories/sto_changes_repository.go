package repositories

import "car-monitoring/internal/data"

type STOChangesRepository struct{}

func NewSTOChangesRepository() *STOChangesRepository {
	return &STOChangesRepository{}
}

func (r *STOChangesRepository) Get(vehicleID string) map[string]interface{} {
	return data.GetSTOChange(vehicleID)
}

func (r *STOChangesRepository) Update(vehicleID string, change map[string]interface{}) error {
	return data.UpdateSTOChange(vehicleID, change)
}
