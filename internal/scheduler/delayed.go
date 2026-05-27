package scheduler

import (
	"context"
	"log"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type DelayedScheduler struct {
	client *goredis.Client
}

func NewDelayedScheduler(client *goredis.Client) *DelayedScheduler {
	return &DelayedScheduler{
		client: client,
	}
}

func (s *DelayedScheduler) Start(ctx context.Context) {
	log.Println("delayed job scheduler started")
	
	// Tick every 1 second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("delayed job scheduler shutting down")
			return
		case <-ticker.C:
			s.processDelayedJobs(ctx)
		}
	}
}

func (s *DelayedScheduler) processDelayedJobs(ctx context.Context) {
	now := time.Now().Unix()

	// Get all jobs where the timestamp is NOW or in the PAST
	opt := &goredis.ZRangeBy{
		Min: "-inf",
		Max: strconv.FormatInt(now, 10),
	}

	jobs, err := s.client.ZRangeByScore(ctx, "delayed_jobs", opt).Result()
	if err != nil {
		log.Printf("scheduler error fetching delayed jobs: %v\n", err)
		return
	}

	for _, jobID := range jobs {
		// Use a Redis Pipeline (Transaction) to atomically remove from delayed queue
		// and push to the main queue so we don't lose jobs if it crashes midway.
		pipe := s.client.TxPipeline()
		pipe.ZRem(ctx, "delayed_jobs", jobID)
		pipe.LPush(ctx, "jobs", jobID)

		_, err := pipe.Exec(ctx)
		if err != nil {
			log.Printf("scheduler error moving job %s: %v\n", jobID, err)
			continue
		}
		
		log.Printf("scheduler moved job %s back to main queue\n", jobID)
	}
}