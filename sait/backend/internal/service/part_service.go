package service

import (
	"context"
	"errors"
	"strings"

	"bazarpo-backend/internal/model"
	"bazarpo-backend/internal/repo"
	"go.mongodb.org/mongo-driver/bson"
)

type PartService struct {
	repo *repo.Repository
}

type ListPartsParams struct {
	VIN      string
	Issue    string
	Search   string
	Make     string
	Model    string
	Year     int
	Category string
}

func NewPartService(r *repo.Repository) *PartService {
	return &PartService{repo: r}
}

func (s *PartService) findCarByVIN(ctx context.Context, vin string) (*model.CarDoc, error) {
	vin = strings.ToUpper(strings.TrimSpace(vin))
	if vin == "" {
		return nil, errors.New("empty vin")
	}
	return s.repo.FindCarByVIN(ctx, vin)
}

func matchesVehicle(p model.PartDoc, make, model string, year int, vin string) bool {
	make = strings.TrimSpace(make)
	model = strings.TrimSpace(model)
	vin = strings.ToUpper(strings.TrimSpace(vin))

	switch strings.ToLower(p.Compat.Type) {
	case "universal":
		return true
	case "vin":
		for _, v := range p.Compat.VINs {
			if strings.ToUpper(strings.TrimSpace(v)) == vin && vin != "" {
				return true
			}
		}
		return false
	case "vehicle":
		for _, vf := range p.Compat.Vehicles {
			if strings.EqualFold(vf.Make, make) && strings.EqualFold(vf.Model, model) {
				if year == 0 {
					return true
				}
				if year >= vf.YearFrom && year <= vf.YearTo {
					return true
				}
			}
		}
		return false
	default:
		return true
	}
}

func (s *PartService) ListParts(ctx context.Context, params ListPartsParams) ([]model.PartDoc, error) {
	vin := strings.ToUpper(strings.TrimSpace(params.VIN))
	issue := strings.TrimSpace(params.Issue)
	search := strings.TrimSpace(params.Search)
	carMake := strings.TrimSpace(params.Make)
	carModel := strings.TrimSpace(params.Model)
	year := params.Year
	category := strings.TrimSpace(params.Category)

	if vin != "" {
		car, err := s.findCarByVIN(ctx, vin)
		if err == nil {
			carMake, carModel, year = car.Make, car.Model, car.Year
		}
	}

	filter := bson.M{
		"is_visible": true,
		"stock_qty":  bson.M{"$gt": 0},
	}
	if category != "" {
		filter["category"] = category
	}
	if issue != "" {
		filter["recommend_for_issue_codes"] = issue
	}
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"sku": bson.M{"$regex": search, "$options": "i"}},
			{"brand": bson.M{"$regex": search, "$options": "i"}},
			{"type": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	items, err := s.repo.FindParts(ctx, filter, 500)
	if err != nil {
		return nil, err
	}

	if (carMake != "" && carModel != "") || vin != "" {
		out := make([]model.PartDoc, 0, len(items))
		for _, p := range items {
			if matchesVehicle(p, carMake, carModel, year, vin) {
				out = append(out, p)
			}
		}
		items = out
	}

	return items, nil
}
