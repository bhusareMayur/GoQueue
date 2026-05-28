package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
}

// NewHealthHandler creates a new handler injected with DB and Redis dependencies
func NewHealthHandler(dbPool *pgxpool.Pool, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{
		dbPool:      dbPool,
		redisClient: redisClient,
	}
}

// Live handles the Liveness check (/live)
// It only checks if the HTTP server is responsive.
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "alive"}`))
}

// Ready handles the Readiness check (/ready)
// It pings both PostgreSQL and Redis to ensure the app can handle traffic.
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Set a short timeout for the health checks so they don't hang
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// 1. Check PostgreSQL Connection
	if err := h.dbPool.Ping(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status": "unavailable", "error": "database unreachable"}`))
		return
	}

	// 2. Check Redis Connection
	if err := h.redisClient.Ping(ctx).Err(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status": "unavailable", "error": "redis unreachable"}`))
		return
	}

	// Both dependencies are healthy
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ready"}`))
}