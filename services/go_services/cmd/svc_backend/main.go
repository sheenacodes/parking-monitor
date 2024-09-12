package main

import (
	"go_services/cmd/svc_backend/config"
	"go_services/cmd/svc_backend/processors"
	"go_services/pkg/logger"
	"go_services/pkg/rabbitmq"
	"go_services/pkg/redis"
	"go_services/pkg/restapi"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// print config for debug purposes
func printConfig(cfg *config.Config) {
	logger.Log.Info().
		Str("RabbitMQURL", cfg.RabbitMQURL).
		Str("RedisAddress", cfg.RedisAddress).
		Str("RedisPassword", cfg.RedisPassword).
		Int("RedisDB", cfg.RedisDB).
		Str("EntryQueueName", cfg.EntryQueueName).
		Str("ExitQueueName", cfg.ExitQueueName).
		Str("APIURL", cfg.APIURL).
		Msg("Configuration settings")
}

// loadConfig loads the application configuration
func loadConfig() *config.Config {
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.LogLevel)
	printConfig(cfg)
	return cfg
}

// initializeServices initializes the RabbitMQ and Redis clients
func initializeServices(cfg *config.Config) (*rabbitmq.RabbitMQClient, *redis.RedisClient, error) {
	rabbitMQClient, err := rabbitmq.GetRabbitMQClient(cfg.RabbitMQURL)
	if err != nil {
		return nil, nil, err
	}

	redisClient, err := redis.GetRedisClient(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		rabbitMQClient.Close()
		return nil, nil, err
	}

	return rabbitMQClient, redisClient, nil
}

// startMetricsServer starts the Prometheus metrics HTTP server
func startMetricsServer() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Log.Debug().Msg("Starting Prometheus metrics server on :2112/metrics")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			logger.Log.Fatal().Err(err).Msgf("Error starting Prometheus server: %v", err)
		}
	}()
}

// setupEventProcessors sets up the entry and exit event processors
func setupEventProcessors(cfg *config.Config, rabbitMQClient *rabbitmq.RabbitMQClient, redisClient *redis.RedisClient) error {
	// Initialize EntryEventProcessor
	entryEvtProcessor := &processors.EntryEventProcessor{
		DataStore: redisClient,
	}

	// Handle Entry Events
	if err := rabbitMQClient.ConsumeQueue(cfg.EntryQueueName, entryEvtProcessor); err != nil {
		return err
	}
	logger.Log.Debug().Msg("Entry queue consumer set up")

	// Configure and create SummaryPoster implementation
	summaryPoster := &restapi.HTTPClientPoster{
		Client: &http.Client{},
		APIURL: cfg.APIURL,
	}

	// Initialize ExitEventProcessor
	exitEvtProcessor := &processors.ExitEventProcessor{
		DataStore:     redisClient,
		SummaryPoster: summaryPoster,
	}

	// Handle Exit Events
	if err := rabbitMQClient.ConsumeQueue(cfg.ExitQueueName, exitEvtProcessor); err != nil {
		return err
	}
	logger.Log.Debug().Msg("Exit queue consumer set up")

	return nil
}

func main() {
	// Load configuration
	cfg := loadConfig()

	// Initialize services
	rabbitMQClient, redisClient, err := initializeServices(cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to initialize services")
	}
	defer rabbitMQClient.Close()
	defer redisClient.Client.Close()

	// Start the Prometheus metrics server
	startMetricsServer()

	// Set up event processors
	if err := setupEventProcessors(cfg, rabbitMQClient, redisClient); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to set up event processors")
	}

	// Keep the main function running
	select {}
}
