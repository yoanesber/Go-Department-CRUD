package redisutil

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// PushToList pushes a value to a Redis list with a specified key.
// It adds the value to the head of the list.
func PushToList(ctx context.Context, client *redis.Client, key string, value string) error {
	return client.LPush(ctx, key, value).Err()
}

// GetListRange retrieves a range of values from a Redis list with a specified key.
// It returns a slice of strings representing the values in the specified range.
func GetListRange(ctx context.Context, client *redis.Client, key string, start int64, stop int64) ([]string, error) {
	values, err := client.LRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	return values, nil
}

// PopFromList pops a value from a Redis list with a specified key.
// It removes the value from the head of the list and returns the updated list.
// If the list is empty, it returns an empty slice.
func PopFromList(ctx context.Context, client *redis.Client, key string) ([]string, error) {
	_, err := client.LPop(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// Get the updated list after popping the value
	updatedList, err := client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return updatedList, nil
}
