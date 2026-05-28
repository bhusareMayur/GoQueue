package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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
	workerName := fmt.Sprintf("worker-%d", w.id)
	
	// Create a worker-specific logger context
	workerLogger := slog.With("worker_id", workerName)
	workerLogger.Info("worker started")

	for {
		select {
		case <-ctx.Done():
			workerLogger.Info("worker shutting down")
			return
		default:
		}

		jobID, err := w.queue.Consume(ctx)
		if err != nil {
			workerLogger.Error("consume error", "error", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if jobID == "" {
			continue
		}

		parsedID, err := uuid.Parse(jobID)
		if err != nil {
			workerLogger.Error("invalid uuid error", "job_id", jobID, "error", err)
			continue
		}

		jobRec, err := w.service.GetJob(context.Background(), parsedID)
		if err != nil {
			workerLogger.Error("error fetching job from db", "job_id", jobID, "error", err)
			continue
		}

		// Inject correlation ID into all subsequent logs for this job
		jobLogger := workerLogger.With("job_id", jobID)
		if jobRec.CorrelationID != nil {
			jobLogger = jobLogger.With("correlation_id", *jobRec.CorrelationID)
		}

		err = w.service.ClaimJob(context.Background(), parsedID, workerName)
		if err != nil {
			jobLogger.Error("error claiming job", "error", err)
			continue
		}

		jobLogger.Info("processing job started")
		
		start := time.Now() 
		time.Sleep(1 * time.Second) 

		var execErr error
		if rand.Intn(2) == 0 {
			execErr = errors.New("simulated random failure")
		}

		duration := time.Since(start).Seconds()
		metrics.JobProcessingDuration.WithLabelValues(jobRec.Priority).Observe(duration)

		if execErr != nil {
			metrics.JobsProcessed.WithLabelValues("failed", jobRec.Priority, workerName).Inc()
			
			jobLogger.Error("job execution failed", "error", execErr)
			
			retryCount := jobRec.RetryCount + 1

			if retryCount > jobRec.MaxRetries {
				jobLogger.Warn("max retries exceeded, moving to DLQ")
				err = w.service.MoveToDLQ(context.Background(), parsedID, execErr.Error())
				if err != nil {
					jobLogger.Error("error moving job to DLQ", "error", err)
				}
				continue
			}

			delay := time.Duration(math.Pow(2, float64(retryCount))) * time.Second
			nextRunAt := time.Now().Add(delay)

			jobLogger.Info("scheduling retry", "retry_attempt", retryCount, "delay_seconds", delay.Seconds())

			err = w.service.UpdateJobRetry(context.Background(), parsedID, retryCount, execErr.Error(), &nextRunAt, "pending")
			if err != nil {
				jobLogger.Error("error updating retry status", "error", err)
			}

			err = w.queue.EnqueueDelayed(context.Background(), jobID, jobRec.Priority, nextRunAt)
			if err != nil {
				jobLogger.Error("error enqueuing delayed job", "error", err)
			}
			continue
		}

		err = w.service.UpdateJobStatus(context.Background(), parsedID, "completed")
		if err != nil {
			jobLogger.Error("error updating job status to completed", "error", err)
			continue
		}

		metrics.JobsProcessed.WithLabelValues("success", jobRec.Priority, workerName).Inc()
		jobLogger.Info("job completed successfully", "duration_seconds", duration)
	}
}