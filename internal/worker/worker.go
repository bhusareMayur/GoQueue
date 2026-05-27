package worker

import (
	"context"
	"errors"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/bhusareMayur/goqueue/internal/domain/job"
)

type Worker struct {
	id      int
	queue   job.Queue
	service *job.Service
}

func NewWorker(id int, queue job.Queue, service *job.Service) *Worker {
	return &Worker{
		id:      id,
		queue:   queue,
		service: service,
	}
}

func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("worker-%d started\n", w.id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker-%d shutting down\n", w.id)
			return
		default:
		}

		// 1. Consume job ID from Redis
		jobID, err := w.queue.Consume(ctx)
		if err != nil {
			log.Printf("worker-%d consume error: %v\n", w.id, err)
			continue
		}

		if jobID == "" {
			continue
		}

		parsedID, err := uuid.Parse(jobID)
		if err != nil {
			log.Printf("worker-%d invalid uuid error: %v\n", w.id, err)
			continue
		}

		// 2. Fetch job from PostgreSQL
		jobRec, err := w.service.GetJob(context.Background(), parsedID)
		if err != nil {
			log.Printf("worker-%d error fetching job from db: %v\n", w.id, err)
			continue
		}

		// 3. Mark job as processing
		err = w.service.UpdateJobStatus(context.Background(), parsedID, "processing")
		if err != nil {
			log.Printf("worker-%d error updating job status to processing: %v\n", w.id, err)
			continue
		}

		log.Printf("worker-%d processing job: %s\n", w.id, jobID)
		time.Sleep(1 * time.Second) // Simulate work

		// ============================================
		// STEP 4: Simulate Job Failure (~50% chance)
		// ============================================
		var execErr error
		if rand.Intn(2) == 0 {
			execErr = errors.New("simulated random failure")
		}

		// ============================================
		// STEP 5: Retry Logic & Exponential Backoff
		// ============================================
		if execErr != nil {
			log.Printf("worker-%d job failed: %v\n", w.id, execErr)
			
			retryCount := jobRec.RetryCount + 1

			// STEP 8: Check Max Retries Limit
			if retryCount > jobRec.MaxRetries {
				log.Printf("worker-%d max retries exceeded for job %s. Marking permanently failed.\n", w.id, jobID)
				_ = w.service.UpdateJobRetry(context.Background(), parsedID, retryCount, execErr.Error(), nil, "failed")
				continue
			}

			// STEP 6: Calculate Exponential Backoff Delay
			delay := time.Duration(math.Pow(2, float64(retryCount))) * time.Second
			nextRunAt := time.Now().Add(delay)

			log.Printf("retry attempt %d in %v\n", retryCount, delay)

			// STEP 9: Persist Last Error and Next Retry Info
			err = w.service.UpdateJobRetry(context.Background(), parsedID, retryCount, execErr.Error(), &nextRunAt, "pending")
			if err != nil {
				log.Printf("worker-%d error updating retry status: %v\n", w.id, err)
			}

			// STEP 7: Non-Blocking Delayed Retry (Using Redis ZSET)
			err = w.queue.EnqueueDelayed(context.Background(), jobID, nextRunAt)
			if err != nil {
				log.Printf("worker-%d error enqueuing delayed job: %v\n", w.id, err)
			}
			continue
		}

		// 5. Mark completed if success path is hit
		err = w.service.UpdateJobStatus(context.Background(), parsedID, "completed")
		if err != nil {
			log.Printf("worker-%d error updating job status to completed: %v\n", w.id, err)
			continue
		}

		log.Printf("worker-%d job completed: %s\n", w.id, jobID)
	}
}