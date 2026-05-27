package job

import (
	"time"

	"github.com/google/uuid"
)

type DeadJob struct {
	ID         uuid.UUID
	Type       string
	Payload    []byte
	RetryCount int
	LastError  string
	FailedAt   time.Time
}