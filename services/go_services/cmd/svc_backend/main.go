package main

import (
	"go_services/cmd/svc_backend/config"
	"go_services/cmd/svc_backend/processors"
	"go_services/pkg/logger"
	"go_services/pkg/rabbitmq"
	"go_services/pkg/redis"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	// Load configuration
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.LogLevel)

	// Initialize RabbitMQ client
	rabbitMQClient, err := rabbitmq.GetRabbitMQClient(cfg.RabbitMQURL)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to initialize RabbitMQ client")
	}
	defer rabbitMQClient.Close()

	// Initialize Redis client
	redisClient, err := redis.GetRedisClient(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error connecting to Redis")
	}
	defer redisClient.Client.Close()

	// Start the Prometheus metrics HTTP server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Log.Debug().Msg("Starting Prometheus metrics server on :2112/metrics")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			logger.Log.Fatal().Err(err).Msgf("Error starting Prometheus server: %v", err)
		}
	}()

	// Initialize EntryEventProcessor
	entryEvtProcessor := &processors.EntryEventProcessor{
		RedisClient: redisClient,
	}
	// Handle Entry Events
	if err := rabbitMQClient.ConsumeQueue(cfg.EntryQueueName, entryEvtProcessor); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to consume entry events")
	}
	logger.Log.Debug().Msg("Entry queue consumer set up")

	// Initialize ExitEventProcessor
	exitEvtProcessor := &processors.ExitEventProcessor{
		RedisClient: redisClient,
		APIURL:      cfg.APIURL, // Set the API URL from configuration
	}
	// Handle Exit Events
	if err := rabbitMQClient.ConsumeQueue(cfg.ExitQueueName, exitEvtProcessor); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to consume exit events")
	}
	logger.Log.Debug().Msg("Exit queue consumer set up")

	// Keep the main function running
	select {}
}
