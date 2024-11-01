package database

import (
	"context"
	"link-base/internal/config"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient initializes and returns a Redis client using the provided configuration.
//
// Parameters:
//   - cfg: A RedisConfig struct containing the Redis connection details such as address, password,
//     database number, dial timeout, read timeout, write timeout, pool size, and minimum idle connections.
//
// Returns:
//   - *redis.Client: A pointer to the initialized Redis client.
//   - error: An error if the connection to Redis fails.
func NewRedisClient(cfg config.RedisConfig) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:         cfg.Address,
		DB:           cfg.DatabaseNumber,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	}

	client := redis.NewClient(opts)

	_, err := client.Ping(context.TODO()).Result()
	return client, err
}
