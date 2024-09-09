package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"parking-record-keeper/config"

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

type ParkingLog struct {
	VehiclePlate  string    `json:"vehicle_plate"`
	ExitDateTime  time.Time `json:"exit_date_time"`
	EntryDateTime time.Time `json:"entry_date_time"`
	Duration      string    `json:"duration"`
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
	logger.Log.Debug().Msgf("key - %s; field - %s; value - %s", hashKey, fieldName, fieldValue)

	// Store entry time
	if err := eventprocessor.RedisClient.AddFieldToHash(hashKey, fieldName, fieldValue); err != nil {
		return err
	}

	logger.Log.Debug().Msgf("Added to Redis: key - %s; field - %s; value - %s", hashKey, fieldName, fieldValue)
	return nil
}

type ExitEventProcessor struct {
	RedisClient *redis.RedisClient
}

func createJsonBody(log ParkingLog) ([]byte, error) {
	return json.Marshal(log)
}

func generateParkingSummary(vehiclePlate string, exitDateTime time.Time, rClient *redis.RedisClient) (*ParkingLog, error) {

	// Retrieve entry time
	fieldName := "entry_date_time"
	layout := time.RFC3339
	entryDateTime, err := rClient.GetFieldAsTime(vehiclePlate, fieldName, layout)
	if err != nil {
		return nil, fmt.Errorf("error retrieving entry time: %v", err)
	}

	parkingDuration := exitDateTime.Sub(entryDateTime).String()
	logger.Log.Debug().Msgf("Duration: %v", parkingDuration)

	return &ParkingLog{
		VehiclePlate:  vehiclePlate,
		EntryDateTime: entryDateTime,
		ExitDateTime:  exitDateTime,
		Duration:      parkingDuration,
	}, nil

}

func (eventprocessor *ExitEventProcessor) ProcessMessage(msgBody []byte) error {
	var payload ExitEvent
	if err := json.Unmarshal(msgBody, &payload); err != nil {
		return err
	}

	hashKey := payload.VehiclePlate
	fieldName := "exit_date_time"
	fieldValue := payload.ExitDateTime
	logger.Log.Debug().Msgf("key - %s; field - %s; value - %s", hashKey, fieldName, fieldValue)

	// Store exit time
	if err := eventprocessor.RedisClient.AddFieldToHash(hashKey, fieldName, fieldValue); err != nil {
		logger.Log.Error().Err(err).Msg("Failed writing to Redis")
		return err
	}

	logger.Log.Debug().Msgf("Added to Redis: key - %s; field - %s; value - %s", hashKey, fieldName, fieldValue)

	parkingLog, err := generateParkingSummary(payload.VehiclePlate, payload.ExitDateTime, eventprocessor.RedisClient)
	if err != nil {
		return fmt.Errorf("error creating Parking Log: %v", err)
	}

	// Marshal ParkingLog to JSON
	jsonBody, err := createJsonBody(*parkingLog)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Send HTTP POST request to REST API server
	resp, err := http.Post("http://python-server:8000/parkinglog", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode == 201 {
		fmt.Println("ParkingLog posted successfully")
	} else {
		return fmt.Errorf("failed to post ParkingLog: %s", resp.Status)
	}

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

	redisClient, err := redis.GetRedisClient(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error connecting to Redis")
	}
	defer redisClient.Client.Close()

	entryEvtProcessor := &EntryEventProcessor{
		RedisClient: redisClient,
	}

	// Handle Entry Events
	if err := rabbitMQClient.ConsumeQueue(cfg.EntryQueueName, entryEvtProcessor); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to consume entry events")
	}

	logger.Log.Debug().Msg("Entry queue consumer set up")

	exitEvtProcessor := &ExitEventProcessor{
		RedisClient: redisClient,
	}

	// Handle Exit Events
	if err := rabbitMQClient.ConsumeQueue(cfg.ExitQueueName, exitEvtProcessor); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to consume exit events")
	}

	logger.Log.Debug().Msg("Exit queue consumer set up")

	select {}
}
