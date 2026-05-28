package job

import (
	"time"

	"github.com/google/uuid"
)

type DeadJob struct {
	ID         uuid.UUID
	Type       string
	Payload    []byte
	Priority   string     // NEW: Priority level
	RetryCount int
	LastError  string
	FailedAt   time.Time
}