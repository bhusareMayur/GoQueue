package job

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID                  uuid.UUID
	Type                string
	Payload             []byte
	Status              string
	RetryCount          int
	MaxRetries          int
	NextRunAt           *time.Time // pointer because it can be nil
	LastError           string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	WorkerID            *string    // NEW: Tracks which worker took it
	ProcessingStartedAt *time.Time // NEW: Tracks when processing started
}