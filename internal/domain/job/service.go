package job

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/bhusareMayur/goqueue/internal/observability/metrics"
)

type Service struct {
	repo  Repository
	queue Queue
}

func NewService(repo Repository, queue Queue) *Service {
	return &Service{
		repo:  repo,
		queue: queue,
	}
}

// NEW: Backpressure helper
func (s *Service) GetQueueLength(ctx context.Context) (int64, error) {
	return s.queue.GetQueueLength(ctx)
}

func (s *Service) CreateJob(ctx context.Context, jobType string, payload []byte, priority string, idempotencyKey string, correlationID string) (*Job, error) {
	if idempotencyKey != "" {
		existingJob, err := s.repo.GetByIdempotencyKey(ctx, idempotencyKey)
		if err == nil && existingJob != nil {
			return existingJob, nil
		}
	}

	priority = strings.ToLower(priority)
	if priority != "high" && priority != "medium" && priority != "low" {
		priority = "default"
	}

	var idKey *string
	if idempotencyKey != "" {
		idKey = &idempotencyKey
	}
	
	var corrID *string
	if correlationID != "" {
		corrID = &correlationID
	}

	j := &Job{
		ID:             uuid.New(),
		Type:           jobType,
		Payload:        payload,
		Status:         "pending",
		Priority:       priority,
		RetryCount:     0,
		MaxRetries:     5,
		IdempotencyKey: idKey,
		CorrelationID:  corrID,
	}

	event := &OutboxEvent{
		ID:        uuid.New(),
		JobID:     j.ID,
		Priority:  j.Priority,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	// Atomically save both Job and OutboxEvent to Postgres
	if err := s.repo.CreateWithOutbox(ctx, j, event); err != nil {
		return nil, err
	}

	// We NO LONGER enqueue to Redis here. The background Publisher will handle it.
	return j, nil
}

func (s *Service) GetJob(ctx context.Context, id uuid.UUID) (*Job, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateJobStatus(ctx context.Context, id uuid.UUID, status string) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *Service) UpdateJobRetry(ctx context.Context, id uuid.UUID, retryCount int, lastError string, nextRunAt *time.Time, status string) error {
	return s.repo.UpdateRetry(ctx, id, retryCount, lastError, nextRunAt, status)
}

func (s *Service) MoveToDLQ(ctx context.Context, id uuid.UUID, errMessage string) error {
	j, err := s.GetJob(ctx, id)
	if err != nil {
		return err
	}

	deadJob := &DeadJob{
		ID:            id,
		Type:          j.Type,
		Payload:       j.Payload,
		Priority:      j.Priority,
		RetryCount:    j.RetryCount,
		LastError:     errMessage,
		FailedAt:      time.Now(),
		CorrelationID: j.CorrelationID,
	}

	err = s.repo.MoveToDLQ(ctx, deadJob)
	
	if err == nil {
		metrics.DeadLetterJobs.WithLabelValues(j.Priority).Inc()
	}

	return err
}

func (s *Service) ClaimJob(ctx context.Context, id uuid.UUID, workerID string) error {
	return s.repo.ClaimJob(ctx, id, workerID)
}

func (s *Service) GetStuckJobs(ctx context.Context, timeout time.Duration) ([]*Job, error) {
	cutoffTime := time.Now().Add(-timeout)
	return s.repo.GetStuckJobs(ctx, cutoffTime)
}

func (s *Service) RequeueStuckJob(ctx context.Context, id uuid.UUID) error {
	return s.repo.RequeueStuckJob(ctx, id)
}
// Add these to the bottom of the file
func (s *Service) GetPendingOutboxEvents(ctx context.Context, limit int) ([]*OutboxEvent, error) {
	return s.repo.GetPendingOutboxEvents(ctx, limit)
}

func (s *Service) MarkOutboxEventPublished(ctx context.Context, id uuid.UUID) error {
	return s.repo.MarkOutboxEventPublished(ctx, id)
}