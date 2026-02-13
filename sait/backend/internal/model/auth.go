package model

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type RegisterRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateOrderRequest struct {
	Items           []OrderItem `json:"items"`
	ShippingAddress string      `json:"shippingAddress"`
	ContactInfo     string      `json:"contactInfo"`
}

