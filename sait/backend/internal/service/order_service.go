package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"bazarpo-backend/internal/model"
	"bazarpo-backend/internal/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrEmptyItems = errors.New("empty items")
	ErrMissingID  = errors.New("missing id")
	ErrBadID      = errors.New("bad id")
	ErrNoFields   = errors.New("no fields")
)

type StockIssue struct {
	SKU       string `json:"sku"`
	Requested int    `json:"requested"`
	Available int    `json:"available"`
}

type InsufficientStockError struct {
	Issues []StockIssue `json:"issues"`
}

func (e *InsufficientStockError) Error() string {
	return "insufficient stock"
}

type OrderService struct {
	repo *repo.Repository
}

func NewOrderService(r *repo.Repository) *OrderService {
	return &OrderService{repo: r}
}

func (s *OrderService) CreateOrder(ctx context.Context, userHex string, req model.CreateOrderRequest) (primitive.ObjectID, error) {
	if len(req.Items) == 0 {
		return primitive.NilObjectID, ErrEmptyItems
	}
	userID, err := primitive.ObjectIDFromHex(userHex)
	if err != nil {
		return primitive.NilObjectID, ErrUnauthorized
	}

	requestedBySKU := map[string]int{}
	sanitizedItems := make([]model.OrderItem, 0, len(req.Items))
	for _, it := range req.Items {
		sku := strings.TrimSpace(it.SKU)
		if sku == "" || it.Qty <= 0 {
			continue
		}
		it.SKU = sku
		sanitizedItems = append(sanitizedItems, it)
		requestedBySKU[sku] += it.Qty
	}
	if len(sanitizedItems) == 0 {
		return primitive.NilObjectID, ErrEmptyItems
	}

	skus := make([]string, 0, len(requestedBySKU))
	for sku := range requestedBySKU {
		skus = append(skus, sku)
	}
	parts, err := s.repo.FindPartsBySKUs(ctx, skus)
	if err != nil {
		return primitive.NilObjectID, err
	}

	partBySKU := map[string]model.PartDoc{}
	availableBySKU := map[string]int{}
	for _, p := range parts {
		sku := strings.TrimSpace(p.SKU)
		partBySKU[sku] = p
		availableBySKU[sku] = p.StockQty
	}
	issues := make([]StockIssue, 0)
	for sku, qty := range requestedBySKU {
		available := availableBySKU[sku]
		if qty > available {
			issues = append(issues, StockIssue{
				SKU:       sku,
				Requested: qty,
				Available: available,
			})
		}
	}
	if len(issues) > 0 {
		return primitive.NilObjectID, &InsufficientStockError{Issues: issues}
	}

	// Trust price from DB, not from client payload.
	total := 0
	pricedItems := make([]model.OrderItem, 0, len(sanitizedItems))
	for _, it := range sanitizedItems {
		part, ok := partBySKU[it.SKU]
		if !ok {
			continue
		}
		it.PriceKZT = part.PriceKZT
		if strings.TrimSpace(it.Name) == "" {
			it.Name = part.Name
		}
		total += it.PriceKZT * it.Qty
		pricedItems = append(pricedItems, it)
	}
	if len(pricedItems) == 0 {
		return primitive.NilObjectID, ErrEmptyItems
	}

	reserved := map[string]int{}
	for sku, qty := range requestedBySKU {
		ok, reserveErr := s.repo.DecrementPartStockIfEnough(ctx, sku, qty)
		if reserveErr != nil {
			for rSKU, rQty := range reserved {
				_ = s.repo.IncrementPartStock(ctx, rSKU, rQty)
			}
			return primitive.NilObjectID, reserveErr
		}
		if !ok {
			for rSKU, rQty := range reserved {
				_ = s.repo.IncrementPartStock(ctx, rSKU, rQty)
			}
			part, partErr := s.repo.FindPartBySKU(ctx, sku)
			available := 0
			if partErr == nil {
				available = part.StockQty
			}
			return primitive.NilObjectID, &InsufficientStockError{
				Issues: []StockIssue{{
					SKU:       sku,
					Requested: qty,
					Available: available,
				}},
			}
		}
		reserved[sku] = qty
	}

	oid, err := s.repo.InsertOrder(ctx, model.OrderDoc{
		UserID:          userID,
		Items:           pricedItems,
		TotalKZT:        total,
		Status:          "pending",
		ShippingAddress: req.ShippingAddress,
		ContactInfo:     req.ContactInfo,
		CreatedAt:       time.Now(),
	})
	if err != nil {
		for sku, qty := range reserved {
			_ = s.repo.IncrementPartStock(ctx, sku, qty)
		}
		return primitive.NilObjectID, err
	}

	return oid, nil
}

func (s *OrderService) AdminListOrders(ctx context.Context) ([]model.OrderDoc, error) {
	return s.repo.ListOrders(ctx, 200)
}

func (s *OrderService) AdminUpdateOrder(ctx context.Context, id string, patch map[string]any) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrMissingID
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrBadID
	}

	set := bson.M{}
	if v, ok := patch["status"].(string); ok {
		status := strings.ToLower(strings.TrimSpace(v))
		switch status {
		case "pending", "processing", "shipped", "completed", "cancelled":
			set["status"] = status
		}
	}
	if len(set) == 0 {
		return ErrNoFields
	}
	return s.repo.UpdateOrderByID(ctx, oid, set)
}

func (s *OrderService) AdminDeleteOrder(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrMissingID
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrBadID
	}
	return s.repo.DeleteOrderByID(ctx, oid)
}
