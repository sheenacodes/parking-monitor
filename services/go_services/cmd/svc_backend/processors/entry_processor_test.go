package processors

import (
	"encoding/json"
	"errors"
	"go_services/cmd/svc_backend/metrics"
	"go_services/cmd/svc_backend/models"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestEntryEventProcessor_ProcessMessage(t *testing.T) {
	tests := []struct {
		name          string
		msgBody       []byte
		mockError     error
		expectedError bool
		expectedCount float64
	}{
		{
			name:          "Invalid JSON",
			msgBody:       []byte("{invalid json}"),
			expectedError: true,
			expectedCount: 1, // Expect the JSON unmarshal error metric to increment
		},
		{
			name: "Redis Operation Error",
			msgBody: func() []byte {
				entryDateTime, _ := time.Parse(time.RFC3339, "2024-09-11T10:00:00Z") // Correctly parse the time
				payload := models.EntryEvent{
					VehiclePlate:  "ABC123",
					EntryDateTime: entryDateTime,
				}
				b, _ := json.Marshal(payload)
				return b
			}(),
			mockError:     errors.New("redis write error"),
			expectedError: true,
			expectedCount: 1, // Expect the Redis operation error metric to increment
		},
		{
			name: "Successful Processing",
			msgBody: func() []byte {
				entryDateTime, _ := time.Parse(time.RFC3339, "2024-09-11T10:30:00Z") // Correctly parse the time
				payload := models.EntryEvent{
					VehiclePlate:  "XYZ789",
					EntryDateTime: entryDateTime,
				}
				b, _ := json.Marshal(payload)
				return b
			}(),
			expectedError: false,
			expectedCount: 1, // Expect the success metric to increment
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the metrics
			metrics.EventProcessingFails.Reset()
			metrics.EventProcessingSuccesses.Reset()

			// Create a mock DataStore
			mockDataStore := &MockDataStore{
				AddFieldToHashFunc: func(hashKey string, fieldName string, fieldValue time.Time) error {
					return tt.mockError
				},
			}

			// Create an EntryEventProcessor with the mock DataStore
			processor := &EntryEventProcessor{
				DataStore: mockDataStore,
			}

			// Call ProcessMessage with the test message body
			err := processor.ProcessMessage(tt.msgBody)

			// Verify the expected error state
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify the metrics
			if tt.name == "Invalid JSON" {
				count := testutil.ToFloat64(metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "entry", "error_stage": "unmarshal"}))
				assert.Equal(t, tt.expectedCount, count)
			} else if tt.name == "Redis Operation Error" {
				count := testutil.ToFloat64(metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "entry", "error_stage": "redis_write"}))
				assert.Equal(t, tt.expectedCount, count)
			} else if tt.name == "Successful Processing" {
				count := testutil.ToFloat64(metrics.EventProcessingSuccesses.With(prometheus.Labels{"event_type": "entry"}))
				assert.Equal(t, tt.expectedCount, count)
			}
		})
	}
}
