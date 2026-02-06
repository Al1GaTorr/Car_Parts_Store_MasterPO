package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LowStockAlert struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PartID primitive.ObjectID `bson:"part_id" json:"part_id"`
	Name   string             `bson:"name" json:"name"`
	Stock  int                `bson:"stock" json:"stock"`
	At     time.Time          `bson:"at" json:"at"`
}
