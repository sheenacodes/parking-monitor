package rabbitmq

import (
	"encoding/json"

	"go_services/pkg/logger"

	"github.com/rabbitmq/amqp091-go"
)

func PublishEvent(conn *amqp091.Connection, queueName string, eventPayload any) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	body, err := json.Marshal(eventPayload)
	if err != nil {
		logger.Log.Error().Err(err).Msgf("JSON conversion error in %v...", eventPayload)
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
