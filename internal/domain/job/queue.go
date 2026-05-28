package job

import (
	"context"
	"time"
)

type Queue interface {
	// NEW: Added priority parameter
	Enqueue(ctx context.Context, jobID string, priority string) error
	
	Consume(ctx context.Context) (string, error)

	// NEW: Added priority parameter
	EnqueueDelayed(ctx context.Context, jobID string, priority string, runAt time.Time) error
}