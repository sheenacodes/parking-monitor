package rabbitmq

import (
	"encoding/json"
	"generator/logger"
	"time"

	"github.com/rabbitmq/amqp091-go"
	//"github.com/rs/zerolog/log"
)

const (
	maxRetries     = 5                // Maximum number of retries before giving up
	initialBackoff = 2 * time.Second  // Initial delay before retrying
	maxBackoff     = 30 * time.Second // Maximum delay between retries
)

// ConnectToRabbitMQ connects to RabbitMQ with retry logic
func ConnectToRabbitMQ(url string) *amqp091.Connection {
	var conn *amqp091.Connection
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
		logger.Log.Fatal().Err(err).Msg("Failed to connect to RabbitMQ after multiple attempts")
	}

	return conn
}

func PublishEvent(conn *amqp091.Connection, queueName string, eventPayload any) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	body, err := json.Marshal(eventPayload)
	if err != nil {
		logger.Log.Fatal().Err(err).Msgf("JSON conversion error in %v...", eventPayload)
		return err
	}

	err = channel.Publish(
		"",        // exchange
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	logger.Log.Info().Msgf("Published event: %s", body)
	return nil
}
