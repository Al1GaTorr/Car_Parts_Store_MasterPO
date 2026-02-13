package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDoc struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName    string             `bson:"first_name" json:"firstName"`
	LastName     string             `bson:"last_name" json:"lastName"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Role         string             `bson:"role" json:"role"` // user|admin
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
}

type CarDoc struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	VIN       string             `bson:"vin" json:"vin"`
	Make      string             `bson:"make" json:"make"`
	Model     string             `bson:"model" json:"model"`
	Year      int                `bson:"year" json:"year"`
	Engine    string             `bson:"engine,omitempty" json:"engine,omitempty"`
	Fuel      string             `bson:"fuel,omitempty" json:"fuel,omitempty"`
	MileageKM int                `bson:"mileage_km,omitempty" json:"mileage_km,omitempty"`
}

type VehicleFit struct {
	Make     string `bson:"make" json:"make"`
	Model    string `bson:"model" json:"model"`
	YearFrom int    `bson:"year_from" json:"year_from"`
	YearTo   int    `bson:"year_to" json:"year_to"`
}

type Compatibility struct {
	Type     string       `bson:"type" json:"type"` // universal|vin|vehicle
	VINs     []string     `bson:"vins,omitempty" json:"vins,omitempty"`
	Vehicles []VehicleFit `bson:"vehicles,omitempty" json:"vehicles,omitempty"`
}

type PartDoc struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	SKU       string             `bson:"sku" json:"sku"`
	Name      string             `bson:"name" json:"name"`
	Category  string             `bson:"category" json:"category"` // Russian category label
	Type      string             `bson:"type,omitempty" json:"type,omitempty"`
	Brand     string             `bson:"brand,omitempty" json:"brand,omitempty"`
	PriceKZT  int                `bson:"price_kzt" json:"price_kzt"`
	Currency  string             `bson:"currency" json:"currency"`
	StockQty  int                `bson:"stock_qty" json:"stock_qty"`
	IsVisible bool               `bson:"is_visible" json:"is_visible"`
	Images    []string           `bson:"images,omitempty" json:"images,omitempty"`
	Compat    Compatibility      `bson:"compatibility" json:"compatibility"`
	Recommend []string           `bson:"recommend_for_issue_codes,omitempty" json:"recommend_for_issue_codes,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type OrderItem struct {
	SKU      string `bson:"sku" json:"sku"`
	Name     string `bson:"name" json:"name"`
	PriceKZT int    `bson:"price_kzt" json:"price_kzt"`
	Qty      int    `bson:"qty" json:"qty"`
	Image    string `bson:"image,omitempty" json:"image,omitempty"`
}

type OrderDoc struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"user_id" json:"userId"`
	Items           []OrderItem        `bson:"items" json:"items"`
	TotalKZT        int                `bson:"total_kzt" json:"total_kzt"`
	Status          string             `bson:"status" json:"status"`
	ShippingAddress string             `bson:"shipping_address" json:"shippingAddress"`
	ContactInfo     string             `bson:"contact_info" json:"contactInfo"`
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt"`
}

