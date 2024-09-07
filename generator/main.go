package main

import (
	"context"
	"fmt"
	"generator/config"
	"generator/event"
	"generator/logger"
	"generator/rabbitmq"

	"math/rand"
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

func (r *RedisClient) IsSetNotEmpty() (bool, error) {
	// Use the SCARD command to get the number of members in the set
	ctx := context.Background()
	card, err := r.Client.SCard(ctx, "vehicles_in_parking").Result()
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

// GetRandomVehiclePlate retrieves a random vehicle plate from a Redis set
func (r *RedisClient) GetRandomVehiclePlateFromParkedSet() (string, error) {
	// Use SRANDMEMBER to get a random member from the set
	ctx := context.Background()
	plate, err := r.Client.SRandMember(ctx, "vehicles_in_parking").Result()
	if err != nil {
		return "", fmt.Errorf("could not get random member from set: %w", err)
	}
	return plate, nil
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

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.LogLevel)

	// Connect to RabbitMQ
	conn := rabbitmq.ConnectToRabbitMQ(cfg.RabbitMQURL)
	defer conn.Close()

	//ctx := context.Background()

	logger.Log.Info().Msg("Successfully connected to Redis")

	redisClient, err := connectToRedis(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error connecting to Redis:")
	}
	defer func() {
		if err := redisClient.Client.Close(); err != nil {
			logger.Log.Fatal().Err(err).Msg("Failed to close Redis client")
		}
	}()
	logger.Log.Info().Msg("Successfully connected to Redis")

	logger.Log.Info().Msg(cfg.GeneratorMode)
	if cfg.GeneratorMode == "entry" {

		for {
			eventPayload := event.GenerateEntryEvent()
			err := rabbitmq.PublishEvent(conn, cfg.QueueName, eventPayload)
			if err != nil {
				logger.Log.Fatal().Err(err).Msg("Failed to publish event")
			}

			err = redisClient.AddVehicleEntry(eventPayload.VehiclePlate)
			if err != nil {
				logger.Log.Fatal().Err(err).Msg("Error adding vehicle to redis set")
			} else {
				logger.Log.Debug().Msg("Vehicle entry added to redis set successfully")
			}

			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		}

	} else {

		exitRatio := 0.8

		for {
			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

			registeredCarAvailableForExit, err := redisClient.IsSetNotEmpty()
			if err != nil {
				logger.Log.Fatal().Err(err).Msg("Error checking set")
			}

			eventPayload := event.GenerateExitEvent()

			randomProbability := rand.Float64()
			logger.Log.Debug().Msgf("random prob %f", randomProbability)

			if randomProbability <= exitRatio && registeredCarAvailableForExit {

				parkedVehiclePlate, err := redisClient.GetRandomVehiclePlateFromParkedSet()
				logger.Log.Debug().Msgf("random parked vehicle plate %s", parkedVehiclePlate)

				if err == nil {

					eventPayload.VehiclePlate = parkedVehiclePlate
					err = rabbitmq.PublishEvent(conn, cfg.QueueName, eventPayload)
					if err != nil {
						logger.Log.Fatal().Err(err).Msg("Failed to publish event")
					}
					err = redisClient.RemoveVehicleEntry(eventPayload.VehiclePlate)
					if err != nil {
						logger.Log.Fatal().Err(err).Msg("Error removing vehicle from redis set")
					} else {
						logger.Log.Debug().Msg("Vehicle plate removed redis set successfully")
					}

				}
			} else {

				err = rabbitmq.PublishEvent(conn, cfg.QueueName, eventPayload)
				if err != nil {
					logger.Log.Fatal().Err(err).Msg("Failed to publish event")
				}

			}

		}

	}

}
