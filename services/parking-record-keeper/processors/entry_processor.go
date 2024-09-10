package processors

import (
	"encoding/json"
	"parking-record-keeper/metrics"
	"parking-record-keeper/models"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sheenacodes/sharedutils/logger"
	"github.com/sheenacodes/sharedutils/redis"
)

// EntryEventProcessor handles the processing of entry events.
type EntryEventProcessor struct {
	RedisClient *redis.RedisClient
}

// ProcessMessage processes an entry event message.
func (p *EntryEventProcessor) ProcessMessage(msgBody []byte) error {
	start := time.Now() // metrics instrumentation: start timer

	var payload models.EntryEvent
	if err := json.Unmarshal(msgBody, &payload); err != nil {
		// metrics instrumentation: Increment the error counter for JSON unmarshal error
		metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "entry", "error_stage": "unmarshal"}).Inc()

		return err
	}

	hashKey := payload.VehiclePlate
	fieldName := "entry_date_time"
	fieldValue := payload.EntryDateTime
	logger.Log.Debug().Msgf("Storing entry: key - %s; field - %s; value - %s", hashKey, fieldName, fieldValue)

	if err := p.RedisClient.AddFieldToHash(hashKey, fieldName, fieldValue); err != nil {
		// metrics instrumentation: Increment the error counter for Redis operation error
		metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "entry", "error_stage": "redis_write"}).Inc()

		return err
	}

	// metrics instrumentation: Record the duration taken to process the message
	duration := time.Since(start).Seconds()
	metrics.EventProcessingLatency.With(prometheus.Labels{"event_type": "entry"}).Observe(duration)

	// metrics instrumentation:
	metrics.EventProcessingSuccesses.With(prometheus.Labels{"event_type": "entry"}).Inc()

	return nil
}
