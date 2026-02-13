package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bazarpo-backend/internal/model"
	"bazarpo-backend/internal/service"
)

type Handler struct {
	auth   *service.AuthService
	cars   *service.CarService
	parts  *service.PartService
	orders *service.OrderService
	admin  *service.AdminPartService
}

func New(auth *service.AuthService, cars *service.CarService, parts *service.PartService, orders *service.OrderService, admin *service.AdminPartService) *Handler {
	return &Handler{
		auth:   auth,
		cars:   cars,
		parts:  parts,
		orders: orders,
		admin:  admin,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func readJSON(r *http.Request, out any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(out)
}

func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func partIDFromPath(path string) string {
	if strings.HasPrefix(path, "/api/admin/parts/") {
		return strings.TrimPrefix(path, "/api/admin/parts/")
	}
	if strings.HasPrefix(path, "/api/admin/parts/") {
		return strings.TrimPrefix(path, "/api/admin/parts/")
	}
	return path
}

func (h *Handler) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) requireAuth(next func(http.ResponseWriter, *http.Request, *model.Claims)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := h.auth.ParseToken(r.Header.Get("Authorization"))
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
			return
		}
		next(w, r, claims)
	}
}

func (h *Handler) requireAdmin(next func(http.ResponseWriter, *http.Request, *model.Claims)) http.HandlerFunc {
	return h.requireAuth(func(w http.ResponseWriter, r *http.Request, c *model.Claims) {
		if c.Role != "admin" {
			writeJSON(w, http.StatusForbidden, map[string]any{"error": "forbidden"})
			return
		}
		next(w, r, c)
	})
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{"ok": true})
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, 400, map[string]any{"error": "invalid body"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	token, role, err := h.auth.Register(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailPasswordInvalid):
			writeJSON(w, 400, map[string]any{"error": "email/password invalid"})
		case errors.Is(err, service.ErrEmailAlreadyExists):
			writeJSON(w, 409, map[string]any{"error": "email already exists"})
		default:
			writeJSON(w, 500, map[string]any{"error": "db error"})
		}
		return
	}
	writeJSON(w, 201, map[string]any{"token": token, "role": role})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, 400, map[string]any{"error": "invalid body"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	token, role, err := h.auth.Login(ctx, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeJSON(w, 401, map[string]any{"error": "invalid credentials"})
			return
		}
		writeJSON(w, 500, map[string]any{"error": "db error"})
		return
	}
	writeJSON(w, 200, map[string]any{"token": token, "role": role})
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request, c *model.Claims) {
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	u, err := h.auth.Me(ctx, c.UserID)
	if err != nil {
		writeJSON(w, 401, map[string]any{"error": "unauthorized"})
		return
	}
	writeJSON(w, 200, map[string]any{
		"id":        u.ID.Hex(),
		"email":     u.Email,
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"role":      u.Role,
	})
}

func (h *Handler) listParts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	year := 0
	if y := strings.TrimSpace(q.Get("year")); y != "" {
		if parsed, err := strconv.Atoi(y); err == nil && parsed > 0 {
			year = parsed
		}
	}
	params := service.ListPartsParams{
		VIN:      q.Get("vin"),
		Issue:    q.Get("issue"),
		Search:   q.Get("q"),
		Make:     q.Get("make"),
		Model:    q.Get("model"),
		Year:     year,
		Category: q.Get("category"),
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	items, err := h.parts.ListParts(ctx, params)
	if err != nil {
		writeJSON(w, 500, map[string]any{"error": "db error"})
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (h *Handler) listCarMakes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	items, err := h.cars.ListMakes(ctx)
	if err != nil {
		writeJSON(w, 500, map[string]any{"error": "db error"})
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (h *Handler) listCarModelsByMake(w http.ResponseWriter, r *http.Request) {
	make := strings.TrimSpace(r.URL.Query().Get("make"))
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	items, err := h.cars.ListModelsByMake(ctx, make)
	if err != nil {
		if errors.Is(err, service.ErrMakeRequired) {
			writeJSON(w, 400, map[string]any{"error": "make required"})
			return
		}
		writeJSON(w, 500, map[string]any{"error": "db error"})
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (h *Handler) listCarYearsByMakeModel(w http.ResponseWriter, r *http.Request) {
	make := strings.TrimSpace(r.URL.Query().Get("make"))
	model := strings.TrimSpace(r.URL.Query().Get("model"))
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	items, err := h.cars.ListYearsByMakeModel(ctx, make, model)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMakeRequired):
			writeJSON(w, 400, map[string]any{"error": "make required"})
		case errors.Is(err, service.ErrModelRequired):
			writeJSON(w, 400, map[string]any{"error": "model required"})
		default:
			writeJSON(w, 500, map[string]any{"error": "db error"})
		}
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request, c *model.Claims) {
	var req model.CreateOrderRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, 400, map[string]any{"error": "invalid body"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()

	oid, err := h.orders.CreateOrder(ctx, c.UserID, req)
	if err != nil {
		if errors.Is(err, service.ErrEmptyItems) {
			writeJSON(w, 400, map[string]any{"error": "empty items"})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			writeJSON(w, 401, map[string]any{"error": "unauthorized"})
			return
		}
		var stockErr *service.InsufficientStockError
		if errors.As(err, &stockErr) {
			writeJSON(w, 409, map[string]any{
				"error":  "insufficient stock",
				"issues": stockErr.Issues,
			})
			return
		}
		writeJSON(w, 500, map[string]any{"error": "db error"})
		return
	}
	writeJSON(w, 201, map[string]any{"id": oid.Hex()})
}

func (h *Handler) adminListOrders(w http.ResponseWriter, r *http.Request, _ *model.Claims) {
	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()
	items, err := h.orders.AdminListOrders(ctx)
	if err != nil {
		writeJSON(w, 500, map[string]any{"error": "db error"})
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (h *Handler) adminUpdateOrder(w http.ResponseWriter, r *http.Request, _ *model.Claims) {
	id := strings.TrimPrefix(r.URL.Path, "/api/admin/orders/")
	var patch map[string]any
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeJSON(w, 400, map[string]any{"error": "invalid body"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	err := h.orders.AdminUpdateOrder(ctx, id, patch)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMissingID):
			writeJSON(w, 400, map[string]any{"error": "missing id"})
		case errors.Is(err, service.ErrBadID):
			writeJSON(w, 400, map[string]any{"error": "bad id"})
		case errors.Is(err, service.ErrNoFields):
			writeJSON(w, 400, map[string]any{"error": "no fields"})
		default:
			writeJSON(w, 500, map[string]any{"error": "db error"})
		}
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true})
}

func (h *Handler) adminDeleteOrder(w http.ResponseWriter, r *http.Request, _ *model.Claims) {
	id := strings.TrimPrefix(r.URL.Path, "/api/admin/orders/")
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := h.orders.AdminDeleteOrder(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMissingID):
			writeJSON(w, 400, map[string]any{"error": "missing id"})
		case errors.Is(err, service.ErrBadID):
			writeJSON(w, 400, map[string]any{"error": "bad id"})
		default:
			writeJSON(w, 500, map[string]any{"error": "db error"})
		}
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true})
}

func (h *Handler) adminListParts(w http.ResponseWriter, r *http.Request, _ *model.Claims) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	items, err := h.admin.AdminListParts(ctx)
	if err != nil {
		writeJSON(w, 500, map[string]any{"error": "db error"})
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (h *Handler) adminInsertPart(w http.ResponseWriter, r *http.Request, _ *model.Claims) {
	var p model.PartDoc
	if err := readJSON(r, &p); err != nil {
		writeJSON(w, 400, map[string]any{"error": "invalid body"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	if err := h.admin.AdminInsertPart(ctx, p); err != nil {
		if errors.Is(err, service.ErrNoFields) {
			writeJSON(w, 400, map[string]any{"error": "sku/name/category required"})
			return
		}
		writeJSON(w, 500, map[string]any{"error": "db error"})
		return
	}
	writeJSON(w, 201, map[string]any{"ok": true})
}

func (h *Handler) adminUpdatePart(w http.ResponseWriter, r *http.Request, _ *model.Claims) {
	id := partIDFromPath(r.URL.Path)
	var patch map[string]any
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeJSON(w, 400, map[string]any{"error": "invalid body"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	err := h.admin.AdminUpdatePart(ctx, id, patch)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMissingID):
			writeJSON(w, 400, map[string]any{"error": "missing id"})
		case errors.Is(err, service.ErrBadID):
			writeJSON(w, 400, map[string]any{"error": "bad id"})
		case errors.Is(err, service.ErrNoFields):
			writeJSON(w, 400, map[string]any{"error": "no fields"})
		default:
			writeJSON(w, 500, map[string]any{"error": "db error"})
		}
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true})
}

func (h *Handler) adminDeletePart(w http.ResponseWriter, r *http.Request, _ *model.Claims) {
	id := partIDFromPath(r.URL.Path)
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	err := h.admin.AdminDeletePart(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrMissingID) {
			writeJSON(w, 400, map[string]any{"error": "missing id"})
			return
		}
		if errors.Is(err, service.ErrBadID) {
			writeJSON(w, 400, map[string]any{"error": "bad id"})
			return
		}
		writeJSON(w, 500, map[string]any{"error": "db error"})
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true})
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", h.health)
	mux.HandleFunc("/api/parts", h.listParts)
	mux.HandleFunc("/api/cars/makes", func(w http.ResponseWriter, r *http.Request) {
		if !requireMethod(w, r, http.MethodGet) {
			return
		}
		h.listCarMakes(w, r)
	})
	mux.HandleFunc("/api/cars/models", func(w http.ResponseWriter, r *http.Request) {
		if !requireMethod(w, r, http.MethodGet) {
			return
		}
		h.listCarModelsByMake(w, r)
	})
	mux.HandleFunc("/api/cars/years", func(w http.ResponseWriter, r *http.Request) {
		if !requireMethod(w, r, http.MethodGet) {
			return
		}
		h.listCarYearsByMakeModel(w, r)
	})

	mux.HandleFunc("/api/auth/register", h.register)
	mux.HandleFunc("/api/auth/login", h.login)
	mux.HandleFunc("/api/auth/me", h.requireAuth(h.me))

	mux.HandleFunc("/api/orders", h.requireAuth(func(w http.ResponseWriter, r *http.Request, c *model.Claims) {
		if !requireMethod(w, r, http.MethodPost) {
			return
		}
		h.createOrder(w, r, c)
	}))

	mux.HandleFunc("/api/admin/orders", h.requireAdmin(func(w http.ResponseWriter, r *http.Request, c *model.Claims) {
		if !requireMethod(w, r, http.MethodGet) {
			return
		}
		h.adminListOrders(w, r, c)
	}))
	mux.HandleFunc("/api/admin/orders/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch || r.Method == http.MethodPut {
			h.requireAdmin(h.adminUpdateOrder)(w, r)
			return
		}
		if r.Method == http.MethodDelete {
			h.requireAdmin(h.adminDeleteOrder)(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	// Preferred naming.
	mux.HandleFunc("/api/admin/parts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.requireAdmin(h.adminListParts)(w, r)
			return
		}
		if r.Method == http.MethodPost {
			h.requireAdmin(h.adminInsertPart)(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/api/admin/parts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch || r.Method == http.MethodPut {
			h.requireAdmin(h.adminUpdatePart)(w, r)
			return
		}
		if r.Method == http.MethodDelete {
			h.requireAdmin(h.adminDeletePart)(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	return h.withCORS(mux)
}
