package service

import (
	"context"
	"errors"
	"strings"

	"bazarpo-backend/internal/repo"
)

var (
	ErrMakeRequired  = errors.New("make required")
	ErrModelRequired = errors.New("model required")
)

type CarService struct {
	repo *repo.Repository
}

func NewCarService(r *repo.Repository) *CarService {
	return &CarService{repo: r}
}

func (s *CarService) ListMakes(ctx context.Context) ([]string, error) {
	return s.repo.DistinctCarMakes(ctx)
}

func (s *CarService) ListModelsByMake(ctx context.Context, make string) ([]string, error) {
	if strings.TrimSpace(make) == "" {
		return nil, ErrMakeRequired
	}
	return s.repo.DistinctCarModelsByMake(ctx, make)
}

func (s *CarService) ListYearsByMakeModel(ctx context.Context, make, model string) ([]int, error) {
	if strings.TrimSpace(make) == "" {
		return nil, ErrMakeRequired
	}
	if strings.TrimSpace(model) == "" {
		return nil, ErrModelRequired
	}
	return s.repo.DistinctCarYearsByMakeModel(ctx, make, model)
}
