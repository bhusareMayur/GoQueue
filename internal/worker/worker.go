package worker

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/bhusareMayur/goqueue/internal/domain/job"
)

type Worker struct {
	queue   job.Queue
	service *job.Service
}

func NewWorker(
	queue job.Queue,
	service *job.Service,
) *Worker {
	return &Worker{
		queue:   queue,
		service: service,
	}
}

func (w *Worker) Start() {
	log.Println("worker started")

	for {
		// 1. Consume job ID from Redis
		jobID, err := w.queue.Consume(context.Background())
		if err != nil {
			log.Println("consume error:", err)
			continue
		}

		log.Println("received job:", jobID)

		// Parse the UUID
		parsedID, err := uuid.Parse(jobID)
		if err != nil {
			log.Println("invalid uuid error:", err)
			continue
		}

		// 2. Fetch job from PostgreSQL
		_, err = w.service.GetJob(
			context.Background(),
			parsedID,
		)
		if err != nil {
			log.Println("error fetching job from db:", err)
			continue
		}

		// 3. Mark job as processing
		err = w.service.UpdateJobStatus(
			context.Background(),
			parsedID,
			"processing",
		)
		if err != nil {
			log.Println("error updating job status to processing:", err)
			continue
		}

		log.Println("processing job:", jobID)

		// 4. Execute job (Simulating work)
		time.Sleep(2 * time.Second)

		// 5. Mark completed
		err = w.service.UpdateJobStatus(
			context.Background(),
			parsedID,
			"completed",
		)
		if err != nil {
			log.Println("error updating job status to completed:", err)
			continue
		}

		log.Println("job completed:", jobID)
	}
}