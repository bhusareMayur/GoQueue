package reaper

import (
	"context"
	"log"
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
	log.Println("reaper service started")

	// Tick every 10 seconds to check for stuck jobs
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("reaper service shutting down")
			return
		case <-ticker.C:
			r.processStuckJobs(ctx)
		}
	}
}

func (r *Reaper) processStuckJobs(ctx context.Context) {
	// Define timeout duration (e.g., job processing longer than 30s is stuck)
	stuckJobs, err := r.service.GetStuckJobs(ctx, 30*time.Second)
	if err != nil {
		log.Printf("reaper error fetching stuck jobs: %v\n", err)
		return
	}

	for _, j := range stuckJobs {
		workerID := "unknown"
		if j.WorkerID != nil {
			workerID = *j.WorkerID
		}
		
		log.Printf("reaper detected stuck job %s (crashed worker: %s)\n", j.ID, workerID)

		// Check if it has exceeded retries
		newRetryCount := j.RetryCount + 1
		if newRetryCount > j.MaxRetries {
			log.Printf("reaper moving job %s to DLQ (max retries exceeded)\n", j.ID)
			err := r.service.MoveToDLQ(ctx, j.ID, "visibility timeout exceeded")
			if err != nil {
				log.Printf("reaper error moving job %s to DLQ: %v\n", j.ID, err)
			}
			continue
		}

		// Update DB state
		err = r.service.RequeueStuckJob(ctx, j.ID)
		if err != nil {
			log.Printf("reaper error requeueing job %s in db: %v\n", j.ID, err)
			continue
		}

		// Push back to Redis Queue so another worker can pick it up
		err = r.queue.Enqueue(ctx, j.ID.String())
		if err != nil {
			log.Printf("reaper error enqueueing job %s to redis: %v\n", j.ID, err)
			continue
		}

		log.Printf("reaper successfully requeued job %s\n", j.ID)
	}
}