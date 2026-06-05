package redis

import (
	"context"
	"sync/atomic"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Queue struct {
	client       *goredis.Client
	cachedLength atomic.Int64
}

func NewQueue(
	client *goredis.Client,
) *Queue {

	q := &Queue{
		client: client,
	}

	// Start background goroutine to cache queue length every 100ms
	go q.startQueueLengthPinger()

	return q
}

// Background worker to poll Redis so the API doesn't block
func (q *Queue) startQueueLengthPinger() {
	ticker := time.NewTicker(100 * time.Millisecond)
	ctx := context.Background()
	
	for range ticker.C {
		queues := []string{"jobs:high", "jobs:medium", "jobs:low", "jobs"}
		
		pipe := q.client.Pipeline()
		var cmds []*goredis.IntCmd
		
		for _, queue := range queues {
			cmds = append(cmds, pipe.LLen(ctx, queue))
		}
		
		_, err := pipe.Exec(ctx)
		if err == nil {
			var total int64
			for _, cmd := range cmds {
				total += cmd.Val()
			}
			// Update the thread-safe counter
			q.cachedLength.Store(total)
		}
	}
}

// GetQueueLength now reads instantly from memory instead of hitting Redis
func (q *Queue) GetQueueLength(
	ctx context.Context,
) (int64, error) {
	return q.cachedLength.Load(), nil
}

func (q *Queue) Enqueue(
	ctx context.Context,
	jobID string,
	priority string,
) error {
	
	queueName := "jobs"
	if priority != "" && priority != "default" {
		queueName = "jobs:" + priority
	}

	return q.client.LPush(
		ctx,
		queueName,
		jobID,
	).Err()
}

func (q *Queue) Consume(
	ctx context.Context,
) (string, error) {

	result, err := q.client.BRPop(
		ctx,
		5*time.Second,
		"jobs:high", "jobs:medium", "jobs:low", "jobs",
	).Result()

	if err != nil {
		if err == goredis.Nil {
			return "", nil
		}
		return "", err
	}

	return result[1], nil
}

func (q *Queue) EnqueueDelayed(
	ctx context.Context,
	jobID string,
	priority string,
	runAt time.Time,
) error {
	
	member := jobID
	if priority != "" && priority != "default" {
		member = jobID + ":" + priority
	}

	return q.client.ZAdd(
		ctx,
		"delayed_jobs",
		goredis.Z{
			Score:  float64(runAt.Unix()),
			Member: member,
		},
	).Err()
}