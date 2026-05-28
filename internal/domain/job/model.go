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
	Priority            string
	RetryCount          int
	MaxRetries          int
	NextRunAt           *time.Time
	LastError           string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	WorkerID            *string
	ProcessingStartedAt *time.Time
	IdempotencyKey      *string    // NEW: Idempotency Key
}