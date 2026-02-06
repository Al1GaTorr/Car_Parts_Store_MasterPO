package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Warehouse struct {
	ID    primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Parts []primitive.ObjectID `bson:"parts" json:"parts"`
}

func (w *Warehouse) CheckInventory() []primitive.ObjectID { return w.Parts }
func (w *Warehouse) UpdateStock(partID string)            {}
func (w *Warehouse) GenerateReport()                      {}
func (w *Warehouse) AlertLowStock()                       {}
