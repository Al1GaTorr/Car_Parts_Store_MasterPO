package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CategoriesHandler(rp *Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			cats, err := rp.ListCategories(ctx)
			if err != nil {
				WriteError(w, 500, "db error")
				return
			}
			WriteJSON(w, 200, cats)

		case http.MethodPost:
			var in struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			}
			if err := ReadJSON(r, &in); err != nil {
				WriteError(w, 400, "invalid json")
				return
			}
			if strings.TrimSpace(in.Name) == "" {
				WriteError(w, 400, "name is required")
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			c := Category{Name: in.Name, Description: in.Description, PartsList: []primitive.ObjectID{}}
			out, err := rp.CreateCategory(ctx, c)
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

func CategoryByIDHandler(rp *Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
		if idStr == "" {
			WriteError(w, 400, "missing id")
			return
		}
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			WriteError(w, 400, "invalid id")
			return
		}

		switch r.Method {
		case http.MethodGet:
			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			c, err := rp.GetCategory(ctx, id)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					WriteError(w, 404, "not found")
					return
				}
				WriteError(w, 500, "db error")
				return
			}
			WriteJSON(w, 200, c)

		case http.MethodPut:
			var in struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			}
			if err := ReadJSON(r, &in); err != nil {
				WriteError(w, 400, "invalid json")
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			c, err := rp.UpdateCategory(ctx, id, in.Name, in.Description)
			if err != nil {
				WriteError(w, 500, "db error")
				return
			}
			WriteJSON(w, 200, c)

		case http.MethodDelete:
			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			if err := rp.DeleteCategory(ctx, id); err != nil {
				WriteError(w, 500, "db error")
				return
			}
			WriteJSON(w, 200, map[string]string{"deleted": idStr})

		default:
			WriteError(w, 405, "method not allowed")
		}
	}
}
