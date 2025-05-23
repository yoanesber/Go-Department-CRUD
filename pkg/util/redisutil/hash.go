package redisutil

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// SetHashField sets a field in a Redis hash with a specified key and value.
// It adds the field to the hash if it doesn't exist, or updates it if it does.
func SetHashField(ctx context.Context, client *redis.Client, key, field, value string) error {
	return client.HSet(ctx, key, field, value).Err()
}

// GetHashField retrieves a field from a Redis hash with a specified key.
// It returns the value of the field if it exists, or an error if it doesn't.
func GetHashField(ctx context.Context, client *redis.Client, key, field string) (string, error) {
	return client.HGet(ctx, key, field).Result()
}

// GetAllHash retrieves all fields and values from a Redis hash with a specified key.
// It returns a map of field-value pairs.
func GetAllHash(ctx context.Context, client *redis.Client, key string) (map[string]string, error) {
	return client.HGetAll(ctx, key).Result()
}
