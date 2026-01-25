package main

import "time"

type SpareParts struct {
	ID              string    `json:"id"`
	CategoryID      string    `json:"category_id"`
	Compatibility   string    `json:"compatibility"`
	Price           float64   `json:"price"`
	Stock           string    `json:"stock"`
	Description     string    `json:"description"`
	ManufactureDate time.Time `json:"manufacture_date"`
	IsNew           bool      `json:"is_new"`
}

func (s *SpareParts) GetDetails()           {}
func (s *SpareParts) CheckStock()           {}
func (s *SpareParts) UpdatePrice(p float64) { s.Price = p }
func (s *SpareParts) CheckCompatibility()   {}
