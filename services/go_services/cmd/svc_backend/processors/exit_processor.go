package processors

import (
	"encoding/json"
	"go_services/cmd/svc_backend/metrics"
	"go_services/cmd/svc_backend/models"
	"go_services/cmd/svc_backend/summary_logger"
	"go_services/pkg/logger"
	"go_services/pkg/redis"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ExitEventProcessor handles the processing of exit events.
type ExitEventProcessor struct {
	RedisClient *redis.RedisClient
	APIURL      string
}

// ProcessMessage processes an exit event message.
func (p *ExitEventProcessor) ProcessMessage(msgBody []byte) error {
	start := time.Now() // metrics instrumentation: Start time for latency measurement

	// Unmarshal the incoming JSON message to the ExitEvent struct
	var payload models.ExitEvent
	if err := json.Unmarshal(msgBody, &payload); err != nil {
		metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "exit", "error_stage": "unmarshal"}).Inc()
		return err
	}

	hashKey := payload.VehiclePlate
	fieldName := "exit_date_time"
	fieldValue := payload.ExitDateTime
	logger.Log.Debug().Msgf("Storing exit: key - %s; field - %s; value - %s", hashKey, fieldName, fieldValue)

	// Store the exit time in Redis
	if err := p.RedisClient.AddFieldToHash(hashKey, fieldName, fieldValue); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed writing to Redis")
		// metrics instrumentation:
		metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "exit", "error_stage": "redis_write"}).Inc()
		return err
	}

	// Generate the parking summary
	parkingLog, err := summary_logger.GenerateParkingSummary(payload.VehiclePlate, payload.ExitDateTime, p.RedisClient)
	if err != nil {
		// metrics instrumentation:
		metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "exit", "error_stage": "generate_summary"}).Inc()
		return err
	}

	// Post the parking summary to the API
	if err := summary_logger.PostSummary(p.APIURL, *parkingLog); err != nil {
		// metrics instrumentation:
		metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "exit", "error_stage": "post_summary"}).Inc()
		return err
	}

	// metrics instrumentation: Record the duration taken to process the message
	duration := time.Since(start).Seconds()
	metrics.EventProcessingLatency.With(prometheus.Labels{"event_type": "exit"}).Observe(duration)

	// metrics instrumentation:
	metrics.EventProcessingSuccesses.With(prometheus.Labels{"event_type": "exit"}).Inc()

	return nil
}
