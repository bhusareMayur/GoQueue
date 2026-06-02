package redis

import (
	goredis "github.com/redis/go-redis/v9"
)

func NewClient(
	addr string,
) *goredis.Client {

	// BOTTLENECK 3 FIX: Redis Connection Pooling
	// Reusing connections drastically reduces network handshakes under heavy load.
	return goredis.NewClient(
		&goredis.Options{
			Addr:         addr,
			PoolSize:     100, // Maximum number of socket connections
			MinIdleConns: 20,  // Minimum number of idle connections kept alive
		},
	)
}