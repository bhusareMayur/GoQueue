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

func (q *Queue) Enqueue(
	ctx context.Context,
	jobID string,
) error {

	return q.client.LPush(
		ctx,
		"jobs",
		jobID,
	).Err()
}

func (q *Queue) Consume(
	ctx context.Context,
) (string, error) {

	// STEP 6: Problem With BRPOP - Use 5 seconds timeout instead of 0
	result, err := q.client.BRPop(
		ctx,
		5*time.Second,
		"jobs",
	).Result()

	if err != nil {
		// If the error is simply a timeout (no jobs found in 5s), return empty string
		if err == goredis.Nil {
			return "", nil
		}
		return "", err
	}

	return result[1], nil
}