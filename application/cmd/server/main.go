package main

import (
	"car-monitoring/internal/data"
	"car-monitoring/internal/handlers"
	"car-monitoring/internal/models"
	"car-monitoring/internal/services"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Initialize MockTelemetry map
	data.MockTelemetry = make(map[string]models.Telemetry)

	// Prefer loading vehicles from MongoDB (bazarPO.cars). If that fails, fallback to mock_vehicles.json.
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017/bazarPO"
	}
	if err := data.LoadVehiclesFromMongo(mongoURI); err != nil {
		log.Printf("⚠️  Failed to load vehicles from MongoDB (%s): %v", mongoURI, err)
		log.Println("   Falling back to mock_vehicles.json")
		if err2 := data.LoadMockVehicles(); err2 != nil {
			log.Printf("⚠️  Failed to load mock_vehicles.json: %v", err2)
			log.Println("   Using fallback built-in mock data")
			data.InitializeMockData()
		} else {
			log.Printf("✓ Loaded %d vehicles from mock_vehicles.json", len(data.MockVehicles))
		}
	} else {
		log.Printf("✓ Loaded %d vehicles from MongoDB", len(data.MockVehicles))
	}

	// Load error severity map
	if err := data.LoadErrorSeverityMap(); err != nil {
		log.Printf("⚠️  Failed to load error_severity_map.json: %v", err)
	} else {
		log.Println("✓ Loaded error severity map")
	}

	// Load persisted service history
	if err := data.LoadServiceHistory(); err != nil {
		log.Printf("⚠️  Failed to load service_history.json: %v", err)
	} else {
		log.Println("✓ Loaded persisted service history")
	}

	// Load STO changes summary
	if err := data.LoadSTOChanges(); err != nil {
		log.Printf("⚠️  Failed to load sto_changes.json: %v", err)
	} else {
		log.Println("✓ Loaded STO changes summary")
	}

	// Initialize services
	vehicleService := services.NewVehicleService()
	carService := services.NewCarService()
	analysisService := services.NewAnalysisService(carService)
	dashboardService := services.NewDashboardService(vehicleService, carService, analysisService)
	historyService := services.NewServiceHistoryService(carService)
	notificationService := services.NewNotificationService(vehicleService, analysisService)

	// 2GIS Service (use API key from environment or default key)
	twoGISKey := os.Getenv("TWOGIS_API_KEY")
	if twoGISKey == "" {
		// Default API key for development
		twoGISKey = "60bd214c-bcab-413a-a891-ffba2e40a31f"
	}
	sto2GISService := services.NewSTO2GISService(twoGISKey)

	// Legacy STO service (for backward compatibility)
	stoService := services.NewSTOService()

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	adminHandler := handlers.NewAdminHandler(vehicleService)
	carHandler := handlers.NewCarHandler(carService)
	analysisHandler := handlers.NewAnalysisHandler(analysisService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	historyHandler := handlers.NewServiceHistoryHandler(historyService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	repairShopsHandler := handlers.NewRepairShopsHandler(sto2GISService, vehicleService)
	stoHandler := handlers.NewSTOHandler(stoService)

	mux := http.NewServeMux()

	// CORS middleware for API routes
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// API routes with CORS
	mux.Handle("/api/health", corsMiddleware(healthHandler))
	mux.Handle("/api/admin/", corsMiddleware(adminHandler))
	mux.Handle("/api/dashboard/", corsMiddleware(dashboardHandler))
	mux.Handle("/api/service-history/", corsMiddleware(historyHandler))
	mux.Handle("/api/notifications/", corsMiddleware(notificationHandler))
	mux.Handle("/api/repair-shops", corsMiddleware(repairShopsHandler))

	// Legacy car endpoints
	mux.HandleFunc("/api/cars/", func(w http.ResponseWriter, r *http.Request) {
		corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/analysis") {
				analysisHandler.ServeHTTP(w, r)
				return
			}
			carHandler.ServeHTTP(w, r)
		})).ServeHTTP(w, r)
	})

	// Legacy STO endpoint
	mux.Handle("/api/sto", corsMiddleware(stoHandler))

	// Real-time SSE endpoint per VIN
	mux.HandleFunc("/ws/cars/", func(w http.ResponseWriter, r *http.Request) {
		corsMiddleware(http.HandlerFunc(handlers.SSEHandler)).ServeHTTP(w, r)
	})

	// Telemetry POST endpoint
	mux.Handle("/api/telemetry", corsMiddleware(http.HandlerFunc(handlers.TelemetryHandler)))

	// STO panel API (vehicles list, update, add records)
	mux.Handle("/api/sto-panel/", corsMiddleware(http.HandlerFunc(handlers.STOPanelHandler)))

	// Serve React frontend from front/dist
	// Handle SPA routing - all non-API routes serve index.html
	frontendDir := "./front/dist"

	// Check if frontend is built
	_, frontendExists := os.Stat(frontendDir)
	if os.IsNotExist(frontendExists) {
		log.Printf("⚠️  Frontend not built. Run 'cd front && npm run build' to build")
		log.Printf("   For development, run frontend separately: 'cd front && npm run dev'")
		log.Printf("   Backend API is available at http://localhost:8081/api")

		// Serve a helpful message for root path
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				return // API routes handled above
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>Frontend Not Built</title></head>
<body style="font-family: Arial, sans-serif; padding: 40px; text-align: center; max-width: 600px; margin: 0 auto;">
	<h1>Frontend Not Built</h1>
	<p>Please build the frontend first:</p>
	<pre style="background: #f0f0f0; padding: 20px; border-radius: 5px; display: inline-block;">cd front && npm run build</pre>
	<p>Or run in development mode:</p>
	<pre style="background: #f0f0f0; padding: 20px; border-radius: 5px; display: inline-block;">cd front && npm run dev</pre>
	<p>API is available at <a href="/api/health">/api/health</a></p>
</body>
</html>`))
			if err != nil {
				return
			}
		})
	} else {
		// Serve static files
		fs := http.FileServer(http.Dir(frontendDir))

		// SPA routing: serve index.html for all non-API routes
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Don't handle API routes here - they're already handled above
			if strings.HasPrefix(r.URL.Path, "/api/") {
				return
			}

			// Root path or empty - serve index.html
			if r.URL.Path == "/" || r.URL.Path == "" {
				indexPath := frontendDir + "/index.html"
				http.ServeFile(w, r, indexPath)
				return
			}

			// Check if it's a request for a static asset (has file extension)
			// Try to serve the file if it exists
			fullPath := frontendDir + r.URL.Path
			if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
				// File exists, serve it
				fs.ServeHTTP(w, r)
				return
			}

			// For SPA routing, all other routes serve index.html
			// This allows React Router to handle client-side routing
			indexPath := frontendDir + "/index.html"
			if _, err := os.Stat(indexPath); err == nil {
				http.ServeFile(w, r, indexPath)
			} else {
				http.NotFound(w, r)
			}
		})
		log.Printf("✓ Serving React frontend from ./front/dist")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Server started on http://localhost:%s", port)
	log.Println("API endpoints:")
	log.Println("  GET /api/health")
	log.Println("  GET /api/dashboard/{carId}")
	log.Println("  GET /api/service-history/{carId}")
	log.Println("  GET /api/notifications/{carId}")
	log.Println("  GET /api/repair-shops?lat=43.2220&lon=76.8512&radius=10")
	log.Println("  GET /api/cars")
	log.Println("  GET /api/cars/{id}")
	log.Println("  GET /api/cars/{id}/telemetry")
	log.Println("  GET /api/cars/{id}/analysis")

	if twoGISKey == "" {
		log.Println("⚠️  2GIS API key not set (TWOGIS_API_KEY), using mock data")
	} else {
		log.Println("✓ 2GIS API integration enabled")
	}

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
