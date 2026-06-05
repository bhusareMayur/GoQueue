package publisher

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/bhusareMayur/goqueue/internal/domain/job"
	"github.com/bhusareMayur/goqueue/internal/observability/metrics"
)

type Publisher struct {
	service *job.Service
	queue   job.Queue
}

func NewPublisher(service *job.Service, queue job.Queue) *Publisher {
	return &Publisher{
		service: service,
		queue:   queue,
	}
}

func (p *Publisher) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	
	logger := slog.With("component", "outbox_publisher")
	logger.Info("outbox publisher started")

	// Poll frequently for low latency
	ticker := time.NewTicker(200 * time.Millisecond) 
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("outbox publisher shutting down")
			return
		case <-ticker.C:
			p.processOutbox(ctx, logger)
		}
	}
}

func (p *Publisher) processOutbox(ctx context.Context, logger *slog.Logger) {
	events, err := p.service.GetPendingOutboxEvents(ctx, 100)
	if err != nil {
		logger.Error("error fetching outbox events", "error", err)
		return
	}

	for _, event := range events {
		// 1. Publish to Redis
		err := p.queue.Enqueue(ctx, event.JobID.String(), event.Priority)
		if err != nil {
			logger.Error("error enqueueing job to redis", "job_id", event.JobID, "error", err)
			continue
		}

		// 2. Mark as Published
		err = p.service.MarkOutboxEventPublished(ctx, event.ID)
		if err != nil {
			logger.Error("error marking outbox event as published", "event_id", event.ID, "error", err)
			continue
		}
		
		metrics.JobsEnqueued.WithLabelValues(event.Priority).Inc()
		logger.Debug("published outbox event to redis", "job_id", event.JobID)
	}
}