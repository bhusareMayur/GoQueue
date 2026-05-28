package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/bhusareMayur/goqueue/internal/api/handlers"
)

func NewRouter(
	jobHandler *handlers.JobHandler,
	healthHandler *handlers.HealthHandler, // NEW: Inject HealthHandler
) http.Handler {

	r := chi.NewRouter()

	// Expose metrics endpoint
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	// NEW: Expose Health Check endpoints
	r.Get("/live", healthHandler.Live)
	r.Get("/ready", healthHandler.Ready)

	// Job Endpoints
	r.Post(
		"/jobs",
		jobHandler.CreateJob,
	)

	return r
}