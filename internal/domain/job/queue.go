package job

import (
	"context"
	"time"
)

type Queue interface {
	Enqueue(ctx context.Context, jobID string) error
	
	Consume(ctx context.Context) (string, error)

	// NEW: Add ability to enqueue a job for the future
	EnqueueDelayed(ctx context.Context, jobID string, runAt time.Time) error
}