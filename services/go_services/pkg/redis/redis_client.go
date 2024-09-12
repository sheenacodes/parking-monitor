package redis

import (
	"context"
	"fmt"
	"time"

	"go_services/pkg/logger"

	"github.com/redis/go-redis/v9"
)

const (
	maxRetries     = 5                // Maximum number of retries before giving up
	initialBackoff = 2 * time.Second  // Initial delay before retrying
	maxBackoff     = 30 * time.Second // Maximum delay between retries
)

// RedisClient is a wrapper around the redis.Client to hold the instance
type RedisClient struct {
	Client *redis.Client
}

// Function to connect to Redis with retry and exponential backoff
func GetRedisClient(addr string, pword string, database int) (*RedisClient, error) {
	var client *redis.Client
	var err error

	for retries := 0; retries < maxRetries; retries++ {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pword,
			DB:       database,
		})
		logger.Log.Info().Msg(" connecting to Redis")
		ctx := context.Background()
		_, err = client.Ping(ctx).Result()
		if err == nil {
			// Successful connection
			logger.Log.Info().Msg("Successfully connected to Redis")
			return &RedisClient{Client: client}, nil
		} else {
			backoff := time.Duration((1 << retries) * int(initialBackoff))
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			logger.Log.Warn().Err(err).Msgf("Failed to connect to Redis, retrying in %v...", backoff)
			time.Sleep(backoff)
		}
	}

	logger.Log.Fatal().Err(err).Msg("Failed to connect to Redis after multiple attempts")

	return nil, fmt.Errorf("failed to connect to Redis after %d attempts: %v", maxRetries, err)

}
