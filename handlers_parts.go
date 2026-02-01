package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func PartsHandler(rp *Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			q := r.URL.Query()

			var catID *primitive.ObjectID
			if v := q.Get("category_id"); v != "" {
				id, err := primitive.ObjectIDFromHex(v)
				if err != nil {
					WriteError(w, 400, "invalid category_id")
					return
				}
				catID = &id
			}

			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			parts, err := rp.ListPartsFiltered(ctx, catID, q.Get("car_model"), q.Get("brand"), q.Get("q"), q.Get("compatibility"))
			if err != nil {
				WriteError(w, 500, "db error")
				return
			}
			WriteJSON(w, 200, parts)

		case http.MethodPost:
			var in struct {
				CategoryID      string  `json:"category_id"`
				Brand           string  `json:"brand"`
				CarModel        string  `json:"car_model"`
				Compatibility   string  `json:"compatibility"`
				Price           float64 `json:"price"`
				Stock           int     `json:"stock"`
				Description     string  `json:"description"`
				ManufactureDate string  `json:"manufacture_date"`
				IsNew           bool    `json:"is_new"`
			}
			if err := ReadJSON(r, &in); err != nil {
				WriteError(w, 400, "invalid json")
				return
			}
			cid, err := primitive.ObjectIDFromHex(in.CategoryID)
			if err != nil {
				WriteError(w, 400, "invalid category_id")
				return
			}
			if strings.TrimSpace(in.Brand) == "" || strings.TrimSpace(in.CarModel) == "" {
				WriteError(w, 400, "brand and car_model are required")
				return
			}
			if in.Price <= 0 {
				WriteError(w, 400, "price must be > 0")
				return
			}
			if in.Stock < 0 {
				WriteError(w, 400, "stock must be >= 0")
				return
			}

			md := time.Time{}
			if in.ManufactureDate != "" {
				tm, err := time.Parse("2006-01-02", in.ManufactureDate)
				if err != nil {
					WriteError(w, 400, "manufacture_date must be YYYY-MM-DD")
					return
				}
				md = tm
			}

			p := SparePart{
				CategoryID:      cid,
				Brand:           in.Brand,
				CarModel:        in.CarModel,
				Compatibility:   in.Compatibility,
				Price:           in.Price,
				Stock:           in.Stock,
				Description:     in.Description,
				ManufactureDate: md,
				IsNew:           in.IsNew,
				IsActive:        true,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			out, err := rp.CreatePart(ctx, p)
			if err != nil {
				WriteError(w, 500, "db error")
				return
			}
			WriteJSON(w, 201, out)

		default:
			WriteError(w, 405, "method not allowed")
		}
	}
}

func PartByIDHandler(rp *Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/parts/")
		if path == "" {
			WriteError(w, 400, "missing id")
			return
		}

		// /parts/{id}/availability
		if strings.HasSuffix(path, "/availability") {
			idStr := strings.TrimSuffix(path, "/availability")
			idStr = strings.TrimSuffix(idStr, "/")
			id, err := primitive.ObjectIDFromHex(idStr)
			if err != nil { WriteError(w, 400, "invalid id"); return }

			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			p, err := rp.GetPart(ctx, id)
			if err != nil {
				if err == mongo.ErrNoDocuments { WriteError(w, 404, "not found"); return }
				WriteError(w, 500, "db error")
				return
			}
			WriteJSON(w, 200, map[string]any{
				"part_id":   p.ID,
				"stock":     p.Stock,
				"available": p.Stock > 0 && p.IsActive,
			})
			return
		}

		idStr := strings.Split(path, "/")[0]
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			WriteError(w, 400, "invalid id")
			return
		}

		switch r.Method {
		case http.MethodGet:
			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			p, err := rp.GetPart(ctx, id)
			if err != nil {
				if err == mongo.ErrNoDocuments { WriteError(w, 404, "not found"); return }
				WriteError(w, 500, "db error")
				return
			}
			WriteJSON(w, 200, p)

		case http.MethodPut:
			var in struct {
				CategoryID    string  `json:"category_id"`
				Brand         string  `json:"brand"`
				CarModel      string  `json:"car_model"`
				Compatibility string  `json:"compatibility"`
				Price         float64 `json:"price"`
				Stock         int     `json:"stock"`
				Description   string  `json:"description"`
				IsNew         bool    `json:"is_new"`
				IsActive      bool    `json:"is_active"`
			}
			if err := ReadJSON(r, &in); err != nil { WriteError(w, 400, "invalid json"); return }

			upd := bson.M{
				"brand": in.Brand,
				"car_model": in.CarModel,
				"compatibility": in.Compatibility,
				"price": in.Price,
				"stock": in.Stock,
				"description": in.Description,
				"is_new": in.IsNew,
				"is_active": in.IsActive,
			}
			if in.CategoryID != "" {
				cid, err := primitive.ObjectIDFromHex(in.CategoryID)
				if err != nil { WriteError(w, 400, "invalid category_id"); return }
				upd["category_id"] = cid
			}

			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			p, err := rp.UpdatePart(ctx, id, upd)
			if err != nil {
				WriteError(w, 500, "db error")
				return
			}
			WriteJSON(w, 200, p)

		case http.MethodPatch:
			var in map[string]any
			if err := ReadJSON(r, &in); err != nil { WriteError(w, 400, "invalid json"); return }

			upd := bson.M{}
			if v, ok := in["brand"]; ok { upd["brand"] = toString(v) }
			if v, ok := in["car_model"]; ok { upd["car_model"] = toString(v) }
			if v, ok := in["compatibility"]; ok { upd["compatibility"] = toString(v) }
			if v, ok := in["description"]; ok { upd["description"] = toString(v) }
			if v, ok := in["is_new"]; ok { if b, ok2 := v.(bool); ok2 { upd["is_new"] = b } }
			if v, ok := in["is_active"]; ok { if b, ok2 := v.(bool); ok2 { upd["is_active"] = b } }
			if v, ok := in["price"]; ok {
				f, _ := strconv.ParseFloat(toString(v), 64)
				if f > 0 { upd["price"] = f }
			}
			if v, ok := in["stock"]; ok {
				i, err := strconv.Atoi(toString(v))
				if err == nil { upd["stock"] = i }
			}
			if v, ok := in["category_id"]; ok {
				cid, err := primitive.ObjectIDFromHex(toString(v))
				if err != nil { WriteError(w, 400, "invalid category_id"); return }
				upd["category_id"] = cid
			}

			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			p, err := rp.UpdatePart(ctx, id, upd)
			if err != nil {
				WriteError(w, 500, "db error")
				return
			}
			WriteJSON(w, 200, p)

		case http.MethodDelete:
			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			if err := rp.DeletePart(ctx, id); err != nil {
				WriteError(w, 500, "db error")
				return
			}
			WriteJSON(w, 200, map[string]string{"deleted": idStr})

		default:
			WriteError(w, 405, "method not allowed")
		}
	}
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		return strconv.Itoa(t)
	default:
		return ""
	}
}
