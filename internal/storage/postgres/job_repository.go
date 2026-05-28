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
		id, type, payload, status, priority, retry_count, max_retries, next_run_at, last_error, idempotency_key
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(
		ctx, query,
		j.ID, j.Type, j.Payload, j.Status, j.Priority,
		j.RetryCount, j.MaxRetries, j.NextRunAt, j.LastError, j.IdempotencyKey,
	)
	return err
}

func (r *JobRepository) GetByID(ctx context.Context, id uuid.UUID) (*job.Job, error) {
	query := `
	SELECT
		id, type, payload, status, priority, retry_count, max_retries, next_run_at, last_error, created_at, updated_at, worker_id, processing_started_at, idempotency_key
	FROM jobs
	WHERE id = $1
	`

	var j job.Job
	var lastError *string 

	err := r.db.QueryRow(ctx, query, id).Scan(
		&j.ID, &j.Type, &j.Payload, &j.Status, &j.Priority, 
		&j.RetryCount, &j.MaxRetries, &j.NextRunAt, &lastError,
		&j.CreatedAt, &j.UpdatedAt, &j.WorkerID, &j.ProcessingStartedAt, &j.IdempotencyKey,
	)

	if err != nil {
		return nil, err
	}

	if lastError != nil {
		j.LastError = *lastError
	}

	return &j, nil
}

func (r *JobRepository) GetByIdempotencyKey(ctx context.Context, key string) (*job.Job, error) {
	query := `
	SELECT
		id, type, payload, status, priority, retry_count, max_retries, next_run_at, last_error, created_at, updated_at, worker_id, processing_started_at, idempotency_key
	FROM jobs
	WHERE idempotency_key = $1
	`

	var j job.Job
	var lastError *string 

	err := r.db.QueryRow(ctx, query, key).Scan(
		&j.ID, &j.Type, &j.Payload, &j.Status, &j.Priority,
		&j.RetryCount, &j.MaxRetries, &j.NextRunAt, &lastError,
		&j.CreatedAt, &j.UpdatedAt, &j.WorkerID, &j.ProcessingStartedAt, &j.IdempotencyKey,
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
	SET status = $1, updated_at = $2
	WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, status, time.Now(), id)
	return err
}

func (r *JobRepository) UpdateRetry(
	ctx context.Context, id uuid.UUID, retryCount int, lastError string, nextRunAt *time.Time, status string,
) error {
	query := `
	UPDATE jobs
	SET retry_count = $1, last_error = $2, next_run_at = $3, status = $4, updated_at = $5
	WHERE id = $6
	`
	_, err := r.db.Exec(ctx, query, retryCount, lastError, nextRunAt, status, time.Now(), id)
	return err
}

func (r *JobRepository) MoveToDLQ(ctx context.Context, dj *job.DeadJob) error {
	tx, err := r.db.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	insertDLQQuery := `
	INSERT INTO dead_jobs (id, type, payload, priority, retry_count, failed_at, last_error)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(ctx, insertDLQQuery, dj.ID, dj.Type, dj.Payload, dj.Priority, dj.RetryCount, dj.FailedAt, dj.LastError)
	if err != nil { return err }

	updateJobQuery := `
	UPDATE jobs
	SET status = 'failed', last_error = $1, updated_at = $2
	WHERE id = $3
	`
	_, err = tx.Exec(ctx, updateJobQuery, dj.LastError, time.Now(), dj.ID)
	if err != nil { return err }

	return tx.Commit(ctx)
}

func (r *JobRepository) ClaimJob(ctx context.Context, id uuid.UUID, workerID string) error {
	now := time.Now()
	query := `
	UPDATE jobs
	SET status = 'processing', worker_id = $1, processing_started_at = $2, updated_at = $2
	WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, workerID, now, id)
	return err
}

func (r *JobRepository) GetStuckJobs(ctx context.Context, cutoffTime time.Time) ([]*job.Job, error) {
	query := `
	SELECT id, type, payload, status, priority, retry_count, max_retries, next_run_at, last_error, created_at, updated_at, worker_id, processing_started_at, idempotency_key
	FROM jobs
	WHERE status = 'processing' AND processing_started_at < $1
	`
	rows, err := r.db.Query(ctx, query, cutoffTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stuckJobs []*job.Job
	for rows.Next() {
		var j job.Job
		var lastError *string
		err := rows.Scan(
			&j.ID, &j.Type, &j.Payload, &j.Status, &j.Priority,
			&j.RetryCount, &j.MaxRetries, &j.NextRunAt, &lastError,
			&j.CreatedAt, &j.UpdatedAt, &j.WorkerID, &j.ProcessingStartedAt, &j.IdempotencyKey,
		)
		if err != nil {
			return nil, err
		}
		if lastError != nil {
			j.LastError = *lastError
		}
		stuckJobs = append(stuckJobs, &j)
	}

	return stuckJobs, nil
}

func (r *JobRepository) RequeueStuckJob(ctx context.Context, id uuid.UUID) error {
	query := `
	UPDATE jobs
	SET status = 'pending', retry_count = retry_count + 1, worker_id = NULL, processing_started_at = NULL, updated_at = $1
	WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}