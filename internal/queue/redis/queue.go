package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Queue struct {
	client *goredis.Client
}

func NewQueue(
	client *goredis.Client,
) *Queue {

	return &Queue{
		client: client,
	}
}

// NEW: Backpressure support - gets total length of all job queues
func (q *Queue) GetQueueLength(
	ctx context.Context,
) (int64, error) {
	queues := []string{"jobs:high", "jobs:medium", "jobs:low", "jobs"}
	
	pipe := q.client.Pipeline()
	var cmds []*goredis.IntCmd
	
	for _, queue := range queues {
		cmds = append(cmds, pipe.LLen(ctx, queue))
	}
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	
	var total int64
	for _, cmd := range cmds {
		total += cmd.Val()
	}
	
	return total, nil
}

func (q *Queue) Enqueue(
	ctx context.Context,
	jobID string,
	priority string,
) error {
	
	// Default to 'jobs' if priority is not specified, otherwise use 'jobs:<priority>'
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

	// STRICT PRIORITY SCHEDULING:
	// BRPOP checks keys in order. It will always drain jobs:high before touching jobs:medium.
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

	// result[0] is the queue name it popped from, result[1] is the job ID
	return result[1], nil
}

func (q *Queue) EnqueueDelayed(
	ctx context.Context,
	jobID string,
	priority string,
	runAt time.Time,
) error {
	
	// Embed priority into the member string so the scheduler knows where to push it later
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