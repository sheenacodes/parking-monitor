package processors

import (
	"encoding/json"
	"fmt"
	"go_services/cmd/svc_backend/metrics"
	"go_services/cmd/svc_backend/models"
	"go_services/pkg/logger"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ExitEventProcessor handles the processing of exit events.
type ExitEventProcessor struct {
	DataStore     DataStore
	SummaryPoster SummaryPoster
}

// ProcessMessage processes an exit event message.
func (p *ExitEventProcessor) ProcessMessage(msgBody []byte) error {
	start := time.Now() // metrics instrumentation: Start time for latency measurement

	// Unmarshal the incoming JSON message to the ExitEvent struct
	var payload models.ExitEvent
	if err := json.Unmarshal(msgBody, &payload); err != nil {
		metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "exit", "error_stage": "json_unmarshal"}).Inc()
		return err
	}

	hashKey := payload.VehiclePlate
	fieldName := "exit_date_time"
	fieldValue := payload.ExitDateTime
	logger.Log.Debug().Msgf("Storing exit: key - %s; field - %s; value - %s", hashKey, fieldName, fieldValue)

	// Store the exit time
	if err := p.DataStore.AddFieldToHash(hashKey, fieldName, fieldValue); err != nil {
		logger.Log.Error().Err(err).Msg("Failed writing to datastore")
		// metrics instrumentation:
		metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "exit", "error_stage": "db_write_error"}).Inc()
		return err
	}

	fieldName = "entry_date_time"
	layout := time.RFC3339
	entryDateTime, err := p.DataStore.GetFieldAsTime(payload.VehiclePlate, fieldName, layout)
	if err != nil {
		// metrics instrumentation:
		metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "exit", "error_stage": "db_read_error"}).Inc()
		return fmt.Errorf("error retrieving entry time: %v", err)
	}

	// Generate the parking summary
	parkingLog, err := GenerateParkingSummary(payload.VehiclePlate, payload.ExitDateTime, entryDateTime)
	if err != nil {
		// metrics instrumentation:
		metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "exit", "error_stage": "generate_summary"}).Inc()
		return err
	}

	// Post the parking summary to the API
	if err := p.SummaryPoster.PostSummary(*parkingLog); err != nil {
		// metrics instrumentation:
		metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "exit", "error_stage": "api_post_error"}).Inc()
		return err
	}

	// metrics instrumentation: Record the duration taken to process the message
	duration := time.Since(start).Seconds()
	metrics.EventProcessingLatency.With(prometheus.Labels{"event_type": "exit"}).Observe(duration)

	// metrics instrumentation:
	metrics.EventProcessingSuccesses.With(prometheus.Labels{"event_type": "exit"}).Inc()

	return nil
}
