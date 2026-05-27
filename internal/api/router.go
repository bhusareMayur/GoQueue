package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/bhusareMayur/goqueue/internal/api/handlers"
)

func NewRouter(
	jobHandler *handlers.JobHandler,
) http.Handler {

	r := chi.NewRouter()

	r.Post(
		"/jobs",
		jobHandler.CreateJob,
	)

	return r
}