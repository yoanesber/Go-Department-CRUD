package redisutil

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// AddToSet adds one or more members to a Redis Set
// If the key does not exist, it will be created.
func AddToSet(ctx context.Context, client *redis.Client, key string, members ...string) error {
	return client.SAdd(ctx, key, members).Err()
}

// GetSetMembers retrieves all members of a Redis Set
// It returns a slice of strings representing the members of the set.
func GetSetMembers(ctx context.Context, client *redis.Client, key string) ([]string, error) {
	return client.SMembers(ctx, key).Result()
}
