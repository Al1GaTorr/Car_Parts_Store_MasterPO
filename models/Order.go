package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerID primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	Items      []OrderItem        `bson:"items" json:"items"`
	IsPaid     bool               `bson:"is_paid" json:"is_paid"`
	TotalPrice float64            `bson:"total_price" json:"total_price"`
	Status     string             `bson:"status" json:"status"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

func (o *Order) CreateOrder()               {}
func (o *Order) UpdateStatus(status string) { o.Status = status }
func (o *Order) CalculateTotal() float64 {
	var t float64
	for _, it := range o.Items {
		t += it.Price * float64(it.Quantity)
	}
	return t
}
func (o *Order) Cancel() { o.Status = "canceled" }
