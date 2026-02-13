package service

import (
	"context"
	"strings"
	"time"

	"bazarpo-backend/internal/model"
	"bazarpo-backend/internal/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminPartService struct {
	repo *repo.Repository
}

func NewAdminPartService(r *repo.Repository) *AdminPartService {
	return &AdminPartService{repo: r}
}

func (s *AdminPartService) AdminListParts(ctx context.Context) ([]model.PartDoc, error) {
	return s.repo.ListPartsForAdmin(ctx, 1000)
}

func (s *AdminPartService) AdminInsertPart(ctx context.Context, p model.PartDoc) error {
	if strings.TrimSpace(p.SKU) == "" || strings.TrimSpace(p.Name) == "" || strings.TrimSpace(p.Category) == "" {
		return ErrNoFields
	}
	p.SKU = strings.TrimSpace(p.SKU)
	p.Name = strings.TrimSpace(p.Name)
	p.Currency = "KZT"
	p.CreatedAt = time.Now()
	return s.repo.InsertPart(ctx, p)
}

func (s *AdminPartService) AdminUpdatePart(ctx context.Context, id string, patch map[string]any) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrMissingID
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrBadID
	}

	allowed := map[string]bool{
		"name": true, "price_kzt": true, "stock_qty": true, "is_visible": true, "images": true,
		"category": true, "type": true, "brand": true, "compatibility": true, "recommend_for_issue_codes": true,
	}
	set := bson.M{}
	for k, v := range patch {
		if allowed[k] {
			set[k] = v
		}
	}
	if len(set) == 0 {
		return ErrNoFields
	}
	return s.repo.UpdatePartByID(ctx, oid, set)
}

func (s *AdminPartService) AdminDeletePart(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrMissingID
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrBadID
	}
	return s.repo.DeletePartByID(ctx, oid)
}
