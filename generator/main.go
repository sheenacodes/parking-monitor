package main

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitializeLogger sets up the zerolog logger based on environment variables
func InitializeLogger() {
	// Set the time format for zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Read log level from environment variable
	logLevel := os.Getenv("LOG_LEVEL")
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "WARN", "WARNING":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel) // Default to INFO
	}

	// Set up logging to console with pretty output if needed
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(output)
}

func main() {
	// Initialize logger
	InitializeLogger()

	// Get RabbitMQ URL from environment variable
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		log.Fatal().Msg("RABBITMQ_URL environment variable is not set")
	}

	log.Info().Str("rabbitmq_url", rabbitMQURL).Msg("Connecting to RabbitMQ")

	// Example usage of logs
	log.Debug().Msg("This is a debug message")
	log.Info().Msg("This is an info message")
	log.Warn().Msg("This is a warning message")
	log.Error().Msg("This is an error message")

	// Structured logging example
	log.Info().
		Str("component", "main").
		Str("status", "starting").
		Msg("Service is starting")

	// Simulate error handling
	err := simulateError()
	if err != nil {
		log.Error().Err(err).Msg("An error occurred")
	}
}

// Simulate an error function
func simulateError() error {
	// Simulate an error for demonstration purposes
	return os.ErrInvalid
}
