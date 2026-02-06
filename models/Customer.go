package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Customer struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string             `bson:"first_name" json:"first_name"`
	Phone     string             `bson:"phone" json:"phone"`
	Email     string             `bson:"email" json:"email"`
}

func (c *Customer) Register()                 {}
func (c *Customer) Login()                    {}
func (c *Customer) UpdateProfile()            {}
func (c *Customer) ViewOrderHistory() []Order { return nil }
