package dataredis

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/yoanesber/Go-Department-CRUD/pkg/contextdata/dbcontext"
	"github.com/yoanesber/Go-Department-CRUD/pkg/logger"
	"github.com/yoanesber/Go-Department-CRUD/pkg/util/redisutil"
)

// Interface for the DataRedisService
// This interface defines the methods that the DataRedisService should implement
type DataRedisService interface {
	GetStringValue(ctx context.Context, key string) (string, error)
	GetJSONValue(ctx context.Context, key string) (interface{}, error)
}

// This struct defines the DataRedisService
type dataRedisService struct{}

// NewDataRedisService creates a new instance of DataRedisService
// It initializes the dataRedisService struct and returns it.
func NewDataRedisService() DataRedisService {
	return &dataRedisService{}
}

// GetStringValue retrieves a string value from Redis by its key
func (s *dataRedisService) GetStringValue(ctx context.Context, key string) (string, error) {
	// Get the Redis client from the context
	redisClient := dbcontext.GetRedisClient(ctx)
	if redisClient == nil {
		logger.Error("redis client is nil")
		return "", errors.New("redis client is nil")
	}

	// Retrieve the string value from Redis
	value, err := redisutil.Get(ctx, redisClient, key)
	if err == redis.Nil {
		logger.Error("key does not exist in Redis")
		return "", errors.New("key does not exist in Redis")
	}

	if err != nil {
		logger.Error(fmt.Sprintf("failed to get string value from Redis: %v", err))
		return "", err
	}

	// Check if the value is empty
	if value == "" {
		logger.Error("value is empty")
		return "", errors.New("value is empty")
	}

	return value, nil
}

// GetJSONValue retrieves a JSON value from Redis by its key
func (s *dataRedisService) GetJSONValue(ctx context.Context, key string) (interface{}, error) {
	// Get the Redis client from the context
	redisClient := dbcontext.GetRedisClient(ctx)
	if redisClient == nil {
		logger.Error("redis client is nil")
		return nil, errors.New("redis client is nil")
	}

	// Retrieve the JSON value from Redis
	value, err := redisutil.GetJSON[any](ctx, redisClient, key)
	if err == redis.Nil {
		logger.Error("key does not exist in Redis")
		return "", errors.New("key does not exist in Redis")
	}

	if err != nil {
		logger.Error(fmt.Sprintf("failed to get JSON value from Redis: %v", err))
		return nil, err
	}

	return value, nil
}
