package redisutil

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// Set sets a string value in Redis with a specified key and TTL.
func Set(ctx context.Context, client *redis.Client, key string, value string, ttl time.Duration) error {
	return client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a string value from Redis with a specified key.
func Get(ctx context.Context, client *redis.Client, key string) (string, error) {
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// SetJSON sets a JSON value in Redis with a specified key and TTL.
// It marshals the value into JSON format and stores it in Redis.
func SetJSON(ctx context.Context, client *redis.Client, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return client.Set(ctx, key, data, ttl).Err()
}

// GetJSON retrieves a JSON value from Redis with a specified key.
// It unmarshals the JSON data into the provided value.
func GetJSON[T any](ctx context.Context, client *redis.Client, key string) (*T, error) {
	data, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteKey deletes a key from Redis.
func DeleteKey(ctx context.Context, client *redis.Client, key string) error {
	return client.Del(ctx, key).Err()
}
