package main

type OrderItem struct {
	ID       string  `json:"id"`
	OrderID  string  `json:"order_id"`
	PartID   string  `json:"part_id"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func (oi *OrderItem) AddItem() {
	// logic
}

func (oi *OrderItem) RemoveItem() {
	// logic
}

func (oi *OrderItem) UpdateQuantity(qty int) {
	if qty > 0 {
		oi.Quantity = qty
	}
}
