package worker

import (
	"context"
	"log"
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

func NewWorker(
	id int,
	queue job.Queue,
	service *job.Service,
) *Worker {
	return &Worker{
		id:      id,
		queue:   queue,
		service: service,
	}
}

// STEP 4: Worker Accepts Context (and WaitGroup)
func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	// Let main function know this worker is done when the function exits
	defer wg.Done()
	log.Printf("worker-%d started\n", w.id)

	for {
		// STEP 5: Add Context Check in Loop
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

		// Handle the 5-second BRPop timeout cleanly
		if jobID == "" {
			continue
		}

		log.Printf("worker-%d received job: %s\n", w.id, jobID)

		// Parse the UUID
		parsedID, err := uuid.Parse(jobID)
		if err != nil {
			log.Printf("worker-%d invalid uuid error: %v\n", w.id, err)
			continue
		}

		// 2. Fetch job from PostgreSQL
		_, err = w.service.GetJob(
			context.Background(),
			parsedID,
		)
		if err != nil {
			log.Printf("worker-%d error fetching job from db: %v\n", w.id, err)
			continue
		}

		// 3. Mark job as processing
		err = w.service.UpdateJobStatus(
			context.Background(),
			parsedID,
			"processing",
		)
		if err != nil {
			log.Printf("worker-%d error updating job status to processing: %v\n", w.id, err)
			continue
		}

		log.Printf("worker-%d processing job: %s\n", w.id, jobID)

		// 4. Execute job (Simulating work)
		time.Sleep(2 * time.Second)

		// 5. Mark completed
		err = w.service.UpdateJobStatus(
			context.Background(),
			parsedID,
			"completed",
		)
		if err != nil {
			log.Printf("worker-%d error updating job status to completed: %v\n", w.id, err)
			continue
		}

		log.Printf("worker-%d job completed: %s\n", w.id, jobID)
	}
}