package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bhusareMayur/goqueue/internal/domain/job"
)

type JobRepository struct {
	db *pgxpool.Pool
}

func NewJobRepository(
	db *pgxpool.Pool,
) *JobRepository {
	return &JobRepository{
		db: db,
	}
}

func (r *JobRepository) Create(
	ctx context.Context,
	j *job.Job,
) error {

	query := `
	INSERT INTO jobs (
		id,
		type,
		payload,
		status
	)
	VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		j.ID,
		j.Type,
		j.Payload,
		j.Status,
	)

	return err
}

func (r *JobRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*job.Job, error) {

	query := `
	SELECT
		id,
		type,
		payload,
		status,
		created_at,
		updated_at
	FROM jobs
	WHERE id = $1
	`

	var j job.Job

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&j.ID,
		&j.Type,
		&j.Payload,
		&j.Status,
		&j.CreatedAt,
		&j.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &j, nil
}

func (r *JobRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
) error {

	query := `
	UPDATE jobs
	SET
		status = $1,
		updated_at = NOW()
	WHERE id = $2
	`

	_, err := r.db.Exec(
		ctx,
		query,
		status,
		id,
	)

	return err
}