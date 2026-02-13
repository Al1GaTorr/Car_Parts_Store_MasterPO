# Car Monitoring & Diagnostics - React Frontend + Go Backend

A sophisticated car monitoring application with a React frontend and Go backend, featuring vehicle telemetry, diagnostics, maintenance tracking, and 2GIS service station integration.

## 🏗️ Architecture

- **Frontend**: React + TypeScript + Vite (in `front/` folder)
- **Backend**: Go REST API (in `cmd/server/`)
- **Integration**: 2GIS API for service stations

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+ and npm
- (Optional) 2GIS API key for service station data

### Development Setup

1. **Start the Backend**:
   ```bash
   # Install Go dependencies
   go mod tidy
   
   # Run backend server
   go run cmd/server/main.go
   ```
   Backend runs on `http://localhost:8080`

2. **Start the Frontend** (in a separate terminal):
   ```bash
   cd front
   npm install
   npm run dev
   ```
   Frontend runs on `http://localhost:5173` (Vite dev server)

   The Vite config is set up to proxy `/api` requests to the backend.

### Production Build

1. **Build the Frontend**:
   ```bash
   cd front
   npm run build
   ```
   This creates `front/dist/` with the production build.

2. **Run the Backend** (serves the built frontend):
   ```bash
   go run cmd/server/main.go
   ```
   The backend will automatically serve the React app from `front/dist/`.

## 📡 API Endpoints

All API endpoints are prefixed with `/api`:

### Dashboard
- `GET /api/dashboard/{carId}` - Get dashboard data (car info, health metrics, alerts)

### Service History
- `GET /api/service-history/{carId}` - Get maintenance history

### Notifications
- `GET /api/notifications/{carId}` - Get alerts and notifications

### Repair Shops (2GIS)
- `GET /api/repair-shops?lat=43.2220&lon=76.8512&radius=10` - Get nearby service stations

### Legacy Endpoints
- `GET /api/cars` - List all cars
- `GET /api/cars/{id}` - Get car details
- `GET /api/cars/{id}/telemetry` - Get telemetry data
- `GET /api/cars/{id}/analysis` - Get analysis and maintenance recommendations

## 🔧 Configuration

### Environment Variables

- `PORT` - Backend server port (default: 8080)
- `TWOGIS_API_KEY` - 2GIS API key for service station data (optional, uses mock data if not set)

```bash
export TWOGIS_API_KEY=your_api_key_here
go run cmd/server/main.go
```

## 📁 Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Backend entry point
├── internal/
│   ├── handlers/                # HTTP handlers
│   ├── services/                # Business logic
│   ├── models/                  # Data models
│   └── data/                    # Mock data
├── front/                       # React frontend
│   ├── src/
│   │   └── app/
│   │       ├── components/      # React components
│   │       └── App.tsx
│   ├── dist/                    # Production build (generated)
│   └── vite.config.ts
└── go.mod
```

## 🎨 Frontend Features

- **Dashboard**: Vehicle health metrics, oil change countdown, alerts
- **Damage Assessment**: AI-powered damage detection (UI ready)
- **Service History**: Complete maintenance records timeline
- **Service Map**: Nearby repair shops with 2GIS integration

## 🔄 Development Workflow

1. **Backend changes**: Edit Go files, restart backend server
2. **Frontend changes**: Edit React files, Vite hot-reloads automatically
3. **API integration**: Frontend calls `/api/*` endpoints (proxied by Vite in dev)

## 📝 Notes

- The old `web/` folder is deprecated - use the React frontend in `front/`
- In development, frontend and backend run separately (frontend on 5173, backend on 8080)
- In production, backend serves the built frontend from `front/dist/`
- All API endpoints support CORS for frontend integration
- 2GIS integration falls back to mock data if API key is not provided

## 🚨 Troubleshooting

**Frontend not loading?**
- Make sure you've built it: `cd front && npm run build`
- Or run dev server separately: `cd front && npm run dev`

**API calls failing?**
- Check backend is running on port 8080
- In dev mode, Vite proxies `/api` to backend automatically
- Check CORS headers are set correctly

**2GIS not working?**
- Set `TWOGIS_API_KEY` environment variable
- System falls back to mock data if key is missing
