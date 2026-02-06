package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SparePart struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CategoryID      primitive.ObjectID `bson:"category_id" json:"category_id"`
	Brand           string             `bson:"brand" json:"brand"`
	CarModel        string             `bson:"car_model" json:"car_model"`
	Compatibility   string             `bson:"compatibility" json:"compatibility"`
	Price           float64            `bson:"price" json:"price"`
	Stock           int                `bson:"stock" json:"stock"`
	Description     string             `bson:"description" json:"description"`
	ManufactureDate time.Time          `bson:"manufacture_date" json:"manufacture_date"`
	IsNew           bool               `bson:"is_new" json:"is_new"`
	IsActive        bool               `bson:"is_active" json:"is_active"`
}

func (s *SparePart) GetDetails()              {}
func (s *SparePart) CheckStock() bool         { return s.Stock > 0 && s.IsActive }
func (s *SparePart) UpdatePrice(p float64)    { s.Price = p }
func (s *SparePart) CheckCompatibility() bool { return s.Compatibility != "" }
