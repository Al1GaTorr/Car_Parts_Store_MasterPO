package main

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

func AlertsHandler(rp *Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			WriteError(w, 405, "method not allowed")
			return
		}

		limit := int64(0)
		if v := r.URL.Query().Get("limit"); v != "" {
			n, err := strconv.Atoi(v)
			if err == nil && n > 0 {
				limit = int64(n)
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		alerts, err := rp.ListAlerts(ctx, limit)
		if err != nil {
			WriteError(w, 500, "db error")
			return
		}
		WriteJSON(w, 200, alerts)
	}
}
