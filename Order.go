package main

import "time"

type Order struct {
	ID         string      `json:"id"`
	CustomerID string      `json:"customer_id"`
	Items      []OrderItem `json:"items"`
	IsServed   bool        `json:"is_served"`
	DateOffice time.Time   `json:"date_office"`
}

func (o *Order) CreateOrder()             {}
func (o *Order) UpdateStatus(status bool) { o.IsServed = status }
func (o *Order) CalculateTotal() float64 {
	var total float64
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}
func (o *Order) Cancel() {}
