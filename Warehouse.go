package main

type Warehouse struct {
	ID    string   `json:"id"`
	Parts []string `json:"parts"`
}

func (w *Warehouse) CheckInventory() []string {
	return w.Parts
}

func (w *Warehouse) UpdateStock(partID string) {
	// logic to update stock of a specific spare part
}

func (w *Warehouse) GenerateReport() {
	// logic to generate inventory report
}

func (w *Warehouse) AlertLowStock() {
	// logic to notify about low stock items
}
