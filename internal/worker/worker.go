package worker

import (
	"context"
	"errors"
	"fmt"
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

	workerName := fmt.Sprintf("worker-%d", w.id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("%s shutting down\n", workerName)
			return
		default:
		}

		// 1. Consume job ID from Redis
		jobID, err := w.queue.Consume(ctx)
		if err != nil {
			log.Printf("%s consume error: %v\n", workerName, err)
			continue
		}

		if jobID == "" {
			continue
		}

		parsedID, err := uuid.Parse(jobID)
		if err != nil {
			log.Printf("%s invalid uuid error: %v\n", workerName, err)
			continue
		}

		// 2. Fetch job from PostgreSQL
		jobRec, err := w.service.GetJob(context.Background(), parsedID)
		if err != nil {
			log.Printf("%s error fetching job from db: %v\n", workerName, err)
			continue
		}

		// ============================================
		// 3. NEW: Claim job for Visibility Timeout
		// ============================================
		err = w.service.ClaimJob(context.Background(), parsedID, workerName)
		if err != nil {
			log.Printf("%s error claiming job: %v\n", workerName, err)
			continue
		}

		log.Printf("%s processing job: %s\n", workerName, jobID)
		
		// Simulate work
		time.Sleep(1 * time.Second) 

		// *To test the Reaper, uncomment this line below to simulate a HARD CRASH!*
		// if rand.Intn(3) == 0 { log.Fatalf("%s CRASHED mid-job!", workerName) }

		// 4. Simulate Job Failure (~50% chance)
		var execErr error
		if rand.Intn(2) == 0 {
			execErr = errors.New("simulated random failure")
		}

		// 5. Retry Logic & Exponential Backoff
		if execErr != nil {
			log.Printf("%s job failed: %v\n", workerName, execErr)
			
			retryCount := jobRec.RetryCount + 1

			// 8. Check Max Retries Limit & Move to DLQ
			if retryCount > jobRec.MaxRetries {
				log.Printf("%s max retries exceeded for job %s. Moving to DLQ.\n", workerName, jobID)
				
				err = w.service.MoveToDLQ(context.Background(), parsedID, execErr.Error())
				if err != nil {
					log.Printf("%s error moving job to DLQ: %v\n", workerName, err)
				}
				continue
			}

			// 6. Calculate Exponential Backoff Delay
			delay := time.Duration(math.Pow(2, float64(retryCount))) * time.Second
			nextRunAt := time.Now().Add(delay)

			log.Printf("retry attempt %d in %v\n", retryCount, delay)

			// 9. Persist Last Error and Next Retry Info
			err = w.service.UpdateJobRetry(context.Background(), parsedID, retryCount, execErr.Error(), &nextRunAt, "pending")
			if err != nil {
				log.Printf("%s error updating retry status: %v\n", workerName, err)
			}

			// 7. Non-Blocking Delayed Retry (Using Redis ZSET)
			err = w.queue.EnqueueDelayed(context.Background(), jobID, nextRunAt)
			if err != nil {
				log.Printf("%s error enqueuing delayed job: %v\n", workerName, err)
			}
			continue
		}

		// 5. Mark completed if success path is hit
		err = w.service.UpdateJobStatus(context.Background(), parsedID, "completed")
		if err != nil {
			log.Printf("%s error updating job status to completed: %v\n", workerName, err)
			continue
		}

		log.Printf("%s job completed: %s\n", workerName, jobID)
	}
}