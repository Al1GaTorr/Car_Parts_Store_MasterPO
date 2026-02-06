package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID  primitive.ObjectID `bson:"order_id" json:"order_id"`
	PartID   primitive.ObjectID `bson:"part_id" json:"part_id"`
	Price    float64            `bson:"price" json:"price"`
	Quantity int                `bson:"quantity" json:"quantity"`
}

func (oi *OrderItem) AddItem()    {}
func (oi *OrderItem) RemoveItem() {}
func (oi *OrderItem) UpdateQuantity(q int) {
	if q > 0 {
		oi.Quantity = q
	}
}
