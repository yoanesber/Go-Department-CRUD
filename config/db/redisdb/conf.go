package redisdb

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/yoanesber/Go-Department-CRUD/pkg/logger"

	"github.com/go-redis/redis/v8" // Redis client for Go
)

var (
	RedisClient *redis.Client
	RedisDB     string
	RedisHost   string
	RedisPort   string
	RedisUser   string
	RedisPass   string
)

// LoadEnv loads Redis configuration from environment variables
func LoadEnv() {
	RedisDB = os.Getenv("REDIS_DB")
	RedisHost = os.Getenv("REDIS_HOST")
	RedisPort = os.Getenv("REDIS_PORT")
	RedisUser = os.Getenv("REDIS_USER")
	RedisPass = os.Getenv("REDIS_PASS")
}

// InitRedis initializes the Redis client using environment variables
// It constructs the connection string and calls ConnectRedis to establish the connection
func InitRedis() {
	// Initialize the Redis client
	redisDb, _ := strconv.Atoi(RedisDB)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", RedisHost, RedisPort),
		Username: RedisUser,
		Password: RedisPass,
		DB:       redisDb,
		// DialTimeout:        10 * time.Second,
		// ReadTimeout:        30 * time.Second,
		// WriteTimeout:       30 * time.Second,
		// PoolSize:           10,
		// PoolTimeout:        30 * time.Second,
		// IdleTimeout:        500 * time.Millisecond,
		// IdleCheckFrequency: 500 * time.Millisecond,
		// TLSConfig: &tls.Config{
		// 	InsecureSkipVerify: true,
		// },
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to Redis: %v", err))
		return
	}

	logger.Info("Connected to Redis")
}

// GetRedisClient returns the Redis client instance
func GetRedisClient() *redis.Client {
	return RedisClient
}
