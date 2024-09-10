package rabbitmq

import (
	"time"

	"go_services/pkg/logger"

	"github.com/rabbitmq/amqp091-go"
)

const (
	maxRetries     = 7                 // Maximum number of retries before giving up
	initialBackoff = 3 * time.Second   // Initial delay before retrying
	maxBackoff     = 200 * time.Second // Maximum delay between retries
)

// RabbitMQClient holds the RabbitMQ connection and channel
type RabbitMQClient struct {
	Connection *amqp091.Connection
	Channel    *amqp091.Channel
}

// NewRabbitMQClient creates a new RabbitMQ client with retry logic
func GetRabbitMQClient(url string) (*RabbitMQClient, error) {
	var conn *amqp091.Connection
	var ch *amqp091.Channel
	var err error

	for retries := 0; retries < maxRetries; retries++ {
		conn, err = amqp091.Dial(url)
		if err != nil {
			backoff := time.Duration((1 << retries) * int(initialBackoff))
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			logger.Log.Warn().Err(err).Msgf("Failed to connect to RabbitMQ, retrying in %v...", backoff)
			time.Sleep(backoff)
		} else {
			logger.Log.Info().Msg("Successfully connected to RabbitMQ")
			break
		}
	}

	if err != nil {
		return nil, err
	}

	// Create a new channel
	ch, err = conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQClient{
		Connection: conn,
		Channel:    ch,
	}, nil
}

// Close cleans up the RabbitMQ connection and channel
func (r *RabbitMQClient) Close() {
	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to close RabbitMQ channel")
		}
	}
	if r.Connection != nil {
		if err := r.Connection.Close(); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to close RabbitMQ connection")
		}
	}
}
