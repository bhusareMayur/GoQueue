package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/bhusareMayur/goqueue/internal/domain/job"
)

type CreateJobRequest struct {
	Type     string          `json:"type"`
	Payload  json.RawMessage `json:"payload"`
	Priority string          `json:"priority"`
}

type UpdateJobStatusRequest struct {
	Status string `json:"status"`
}

type JobHandler struct {
	service *job.Service
}

func NewJobHandler(
	service *job.Service,
) *JobHandler {
	return &JobHandler{
		service: service,
	}
}

func (h *JobHandler) CreateJob(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req CreateJobRequest

	if err := json.NewDecoder(r.Body).
		Decode(&req); err != nil {

		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)

		return
	}

	// NEW: Extract Idempotency Key from Request Headers
	idempotencyKey := r.Header.Get("Idempotency-Key")

	j, err := h.service.CreateJob(
		r.Context(),
		req.Type,
		req.Payload,
		req.Priority,
		idempotencyKey,
	)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().
		Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(j)
}

func (h *JobHandler) GetJob(
	w http.ResponseWriter,
	r *http.Request,
) {

	idParam := r.URL.Query().Get("id")

	if idParam == "" {

		http.Error(
			w,
			"missing job id",
			http.StatusBadRequest,
		)

		return
	}

	jobID, err := uuid.Parse(idParam)

	if err != nil {

		http.Error(
			w,
			"invalid uuid",
			http.StatusBadRequest,
		)

		return
	}

	j, err := h.service.GetJob(
		r.Context(),
		jobID,
	)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().
		Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(j)
}

func (h *JobHandler) UpdateJobStatus(
	w http.ResponseWriter,
	r *http.Request,
) {

	idParam := r.URL.Query().Get("id")

	if idParam == "" {

		http.Error(
			w,
			"missing job id",
			http.StatusBadRequest,
		)

		return
	}

	jobID, err := uuid.Parse(idParam)

	if err != nil {

		http.Error(
			w,
			"invalid uuid",
			http.StatusBadRequest,
		)

		return
	}

	var req UpdateJobStatusRequest

	if err := json.NewDecoder(r.Body).
		Decode(&req); err != nil {

		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)

		return
	}

	err = h.service.UpdateJobStatus(
		r.Context(),
		jobID,
		req.Status,
	)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "job status updated",
		},
	)
}