package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bhusareMayur/goqueue/internal/domain/job"
)

type JobRepository struct {
	db *pgxpool.Pool
}

func NewJobRepository(db *pgxpool.Pool) *JobRepository {
	return &JobRepository{
		db: db,
	}
}

func (r *JobRepository) Create(ctx context.Context, j *job.Job) error {
	query := `
	INSERT INTO jobs (
		id, type, payload, status, retry_count, max_retries, next_run_at, last_error
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(
		ctx, query,
		j.ID, j.Type, j.Payload, j.Status,
		j.RetryCount, j.MaxRetries, j.NextRunAt, j.LastError,
	)

	return err
}

func (r *JobRepository) GetByID(ctx context.Context, id uuid.UUID) (*job.Job, error) {
	query := `
	SELECT
		id, type, payload, status, retry_count, max_retries, next_run_at, last_error, created_at, updated_at
	FROM jobs
	WHERE id = $1
	`

	var j job.Job
	var lastError *string // pointer allows us to safely scan NULL values from postgres

	err := r.db.QueryRow(ctx, query, id).Scan(
		&j.ID, &j.Type, &j.Payload, &j.Status,
		&j.RetryCount, &j.MaxRetries, &j.NextRunAt, &lastError,
		&j.CreatedAt, &j.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if lastError != nil {
		j.LastError = *lastError
	}

	return &j, nil
}

func (r *JobRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `
	UPDATE jobs
	SET status = $1, updated_at = NOW()
	WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

func (r *JobRepository) UpdateRetry(
	ctx context.Context,
	id uuid.UUID,
	retryCount int,
	lastError string,
	nextRunAt *time.Time,
	status string,
) error {
	query := `
	UPDATE jobs
	SET retry_count = $1, last_error = $2, next_run_at = $3, status = $4, updated_at = NOW()
	WHERE id = $5
	`

	_, err := r.db.Exec(ctx, query, retryCount, lastError, nextRunAt, status, id)
	return err
}

// NEW: MoveToDLQ uses a transaction to ensure both tables are updated safely
func (r *JobRepository) MoveToDLQ(ctx context.Context, dj *job.DeadJob) error {
	// Start a database transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails. If tx.Commit() is called first, rollback does nothing.
	defer tx.Rollback(ctx)

	// 1. Insert into dead_jobs table
	insertDLQQuery := `
	INSERT INTO dead_jobs (id, type, payload, retry_count, failed_at, last_error)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(ctx, insertDLQQuery, dj.ID, dj.Type, dj.Payload, dj.RetryCount, dj.FailedAt, dj.LastError)
	if err != nil {
		return err
	}

	// 2. Update main jobs table status to 'failed'
	updateJobQuery := `
	UPDATE jobs
	SET status = 'failed', last_error = $1, updated_at = NOW()
	WHERE id = $2
	`
	_, err = tx.Exec(ctx, updateJobQuery, dj.LastError, dj.ID)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit(ctx)
}