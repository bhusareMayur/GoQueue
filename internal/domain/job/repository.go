package job

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, job *Job) error
	GetByID(ctx context.Context, id uuid.UUID) (*Job, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	
	// New method for handling failures
	UpdateRetry(ctx context.Context, id uuid.UUID, retryCount int, lastError string, nextRunAt *time.Time, status string) error
	
	// NEW: Method to handle DLQ movement
	MoveToDLQ(ctx context.Context, deadJob *DeadJob) error
}