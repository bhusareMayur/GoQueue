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
	"github.com/bhusareMayur/goqueue/internal/observability/metrics"
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

		jobRec, err := w.service.GetJob(context.Background(), parsedID)
		if err != nil {
			log.Printf("%s error fetching job from db: %v\n", workerName, err)
			continue
		}

		err = w.service.ClaimJob(context.Background(), parsedID, workerName)
		if err != nil {
			log.Printf("%s error claiming job: %v\n", workerName, err)
			continue
		}

		log.Printf("%s processing job: %s\n", workerName, jobID)
		
		// =====================================
		// NEW: Start timing the execution
		// =====================================
		start := time.Now() 
		time.Sleep(1 * time.Second) 

		var execErr error
		if rand.Intn(2) == 0 {
			execErr = errors.New("simulated random failure")
		}

		// =====================================
		// NEW: Record Latency
		// =====================================
		duration := time.Since(start).Seconds()
		metrics.JobProcessingDuration.WithLabelValues(jobRec.Priority).Observe(duration)

		if execErr != nil {
			// NEW: Record failure metric
			metrics.JobsProcessed.WithLabelValues("failed", jobRec.Priority, workerName).Inc()
			
			log.Printf("%s job failed: %v\n", workerName, execErr)
			
			retryCount := jobRec.RetryCount + 1

			if retryCount > jobRec.MaxRetries {
				log.Printf("%s max retries exceeded for job %s. Moving to DLQ.\n", workerName, jobID)
				
				err = w.service.MoveToDLQ(context.Background(), parsedID, execErr.Error())
				if err != nil {
					log.Printf("%s error moving job to DLQ: %v\n", workerName, err)
				}
				continue
			}

			delay := time.Duration(math.Pow(2, float64(retryCount))) * time.Second
			nextRunAt := time.Now().Add(delay)

			log.Printf("retry attempt %d in %v\n", retryCount, delay)

			err = w.service.UpdateJobRetry(context.Background(), parsedID, retryCount, execErr.Error(), &nextRunAt, "pending")
			if err != nil {
				log.Printf("%s error updating retry status: %v\n", workerName, err)
			}

			err = w.queue.EnqueueDelayed(context.Background(), jobID, jobRec.Priority, nextRunAt)
			if err != nil {
				log.Printf("%s error enqueuing delayed job: %v\n", workerName, err)
			}
			continue
		}

		err = w.service.UpdateJobStatus(context.Background(), parsedID, "completed")
		if err != nil {
			log.Printf("%s error updating job status to completed: %v\n", workerName, err)
			continue
		}

		// NEW: Record success metric
		metrics.JobsProcessed.WithLabelValues("success", jobRec.Priority, workerName).Inc()
		log.Printf("%s job completed: %s\n", workerName, jobID)
	}
}