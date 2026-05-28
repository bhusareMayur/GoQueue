package scheduler

import (
	"context"
	"log/slog"
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
	schedLogger := slog.With("component", "delayed_scheduler")
	schedLogger.Info("delayed job scheduler started")
	
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			schedLogger.Info("delayed job scheduler shutting down")
			return
		case <-ticker.C:
			s.processDelayedJobs(ctx, schedLogger)
		}
	}
}

func (s *DelayedScheduler) processDelayedJobs(ctx context.Context, logger *slog.Logger) {
	now := time.Now().Unix()

	opt := &goredis.ZRangeBy{
		Min: "-inf",
		Max: strconv.FormatInt(now, 10),
	}

	jobs, err := s.client.ZRangeByScore(ctx, "delayed_jobs", opt).Result()
	if err != nil {
		logger.Error("error fetching delayed jobs", "error", err)
		time.Sleep(3 * time.Second)
		return
	}

	for _, member := range jobs {
		parts := strings.SplitN(member, ":", 2)
		jobID := parts[0]
		queueName := "jobs"
		
		if len(parts) == 2 {
			queueName = "jobs:" + parts[1]
		}

		pipe := s.client.TxPipeline()
		pipe.ZRem(ctx, "delayed_jobs", member)
		pipe.LPush(ctx, queueName, jobID)

		_, err := pipe.Exec(ctx)
		if err != nil {
			logger.Error("error moving job", "job_id", jobID, "error", err)
			continue
		}
		
		logger.Info("moved delayed job back to queue", "job_id", jobID, "target_queue", queueName)
	}
}