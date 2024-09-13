package processors

import (
	"encoding/json"
	"go_services/cmd/svc_backend/metrics"
	"go_services/cmd/svc_backend/models"
	"time"

	"go_services/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
)

// EntryEventProcessor handles the processing of entry events.
type EntryEventProcessor struct {
	DataStore DataStore
}

// ProcessMessage processes an entry event message.
func (p *EntryEventProcessor) ProcessMessage(msgBody []byte) error {
	start := time.Now() // metrics instrumentation: start timer
	logger.Log.Info().Msg("Process Entry Event")
	var payload models.EntryEvent
	if err := json.Unmarshal(msgBody, &payload); err != nil {
		// metrics instrumentation: Increment the error counter for JSON unmarshal error
		metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "entry", "error_stage": "json_unmarshal"}).Inc()

		return err
	}

	hashKey := payload.VehiclePlate
	fieldName := "entry_date_time"
	fieldValue := payload.EntryDateTime
	logger.Log.Debug().Msgf("Storing entry: key - %s; field - %s; value - %s", hashKey, fieldName, fieldValue)

	if err := p.DataStore.AddFieldToHash(hashKey, fieldName, fieldValue); err != nil {
		// metrics instrumentation: Increment the error counter for Redis operation error
		metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "entry", "error_stage": "db_write_error"}).Inc()

		return err
	}

	logger.Log.Info().Msg("Process Entry Event Success")
	// metrics instrumentation: Record the duration taken to process the message
	duration := time.Since(start).Seconds()
	metrics.EventProcessingLatency.With(prometheus.Labels{"event_type": "entry"}).Observe(duration)

	// metrics instrumentation:
	metrics.EventProcessingSuccesses.With(prometheus.Labels{"event_type": "entry"}).Inc()

	return nil
}
