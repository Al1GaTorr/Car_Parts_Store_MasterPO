# CPS Workspace Run Guide

This repository contains two separate apps:

- `sait` - bazarPO marketplace app (React + Go + MongoDB)
- `application` - car monitoring app (React + Go, Mongo optional)

Use the `sait` project if you want auth, parts catalog, orders, and admin flows.

## Prerequisites

- Go `1.22+` (recommended for both apps)
- Node.js `18+` and npm
- MongoDB `6+` running locally for `sait`

## Project Layout

- `sait/` - frontend
- `sait/backend/` - backend API, auth, seed
- `application/front/` - frontend
- `application/` - backend API

## 1) Run SAIT (Recommended)

### Step 1. Start MongoDB

Run MongoDB locally so `sait/backend` can connect to:

`mongodb://127.0.0.1:27017/bazarPO`

### Step 2. Backend setup

```bash
cd sait/backend
go mod tidy
go run ./cmd/seed
go run ./cmd/server
```

Expected backend URL:

- `http://localhost:8090`
- Health check: `http://localhost:8090/api/health`

### Step 3. Frontend setup

Open a second terminal:

```bash
cd sait
npm install
npm run dev
```

Expected frontend URL:

- `http://localhost:5174`

The frontend proxies `/api/*` requests to `http://localhost:8090`.

### Step 4. Login defaults

Default admin credentials created by backend:

- Email: `admin@cps.local`
- Password: `admin12345`

## SAIT Environment Variables

`sait/backend/internal/model/env.go` uses these variables:

- `MONGO_URI` (default `mongodb://127.0.0.1:27017/bazarPO`)
- `MONGO_DB` (default `bazarPO`)
- `JWT_SECRET` (default `change-me`)
- `PORT` (default `8090`)
- `ADMIN_EMAIL` (default `admin@cps.local`)
- `ADMIN_PASSWORD` (default `admin12345`)

Important: Go does not auto-load `.env` by itself. If you want to load `sait/backend/.env` in shell:

```bash
cd sait/backend
set -a
source .env
set +a
go run ./cmd/server
```

### SAIT API Collection

Postman collection file:

- `sait/backend/postman_collection.json`

Import this file into Postman and set:

- `baseUrl = http://localhost:8090`

## 2) Run APPLICATION (Car Monitoring)

### Step 1. Backend

```bash
cd application
go run cmd/server/main.go
```

Expected backend URL:

- `http://localhost:8081`
- Health check: `http://localhost:8081/api/health`

Notes:

- Backend tries Mongo first (`MONGO_URI`), then falls back to mock data.
- Mongo is optional for this app.

### Step 2. Frontend

Open a second terminal:

```bash
cd application/front
npm install
npm run dev
```

Expected frontend URL:

- `http://localhost:5173`

The frontend proxies `/api/*` to `http://localhost:8081`.

### Optional production-style run (APPLICATION)

```bash
cd application/front
npm run build
cd ..
go run cmd/server/main.go
```

Backend serves static frontend from `application/front/dist`.

## APPLICATION Environment Variables

Common variables:

- `PORT` (default `8081` in code)
- `MONGO_URI` (optional)
- `TWOGIS_API_KEY` (optional)

If using `application/.env`, load it in shell before running backend:

```bash
cd application
set -a
source .env
set +a
go run cmd/server/main.go
```

## Quick Verification Checklist

For `sait`:

- `GET http://localhost:8090/api/health` returns `{ "ok": true }`
- UI opens at `http://localhost:5174`
- Login with `admin@cps.local / admin12345`

For `application`:

- `GET http://localhost:8081/api/health` returns healthy response
- UI opens at `http://localhost:5173`

## Troubleshooting

- Port already in use:
  Run on another port via `PORT=<new_port>` and update frontend proxy target.
- Frontend cannot reach backend:
  Confirm backend is running and proxy target matches backend port.
- Mongo connection errors in `sait`:
  Ensure MongoDB is running and `MONGO_URI`/`MONGO_DB` are correct.
- Seed fails in `sait/backend`:
  Run from `sait/backend` directory so `seed_data/cars.json` resolves correctly.
