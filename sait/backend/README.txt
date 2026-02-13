Backend (Go) — bazarPO

1) Create .env from .env.example
2) Install deps:
   go mod tidy
3) Seed DB:
   go run ./cmd/seed
4) Run server:
   go run ./cmd/server

Health:
  http://localhost:8090/api/health
