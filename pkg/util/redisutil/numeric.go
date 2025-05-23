package redisutil

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// Increment increases a key's value by 1 (or given amount)
// If the key does not exist, it will be created with the specified value.
func Increment(ctx context.Context, client *redis.Client, key string, by int64) (int64, error) {
	return client.IncrBy(ctx, key, by).Result()
}

// Decrement decreases a key's value by 1 (or given amount)
// If the key does not exist, it will be created with the specified value.
func Decrement(ctx context.Context, client *redis.Client, key string, by int64) (int64, error) {
	return client.DecrBy(ctx, key, by).Result()
}
