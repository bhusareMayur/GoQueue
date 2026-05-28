package job

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, job *Job) error
	GetByID(ctx context.Context, id uuid.UUID) (*Job, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*Job, error) // NEW: Idempotency Check
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	
	UpdateRetry(ctx context.Context, id uuid.UUID, retryCount int, lastError string, nextRunAt *time.Time, status string) error
	MoveToDLQ(ctx context.Context, deadJob *DeadJob) error

	ClaimJob(ctx context.Context, id uuid.UUID, workerID string) error
	GetStuckJobs(ctx context.Context, cutoffTime time.Time) ([]*Job, error)
	RequeueStuckJob(ctx context.Context, id uuid.UUID) error
}