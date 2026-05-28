package reaper

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/bhusareMayur/goqueue/internal/domain/job"
)

type Reaper struct {
	service *job.Service
	queue   job.Queue
}

func NewReaper(service *job.Service, queue job.Queue) *Reaper {
	return &Reaper{
		service: service,
		queue:   queue,
	}
}

func (r *Reaper) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	
	reaperLogger := slog.With("component", "reaper")
	reaperLogger.Info("reaper service started")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			reaperLogger.Info("reaper service shutting down")
			return
		case <-ticker.C:
			r.processStuckJobs(ctx, reaperLogger)
		}
	}
}

func (r *Reaper) processStuckJobs(ctx context.Context, logger *slog.Logger) {
	stuckJobs, err := r.service.GetStuckJobs(ctx, 30*time.Second)
	if err != nil {
		logger.Error("error fetching stuck jobs", "error", err)
		return
	}

	for _, j := range stuckJobs {
		workerID := "unknown"
		if j.WorkerID != nil {
			workerID = *j.WorkerID
		}
		
		jobLogger := logger.With("job_id", j.ID)
		if j.CorrelationID != nil {
			jobLogger = jobLogger.With("correlation_id", *j.CorrelationID)
		}
		
		jobLogger.Warn("detected stuck job", "crashed_worker", workerID)

		newRetryCount := j.RetryCount + 1
		if newRetryCount > j.MaxRetries {
			jobLogger.Warn("max retries exceeded, moving to DLQ")
			err := r.service.MoveToDLQ(ctx, j.ID, "visibility timeout exceeded")
			if err != nil {
				jobLogger.Error("error moving job to DLQ", "error", err)
			}
			continue
		}

		err = r.service.RequeueStuckJob(ctx, j.ID)
		if err != nil {
			jobLogger.Error("error requeueing job in db", "error", err)
			continue
		}

		err = r.queue.Enqueue(ctx, j.ID.String(), j.Priority)
		if err != nil {
			jobLogger.Error("error enqueueing job to redis", "error", err)
			continue
		}

		jobLogger.Info("successfully requeued stuck job")
	}
}