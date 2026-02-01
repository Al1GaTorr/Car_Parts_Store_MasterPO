package main

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GET /vehicle/search?car_model=&brand=&compatibility=&category_id=&q=
func VehicleSearchHandler(rp *Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			WriteError(w, 405, "method not allowed")
			return
		}

		q := r.URL.Query()
		var catID *primitive.ObjectID
		if v := q.Get("category_id"); v != "" {
			id, err := primitive.ObjectIDFromHex(v)
			if err != nil { WriteError(w, 400, "invalid category_id"); return }
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
	}
}
