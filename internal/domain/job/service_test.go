package job

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =====================================================================
// MOCK IMPLEMENTATIONS
// =====================================================================

type mockQueue struct {
	enqueueCalls int
}

func (m *mockQueue) Enqueue(ctx context.Context, jobID string, priority string) error {
	m.enqueueCalls++
	return nil
}
func (m *mockQueue) Consume(ctx context.Context) (string, error) { return "", nil }
func (m *mockQueue) EnqueueDelayed(ctx context.Context, jobID string, priority string, runAt time.Time) error {
	return nil
}

type mockRepo struct {
	existingJob *Job
	createCalls int
}

func (m *mockRepo) GetByIdempotencyKey(ctx context.Context, key string) (*Job, error) {
	if m.existingJob != nil && m.existingJob.IdempotencyKey != nil && *m.existingJob.IdempotencyKey == key {
		return m.existingJob, nil
	}
	// Simulate database returning "not found"
	return nil, errors.New("job not found")
}

func (m *mockRepo) Create(ctx context.Context, j *Job) error {
	m.createCalls++
	return nil
}

// Stub remaining required Repository interface methods
func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID) (*Job, error) { return nil, nil }
func (m *mockRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error { return nil }
func (m *mockRepo) UpdateRetry(ctx context.Context, id uuid.UUID, retryCount int, lastError string, nextRunAt *time.Time, status string) error { return nil }
func (m *mockRepo) MoveToDLQ(ctx context.Context, deadJob *DeadJob) error { return nil }
func (m *mockRepo) ClaimJob(ctx context.Context, id uuid.UUID, workerID string) error { return nil }
func (m *mockRepo) GetStuckJobs(ctx context.Context, cutoffTime time.Time) ([]*Job, error) { return nil, nil }
func (m *mockRepo) RequeueStuckJob(ctx context.Context, id uuid.UUID) error { return nil }

// =====================================================================
// TESTS
// =====================================================================

func TestService_CreateJob_Idempotency(t *testing.T) {
	t.Run("first request creates a new job", func(t *testing.T) {
		repo := &mockRepo{}
		queue := &mockQueue{}
		svc := NewService(repo, queue)

		idempotencyKey := "test-key-123"

		// Execute - Added "" as the 6th argument for correlationID
		j, err := svc.CreateJob(context.Background(), "email", []byte(`{}`), "high", idempotencyKey, "")

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, j)
		assert.Equal(t, "email", j.Type)
		assert.NotNil(t, j.IdempotencyKey)
		assert.Equal(t, idempotencyKey, *j.IdempotencyKey)
		
		// Ensure database and queue were called exactly once
		assert.Equal(t, 1, repo.createCalls, "Expected Repo.Create to be called once")
		assert.Equal(t, 1, queue.enqueueCalls, "Expected Queue.Enqueue to be called once")
	})

	t.Run("duplicate request returns existing job (idempotent)", func(t *testing.T) {
		idempotencyKey := "test-key-123"
		
		// Setup mock repo to pretend it already has this job in the database
		existing := &Job{
			ID:             uuid.New(),
			Type:           "email",
			Status:         "completed",
			IdempotencyKey: &idempotencyKey,
		}
		
		repo := &mockRepo{existingJob: existing}
		queue := &mockQueue{}
		svc := NewService(repo, queue)

		// Execute - Added "" as the 6th argument for correlationID
		j, err := svc.CreateJob(context.Background(), "email", []byte(`{}`), "high", idempotencyKey, "")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, existing.ID, j.ID, "Should return the exact same job ID as the existing job")
		assert.Equal(t, "completed", j.Status, "Should return the existing job's updated status")

		// Ensure we DID NOT create a new record in the DB or push to Redis
		assert.Equal(t, 0, repo.createCalls, "Expected Repo.Create to NOT be called")
		assert.Equal(t, 0, queue.enqueueCalls, "Expected Queue.Enqueue to NOT be called")
	})

	t.Run("standard request without key creates normally", func(t *testing.T) {
		repo := &mockRepo{}
		queue := &mockQueue{}
		svc := NewService(repo, queue)

		// Execute with an empty string for the idempotency key AND correlationID
		j, err := svc.CreateJob(context.Background(), "image", []byte(`{}`), "low", "", "")

		// Assert
		require.NoError(t, err)
		assert.Nil(t, j.IdempotencyKey, "IdempotencyKey should be nil for standard jobs")
		
		assert.Equal(t, 1, repo.createCalls, "Expected Repo.Create to be called once")
		assert.Equal(t, 1, queue.enqueueCalls, "Expected Queue.Enqueue to be called once")
	})
}