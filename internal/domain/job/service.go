package job

import (
	"context"
	"time"

	"github.com/google/uuid"
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

func (s *Service) CreateJob(ctx context.Context, jobType string, payload []byte) (*Job, error) {
	j := &Job{
		ID:         uuid.New(),
		Type:       jobType,
		Payload:    payload,
		Status:     "pending",
		RetryCount: 0,
		MaxRetries: 5, // Set default max retries
	}

	if err := s.repo.Create(ctx, j); err != nil {
		return nil, err
	}

	if err := s.queue.Enqueue(ctx, j.ID.String()); err != nil {
		return nil, err
	}

	return j, nil
}

func (s *Service) GetJob(ctx context.Context, id uuid.UUID) (*Job, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateJobStatus(ctx context.Context, id uuid.UUID, status string) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

// UpdateJobRetry saves the failure state and calculates the next run backoff
func (s *Service) UpdateJobRetry(ctx context.Context, id uuid.UUID, retryCount int, lastError string, nextRunAt *time.Time, status string) error {
	return s.repo.UpdateRetry(ctx, id, retryCount, lastError, nextRunAt, status)
}