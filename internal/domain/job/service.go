package job

import (
	"context"
	"strings"
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

func (s *Service) CreateJob(ctx context.Context, jobType string, payload []byte, priority string) (*Job, error) {
	priority = strings.ToLower(priority)
	if priority != "high" && priority != "medium" && priority != "low" {
		priority = "default"
	}

	j := &Job{
		ID:         uuid.New(),
		Type:       jobType,
		Payload:    payload,
		Status:     "pending",
		Priority:   priority,
		RetryCount: 0,
		MaxRetries: 5,
	}

	if err := s.repo.Create(ctx, j); err != nil {
		return nil, err
	}

	if err := s.queue.Enqueue(ctx, j.ID.String(), j.Priority); err != nil {
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

func (s *Service) UpdateJobRetry(ctx context.Context, id uuid.UUID, retryCount int, lastError string, nextRunAt *time.Time, status string) error {
	return s.repo.UpdateRetry(ctx, id, retryCount, lastError, nextRunAt, status)
}

func (s *Service) MoveToDLQ(ctx context.Context, id uuid.UUID, errMessage string) error {
	j, err := s.GetJob(ctx, id)
	if err != nil {
		return err
	}

	deadJob := &DeadJob{
		ID:         id,
		Type:       j.Type,
		Payload:    j.Payload,
		Priority:   j.Priority,
		RetryCount: j.RetryCount,
		LastError:  errMessage,
		FailedAt:   time.Now(),
	}

	return s.repo.MoveToDLQ(ctx, deadJob)
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