package worker

import (
	"context"
	"log"
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

func (w *Worker) Start() {
	log.Printf("worker-%d started", w.id)

	for {
		// 1. Consume job ID from Redis
		jobID, err := w.queue.Consume(context.Background())
		if err != nil {
			log.Printf("worker-%d consume error: %v", w.id, err)
			continue
		}

		log.Printf("worker-%d received job: %s", w.id, jobID)

		// Parse the UUID
		parsedID, err := uuid.Parse(jobID)
		if err != nil {
			log.Printf("worker-%d invalid uuid error: %v", w.id, err)
			continue
		}

		// 2. Fetch job from PostgreSQL
		_, err = w.service.GetJob(
			context.Background(),
			parsedID,
		)
		if err != nil {
			log.Printf("worker-%d error fetching job from db: %v", w.id, err)
			continue
		}

		// 3. Mark job as processing
		err = w.service.UpdateJobStatus(
			context.Background(),
			parsedID,
			"processing",
		)
		if err != nil {
			log.Printf("worker-%d error updating job status to processing: %v", w.id, err)
			continue
		}

		log.Printf("worker-%d processing job: %s", w.id, jobID)

		// 4. Execute job (Simulating work)
		time.Sleep(2 * time.Second)

		// 5. Mark completed
		err = w.service.UpdateJobStatus(
			context.Background(),
			parsedID,
			"completed",
		)
		if err != nil {
			log.Printf("worker-%d error updating job status to completed: %v", w.id, err)
			continue
		}

		log.Printf("worker-%d job completed: %s", w.id, jobID)
	}
}