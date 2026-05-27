package job

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID        uuid.UUID
	Type      string
	Payload   []byte
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}