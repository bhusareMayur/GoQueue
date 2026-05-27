package job

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, job *Job) error

	GetByID(ctx context.Context, id uuid.UUID) (*Job, error)

	UpdateStatus(
		ctx context.Context,
		id uuid.UUID,
		status string,
	) error
}