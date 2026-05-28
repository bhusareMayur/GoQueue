package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/bhusareMayur/goqueue/internal/api/handlers"
)

func NewRouter(
	jobHandler *handlers.JobHandler,
) http.Handler {

	r := chi.NewRouter()

	// NEW: Expose metrics endpoint
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	r.Post(
		"/jobs",
		jobHandler.CreateJob,
	)

	return r
}