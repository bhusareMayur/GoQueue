package job

import "context"

type Queue interface {
	Enqueue(
		ctx context.Context,
		jobID string,
	) error

	Consume(
		ctx context.Context,
	) (string, error)
}