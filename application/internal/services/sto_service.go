package services

import (
	"car-monitoring/internal/models"
	"car-monitoring/internal/repositories"
	"strconv"
	"strings"
)

// STOService handles business logic for service stations
type STOService struct {
	stationRepo *repositories.ServiceStationRepository
}

// NewSTOService creates a new STO service instance
func NewSTOService() *STOService {
	return &STOService{
		stationRepo: repositories.NewServiceStationRepository(),
	}
}

// GetServiceStations returns filtered service stations based on query parameters
func (s *STOService) GetServiceStations(filters map[string]string) []models.ServiceStation {
	stations := s.stationRepo.GetAll()
	filtered := []models.ServiceStation{}

	for _, station := range stations {
		if s.matchesFilters(station, filters) {
			filtered = append(filtered, station)
		}
	}

	return filtered
}

// matchesFilters checks if a station matches all provided filters
func (s *STOService) matchesFilters(station models.ServiceStation, filters map[string]string) bool {
	// Filter by rating (e.g., rating>=4)
	if ratingFilter, ok := filters["rating"]; ok {
		if strings.Contains(ratingFilter, ">=") {
			minRating, err := strconv.ParseFloat(strings.TrimPrefix(ratingFilter, ">="), 64)
			if err == nil && station.Rating < minRating {
				return false
			}
		} else if strings.Contains(ratingFilter, ">") {
			minRating, err := strconv.ParseFloat(strings.TrimPrefix(ratingFilter, ">"), 64)
			if err == nil && station.Rating <= minRating {
				return false
			}
		} else {
			// Exact match
			rating, err := strconv.ParseFloat(ratingFilter, 64)
			if err == nil && station.Rating != rating {
				return false
			}
		}
	}

	// Filter by price level (cheap, medium, expensive)
	if priceLevel, ok := filters["price"]; ok {
		if station.PriceLevel != priceLevel {
			return false
		}
	}

	// Filter by supported brand
	if brand, ok := filters["brand"]; ok {
		brandFound := false
		for _, supportedBrand := range station.SupportedBrands {
			if strings.EqualFold(supportedBrand, brand) {
				brandFound = true
				break
			}
		}
		if !brandFound {
			return false
		}
	}

	// Filter by specialization
	if specialization, ok := filters["specialization"]; ok {
		specFound := false
		for _, spec := range station.Specializations {
			if strings.EqualFold(spec, specialization) {
				specFound = true
				break
			}
		}
		if !specFound {
			return false
		}
	}

	return true
}
