package main

import (
	"carparts/models"
	"context"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func OrdersHandler(rp *Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			WriteError(w, 405, "method not allowed")
			return
		}

		var in struct {
			CustomerID string `json:"customer_id"`
			Items      []struct {
				PartID   string `json:"part_id"`
				Quantity int    `json:"quantity"`
			} `json:"items"`
		}
		if err := ReadJSON(r, &in); err != nil {
			WriteError(w, 400, "invalid json")
			return
		}
		if strings.TrimSpace(in.CustomerID) == "" {
			WriteError(w, 400, "customer_id is required")
			return
		}
		if len(in.Items) == 0 {
			WriteError(w, 400, "items are required")
			return
		}
		cid, err := primitive.ObjectIDFromHex(in.CustomerID)
		if err != nil {
			WriteError(w, 400, "invalid customer_id")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		defer cancel()

		orderItems := make([]models.OrderItem, 0, len(in.Items))
		for _, it := range in.Items {
			pid, err := primitive.ObjectIDFromHex(it.PartID)
			if err != nil {
				WriteError(w, 400, "invalid part_id")
				return
			}

			updatedPart, err := rp.DecreaseStock(ctx, pid, it.Quantity)
			if err != nil {
				WriteError(w, 400, err.Error())
				return
			}

			orderItems = append(orderItems, models.OrderItem{
				OrderID:  primitive.NilObjectID, // will set after insert
				PartID:   updatedPart.ID,
				Price:    updatedPart.Price,
				Quantity: it.Quantity,
			})
		}

		o := models.Order{
			CustomerID: cid,
			Items:      orderItems,
			IsPaid:     false,
			Status:     "created",
			CreatedAt:  time.Now(),
		}
		o.TotalPrice = o.CalculateTotal()

		created, err := rp.CreateOrder(ctx, o)
		if err != nil {
			WriteError(w, 500, "db error")
			return
		}

		// update embedded items with order_id (simple second update)
		for i := range created.Items {
			created.Items[i].OrderID = created.ID
		}
		_ = rp.syncOrderItems(ctx, created)

		WriteJSON(w, 201, created)
	}
}

func (rp *Repo) syncOrderItems(ctx context.Context, o models.Order) error {
	_, err := rp.orders.UpdateOne(ctx, bson.M{"_id": o.ID}, bson.M{"$set": bson.M{"items": o.Items}})
	return err
}

func OrderByIDHandler(rp *Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/orders/")
		idStr := strings.Split(path, "/")[0]
		if idStr == "" {
			WriteError(w, 400, "missing id")
			return
		}
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			WriteError(w, 400, "invalid id")
			return
		}

		// /orders/{id}/status
		if strings.HasSuffix(r.URL.Path, "/status") {
			if r.Method != http.MethodPatch {
				WriteError(w, 405, "method not allowed")
				return
			}
			var in struct {
				Status string `json:"status"`
				IsPaid bool   `json:"is_paid"`
			}
			if err := ReadJSON(r, &in); err != nil {
				WriteError(w, 400, "invalid json")
				return
			}
			if strings.TrimSpace(in.Status) == "" {
				WriteError(w, 400, "status is required")
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			out, err := rp.UpdateOrderStatus(ctx, id, in.Status, in.IsPaid)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					WriteError(w, 404, "not found")
					return
				}
				WriteError(w, 500, "db error")
				return
			}
			WriteJSON(w, 200, out)
			return
		}

		switch r.Method {
		case http.MethodGet:
			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			out, err := rp.GetOrder(ctx, id)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					WriteError(w, 404, "not found")
					return
				}
				WriteError(w, 500, "db error")
				return
			}
			WriteJSON(w, 200, out)

		case http.MethodDelete:
			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			out, err := rp.CancelOrder(ctx, id)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					WriteError(w, 404, "not found")
					return
				}
				WriteError(w, 500, "db error")
				return
			}
			WriteJSON(w, 200, out)

		default:
			WriteError(w, 405, "method not allowed")
		}
	}
}
