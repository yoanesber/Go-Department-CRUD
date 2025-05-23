package dbcontext

import (
	"context"

	"github.com/go-redis/redis/v8" // Redis client for Go
	"gorm.io/gorm"
)

type dbCtxKey struct{}
type redisCtxKey struct{}

var dbKey = dbCtxKey{}
var redisKey = redisCtxKey{}

// InjectDB injects *gorm.DB into context
func InjectDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
}

// GetDB extracts *gorm.DB from context
func GetDB(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(dbKey).(*gorm.DB)
	if !ok {
		return nil
	}
	return db
}

// InjectRedis injects *redis.Client into context
func InjectRedisClient(ctx context.Context, db *redis.Client) context.Context {
	return context.WithValue(ctx, redisKey, db)
}

// GetRedis extracts *redis.Client from context
func GetRedisClient(ctx context.Context) *redis.Client {
	db, ok := ctx.Value(redisKey).(*redis.Client)
	if !ok {
		return nil
	}
	return db
}
