package redis

import (
	"context"

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

	result, err := q.client.BRPop(
		ctx,
		0,
		"jobs",
	).Result()

	if err != nil {
		return "", err
	}

	return result[1], nil
}