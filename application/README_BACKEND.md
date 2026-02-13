# Car Monitoring Backend - Updated for React Frontend

This backend has been updated to work with the React frontend in the `front` folder. It provides API endpoints for dashboard data, service history, repair shops (with 2GIS integration), and notifications.

## 🚀 Features

- **Dashboard API**: Returns car info, health metrics, oil change data, and alerts
- **Service History API**: Returns maintenance history records
- **Repair Shops API**: Integrated with 2GIS API for real service station data (with fallback to mock data)
- **Notifications API**: Returns alerts and maintenance reminders
- **CORS Support**: Enabled for frontend integration
- **2GIS Integration**: Real-time service station data from 2GIS API

## 📡 API Endpoints

### Dashboard
```
GET /api/dashboard/{carId}
```
Returns dashboard data including:
- Car information (model, year, plate, mileage)
- Health metrics (Engine, Battery, Oil Level)
- Oil change countdown
- Recent alerts

**Response Example:**
```json
{
  "carInfo": {
    "model": "Toyota Camry",
    "year": "2020",
    "plate": "ABC-1234",
    "mileage": 45000
  },
  "healthMetrics": [
    {
      "label": "Engine",
      "value": 95,
      "status": "good",
      "icon": "Gauge",
      "color": "text-emerald-400"
    }
  ],
  "oilChangeData": {
    "currentKm": 45000,
    "nextChangeKm": 50000,
    "daysRemaining": 28
  },
  "recentAlerts": [...]
}
```

### Service History
```
GET /api/service-history/{carId}
```
Returns array of service history records.

**Response Example:**
```json
[
  {
    "id": 1,
    "date": "2025-11-15",
    "type": "Oil Change",
    "description": "Full synthetic oil change and filter replacement",
    "mileage": 45234,
    "cost": 15000,
    "shop": "AutoService Premium",
    "location": "Almaty, Abay Ave 150",
    "verified": true,
    "icon": "Droplet",
    "color": "cyan"
  }
]
```

### Repair Shops (2GIS)
```
GET /api/repair-shops?lat=43.2220&lon=76.8512&radius=10&rating=4&price=medium
```
Returns nearby repair shops from 2GIS API.

**Query Parameters:**
- `lat` - Latitude (default: 43.2220 - Almaty)
- `lon` - Longitude (default: 76.8512 - Almaty)
- `radius` - Search radius in km (default: 10)
- `rating` - Minimum rating filter
- `price` - Price level filter (cheap, medium, expensive)

**Response Example:**
```json
[
  {
    "id": 1,
    "name": "AutoService Premium",
    "rating": 4.8,
    "reviews": 234,
    "distance": 1.2,
    "address": "Abay Ave 150, Almaty",
    "phone": "+7 727 123 4567",
    "hours": "Open until 8:00 PM",
    "services": ["Oil Change", "Diagnostics"],
    "verified": true,
    "priceLevel": 2
  }
]
```

### Notifications
```
GET /api/notifications/{carId}
```
Returns notifications and alerts for the car.

**Response Example:**
```json
[
  {
    "id": "alert-high-temp",
    "type": "critical",
    "title": "High Engine Temperature",
    "message": "Engine coolant temperature is 108.5°C",
    "timestamp": 1703520000,
    "read": false,
    "actionable": true
  }
]
```

## 🔧 Configuration

### 2GIS API Integration

To enable 2GIS API integration, set the API key as an environment variable:

```bash
export TWOGIS_API_KEY=your_api_key_here
```

If the API key is not set, the system will automatically fall back to mock data.

**Getting a 2GIS API Key:**
1. Register at https://dev.2gis.com/
2. Create a new application
3. Get your API key from the dashboard
4. Set it as `TWOGIS_API_KEY` environment variable

### Environment Variables

- `PORT` - Server port (default: 8080)
- `TWOGIS_API_KEY` - 2GIS API key for service station data

## 🏃 Running the Server

```bash
# Install dependencies
go mod tidy

# Run server
go run cmd/server/main.go

# Or with environment variables
TWOGIS_API_KEY=your_key PORT=8080 go run cmd/server/main.go
```

## 📁 Frontend Integration

The backend serves the React frontend from:
- Production: `./front/dist` (after `npm run build`)
- Fallback: `./web` (static HTML fallback)

### Development Setup

1. **Backend** (this server):
   ```bash
   go run cmd/server/main.go
   ```

2. **Frontend** (separate terminal):
   ```bash
   cd front
   npm install
   npm run dev
   ```

The frontend dev server runs on port 5173 (Vite default). You can configure the frontend to proxy API requests to the backend on port 8080.

### Production Build

1. Build the frontend:
   ```bash
   cd front
   npm run build
   ```

2. The backend will automatically serve from `./front/dist`

## 🔄 API Compatibility

The backend maintains backward compatibility with the original endpoints:
- `GET /api/cars` - List all cars
- `GET /api/cars/{id}` - Get car details
- `GET /api/cars/{id}/telemetry` - Get telemetry
- `GET /api/cars/{id}/analysis` - Get analysis
- `GET /api/sto` - Legacy service stations (mock data)

## 🚨 Notifications System

The notification system automatically generates alerts based on:
- Critical telemetry thresholds (high temp, low battery)
- Maintenance due dates (oil change, spark plugs, etc.)
- Diagnostic trouble codes (DTC)
- System errors (ABS, airbag, emissions)

Notifications are returned in real-time and include actionable items that users can address.

## 📝 Notes

- The backend uses mock car data by default (5 vehicles)
- 2GIS integration requires an API key; falls back to mock data if unavailable
- All endpoints support CORS for frontend integration
- Service history is generated based on car mileage and maintenance intervals

