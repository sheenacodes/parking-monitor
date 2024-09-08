package redis

import (
	"context"
	"fmt"
	"generator/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	maxRetries     = 5                // Maximum number of retries before giving up
	initialBackoff = 2 * time.Second  // Initial delay before retrying
	maxBackoff     = 30 * time.Second // Maximum delay between retries
	redisSetName   = "vehicles_parked"
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
			logger.Log.Warn().Err(err).Msgf("Failed to connect to RabbitMQ, retrying in %v...", backoff)
			time.Sleep(backoff)
		}
	}

	//if err != nil {
	logger.Log.Fatal().Err(err).Msg("Failed to connect to RabbitMQ after multiple attempts")

	return nil, fmt.Errorf("failed to connect to Redis after %d attempts: %v", maxRetries, err)
	//}

}

func (r *RedisClient) IsSetNotEmpty() (bool, error) {
	// Use the SCARD command to get the number of members in the set
	ctx := context.Background()
	card, err := r.Client.SCard(ctx, redisSetName).Result()
	if err != nil {
		logger.Log.Error().Err(err).Msg("error checking set size")
		return false, fmt.Errorf("error checking set size: %v", err)
	}

	logger.Log.Debug().Msgf("%d Vehicles in Redis Set", card)
	// Return true if the number of members is greater than 0
	return card > 0, nil
}

// AddVehicleEntry adds a vehicle entry to the Redis list of vehicles that have entered but not exited
func (r *RedisClient) AddVehicleEntry(vehiclePlate string) error {
	ctx := context.Background()
	_, err := r.Client.SAdd(ctx, redisSetName, vehiclePlate).Result()
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
	_, err := r.Client.SRem(ctx, redisSetName, vehiclePlate).Result()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to remove vehicle entry from Redis")
		return err
	}
	logger.Log.Debug().Msgf("Vehicle %s removed from Redis", vehiclePlate)
	return nil
}

// GetRandomVehiclePlate retrieves a random vehicle plate from a Redis set
func (r *RedisClient) GetRandomVehiclePlateFromParkedSet() (string, error) {
	// Use SRANDMEMBER to get a random member from the set
	ctx := context.Background()
	plate, err := r.Client.SRandMember(ctx, redisSetName).Result()
	if err != nil {
		return "", fmt.Errorf("could not get random member from set: %w", err)
	}
	return plate, nil
}
