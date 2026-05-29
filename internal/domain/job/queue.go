package job

import (
	"context"
	"time"
)

type Queue interface {
	// NEW: Added GetQueueLength for backpressure
	GetQueueLength(ctx context.Context) (int64, error)

	Enqueue(ctx context.Context, jobID string, priority string) error
	
	Consume(ctx context.Context) (string, error)

	EnqueueDelayed(ctx context.Context, jobID string, priority string, runAt time.Time) error
}