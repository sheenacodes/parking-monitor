package main

import (
	"context"
	"encoding/json"
	"parking-record-keeper/config"
	"time"

	"github.com/sheenacodes/sharedutils/logger"
	"github.com/sheenacodes/sharedutils/rabbitmq"
	"github.com/sheenacodes/sharedutils/redis"
)

type EntryEvent struct {
	ID            string    `json:"id"`
	VehiclePlate  string    `json:"vehicle_plate"`
	EntryDateTime time.Time `json:"entry_date_time"`
}

type ExitEvent struct {
	ID           string    `json:"id"`
	VehiclePlate string    `json:"vehicle_plate"`
	ExitDateTime time.Time `json:"exit_date_time"`
}

type EntryEventProcessor struct {
	RedisClient *redis.RedisClient
}

func (eventprocessor *EntryEventProcessor) ProcessMessage(msgBody []byte) error {
	var payload EntryEvent
	if err := json.Unmarshal(msgBody, &payload); err != nil {
		return err
	}

	hashKey := payload.VehiclePlate
	fieldName := "entry_date_time"
	fieldValue := payload.EntryDateTime
	logger.Log.Debug().Msgf("key - %s; field - %s; value - %s ", hashKey, fieldName, fieldValue)
	// Store entry time

	err := eventprocessor.RedisClient.Client.HSet(context.Background(), hashKey, fieldName, fieldValue).Err()
	if err != nil {
		return err
	}
	logger.Log.Debug().Msgf(" added to redis: key - %s; field - %s; value - %s ", hashKey, fieldName, fieldValue)
	return nil
}

type ExitEventProcessor struct {
	RedisClient *redis.RedisClient
}

func (eventprocessor *ExitEventProcessor) ProcessMessage(msgBody []byte) error {
	var payload ExitEvent
	if err := json.Unmarshal(msgBody, &payload); err != nil {
		return err
	}

	hashKey := payload.VehiclePlate
	fieldName := "exit_date_time"
	fieldValue := payload.ExitDateTime
	logger.Log.Debug().Msgf("key - %s; field - %s; value - %s ", hashKey, fieldName, fieldValue)
	// Store entry time

	err := eventprocessor.RedisClient.Client.HSet(context.Background(), hashKey, fieldName, fieldValue).Err()
	if err != nil {
		return err
	}
	logger.Log.Debug().Msgf(" added to redis: key - %s; field - %s; value - %s ", hashKey, fieldName, fieldValue)

	return nil
}

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.LogLevel)

	rabbitMQClient, err := rabbitmq.GetRabbitMQClient(cfg.RabbitMQURL)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to initialize RabbitMQ client")
	}
	defer rabbitMQClient.Close()

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

	entryEvtProcessor := &EntryEventProcessor{
		RedisClient: redisClient,
	}

	// Handle Entry Events
	err = rabbitMQClient.ConsumeQueue(cfg.EntryQueueName, entryEvtProcessor)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to consume entry events")
	}

	logger.Log.Debug().Msg("Entry queue consumer set up")
	// Handle Exit Events

	exitEvtProcessor := &ExitEventProcessor{
		RedisClient: redisClient,
	}

	err = rabbitMQClient.ConsumeQueue(cfg.ExitQueueName, exitEvtProcessor)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to consume exit events")
	}

	logger.Log.Debug().Msg("Exit queue consumer set up")

	select {}
}
