package scheduler

import (
	"context"
	"log"
	"strconv"
	"strings"
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

	opt := &goredis.ZRangeBy{
		Min: "-inf",
		Max: strconv.FormatInt(now, 10),
	}

	jobs, err := s.client.ZRangeByScore(ctx, "delayed_jobs", opt).Result()
	if err != nil {
		log.Printf("scheduler error fetching delayed jobs: %v\n", err)
		// NEW: Add backoff to the scheduler. 
		// If Redis is down, wait a few seconds before the next tick 
		// to prevent go-redis internal connection pool spam.
		time.Sleep(3 * time.Second)
		return
	}

	for _, member := range jobs {
		// Split the member string to extract JobID and Priority
		parts := strings.SplitN(member, ":", 2)
		jobID := parts[0]
		queueName := "jobs" // default
		
		if len(parts) == 2 {
			queueName = "jobs:" + parts[1]
		}

		pipe := s.client.TxPipeline()
		pipe.ZRem(ctx, "delayed_jobs", member)
		pipe.LPush(ctx, queueName, jobID) // Push to correct priority queue

		_, err := pipe.Exec(ctx)
		if err != nil {
			log.Printf("scheduler error moving job %s: %v\n", jobID, err)
			continue
		}
		
		log.Printf("scheduler moved job %s back to %s\n", jobID, queueName)
	}
}