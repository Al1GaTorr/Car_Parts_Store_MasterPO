package main

import "net/http"

func RegisterRoutes(mux *http.ServeMux, r *Repo) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/categories", CategoriesHandler(r))
	mux.HandleFunc("/categories/", CategoryByIDHandler(r))

	mux.HandleFunc("/parts", PartsHandler(r))
	mux.HandleFunc("/parts/", PartByIDHandler(r))

	mux.HandleFunc("/vehicle/search", VehicleSearchHandler(r))

	mux.HandleFunc("/orders", OrdersHandler(r))
	mux.HandleFunc("/orders/", OrderByIDHandler(r))

	mux.HandleFunc("/alerts", AlertsHandler(r))
}
