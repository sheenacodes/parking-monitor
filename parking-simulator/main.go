package main

import (
	"math/rand"
	"parking-simulator/config"
	"parking-simulator/event"
	"time"

	"github.com/sheenacodes/sharedutils/logger"
	"github.com/sheenacodes/sharedutils/rabbitmq"
	"github.com/sheenacodes/sharedutils/redis"
)

const (
	redisSetName = "parked_vehicles"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.LogLevel)

	// Connect to RabbitMQ
	conn := rabbitmq.ConnectToRabbitMQ(cfg.RabbitMQURL)
	defer conn.Close()

	//ctx := context.Background()

	logger.Log.Info().Msg("Successfully connected to Redis")

	redisClient, err := redis.GetRedisClient(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)
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

			err = redisClient.AddItemToSet(eventPayload.VehiclePlate, redisSetName)
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

			registeredCarAvailableForExit, err := redisClient.IsSetNotEmpty(redisSetName)
			if err != nil {
				logger.Log.Fatal().Err(err).Msg("Error checking set")
			}

			eventPayload := event.GenerateExitEvent()

			randomProbability := rand.Float64()
			logger.Log.Debug().Msgf("random prob %f", randomProbability)

			if randomProbability <= exitRatio && registeredCarAvailableForExit {

				parkedVehiclePlate, err := redisClient.GetRandomItemFromSet(redisSetName)
				logger.Log.Debug().Msgf("random parked vehicle plate %s", parkedVehiclePlate)

				if err == nil {

					eventPayload.VehiclePlate = parkedVehiclePlate
					err = rabbitmq.PublishEvent(conn, cfg.QueueName, eventPayload)
					if err != nil {
						logger.Log.Fatal().Err(err).Msg("Failed to publish event")
					}
					err = redisClient.RemoveItemFromSet(eventPayload.VehiclePlate, redisSetName)
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
