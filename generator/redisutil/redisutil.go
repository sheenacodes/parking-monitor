package redisutil

import (
	"context"
	"fmt"
	"generator/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient is a wrapper around the redis.Client to hold the instance
type RedisClient struct {
	Client *redis.Client
}

const (
	maxRetries     = 5                // Maximum number of retries before giving up
	initialBackoff = 2 * time.Second  // Initial delay before retrying
	maxBackoff     = 30 * time.Second // Maximum delay between retries
)

// Function to connect to Redis with retry and exponential backoff
func connectToRedis(addr string, pword string, database int) (*RedisClient, error) {
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
			logger.Log.Warn().Err(err).Msgf("Failed to connect to RabbitMQ, retrying in %v...", backoff)
			time.Sleep(backoff)
		}
	}

	//if err != nil {
	logger.Log.Fatal().Err(err).Msg("Failed to connect to RabbitMQ after multiple attempts")

	return nil, fmt.Errorf("failed to connect to Redis after %d attempts: %v", maxRetries, err)
	//}

}

// AddVehicleEntry adds a vehicle entry to the Redis list of vehicles that have entered but not exited
func (r *RedisClient) AddVehicleEntry(vehiclePlate string) error {
	ctx := context.Background()
	_, err := r.Client.SAdd(ctx, "vehicles_in_parking", vehiclePlate).Result()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to add vehicle entry to Redis")
		return err
	}
	logger.Log.Debug().Msgf("Vehicle %s added to Redis", vehiclePlate)
	return nil
}

// RemoveVehicleEntry removes a vehicle entry from the Redis list of vehicles that have entered
func (r *RedisClient) RemoveVehicleEntry(vehiclePlate string) error {
	ctx := context.Background()
	_, err := r.Client.SRem(ctx, "vehicles_in_parking", vehiclePlate).Result()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to remove vehicle entry from Redis")
		return err
	}
	logger.Log.Debug().Msgf("Vehicle %s removed from Redis", vehiclePlate)
	return nil
}

// VehicleExists checks if a vehicle has already entered
func (r *RedisClient) VehicleExists(vehiclePlate string) (bool, error) {
	ctx := context.Background()
	exists, err := r.Client.SIsMember(ctx, "vehicles_in_parking", vehiclePlate).Result()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to check vehicle existence in Redis")
		return false, err
	}
	return exists, nil
}
