package handlers

import (
	"encoding/json"
	"log/slog"
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")
	
	// NEW: Extract or Generate Correlation ID
	correlationID := r.Header.Get("X-Correlation-ID")
	if correlationID == "" {
		correlationID = uuid.New().String()
	}

	reqLogger := slog.With("correlation_id", correlationID)
	reqLogger.Info("received create job request", "type", req.Type, "priority", req.Priority)

	j, err := h.service.CreateJob(
		r.Context(),
		req.Type,
		req.Payload,
		req.Priority,
		idempotencyKey,
		correlationID,
	)

	if err != nil {
		reqLogger.Error("failed to create job", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reqLogger.Info("job created successfully", "job_id", j.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(j)
}

func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		http.Error(w, "missing job id", http.StatusBadRequest)
		return
	}

	jobID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	j, err := h.service.GetJob(r.Context(), jobID)
	if err != nil {
		slog.Error("failed to get job", "job_id", jobID, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(j)
}

func (h *JobHandler) UpdateJobStatus(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		http.Error(w, "missing job id", http.StatusBadRequest)
		return
	}

	jobID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	var req UpdateJobStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateJobStatus(r.Context(), jobID, req.Status)
	if err != nil {
		slog.Error("failed to update job status", "job_id", jobID, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "job status updated",
	})
}