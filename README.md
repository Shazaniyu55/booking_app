# Booking API (Go + chi + MongoDB)

## Setup
    docker compose up -d
    go get go.mongodb.org/mongo-driver/v2@latest
    go get github.com/go-chi/chi/v5@latest
    go mod tidy
    go run ./cmd/api

Env vars (optional): PORT=8080, MONGO_URI=mongodb://localhost:27017, MONGO_DB=booking_app

## Endpoints
- GET  /health
- POST /api/v1/resources        GET /api/v1/resources        GET /api/v1/resources/{id}
- POST /api/v1/bookings         GET /api/v1/bookings?resource_id=   GET /api/v1/bookings/{id}
- POST /api/v1/bookings/{id}/cancel
