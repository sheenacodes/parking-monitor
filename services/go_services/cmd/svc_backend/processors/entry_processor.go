package processors

import (
	"encoding/json"
	"go_services/cmd/svc_backend/metrics"
	"go_services/cmd/svc_backend/models"
	"time"

	"go_services/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
)

type DataStore interface {
	AddFieldToHash(hashKey string, fieldName string, fieldValue time.Time) error
}

// EntryEventProcessor handles the processing of entry events.
type EntryEventProcessor struct {
	dataStore DataStore
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

	if err := p.dataStore.AddFieldToHash(hashKey, fieldName, fieldValue); err != nil {
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
