package processors

import (
	"encoding/json"
	"errors"
	"go_services/cmd/svc_backend/metrics"
	"go_services/cmd/svc_backend/models"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

// TestExitEventProcessor_Success tests the successful processing of an exit event message.
func TestExitEventProcessor_Success(t *testing.T) {
	mockDataStore := &MockDataStore{
		AddFieldToHashFunc: func(key, field string, value time.Time) error {
			return nil
		},
		GetFieldAsTimeFunc: func(key, field, layout string) (time.Time, error) {
			return time.Now().Add(-1 * time.Hour), nil // Return a fixed entry time one hour ago
		},
	}

	mockSummaryPoster := &MockSummaryPoster{
		PostSummaryFunc: func(data interface{}) error {
			return nil
		},
	}

	processor := ExitEventProcessor{
		DataStore:     mockDataStore,
		SummaryPoster: mockSummaryPoster,
	}

	// Create a mock exit event payload
	payload := models.ExitEvent{
		VehiclePlate: "ABC123",
		ExitDateTime: time.Now(),
	}
	msgBody, _ := json.Marshal(payload)

	// Call the ProcessMessage method
	err := processor.ProcessMessage(msgBody)

	// Verify no error occurred
	assert.NoError(t, err)

	// Verify metrics are incremented correctly
	assert.Equal(t, 1, testutil.CollectAndCount(metrics.EventProcessingSuccesses))
	assert.Equal(t, 1, testutil.CollectAndCount(metrics.EventProcessingLatency))
}

// TestExitEventProcessor_FailureOnUnmarshal tests the failure scenario when JSON unmarshaling fails.
func TestExitEventProcessor_FailureOnUnmarshal(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockSummaryPoster := &MockSummaryPoster{}

	processor := ExitEventProcessor{
		DataStore:     mockDataStore,
		SummaryPoster: mockSummaryPoster,
	}

	// Create an invalid JSON payload
	msgBody := []byte(`{invalid json}`)

	// Call the ProcessMessage method
	err := processor.ProcessMessage(msgBody)

	// Verify that an error occurred
	assert.Error(t, err)

	// Verify that the correct metric was incremented
	assert.Equal(t, 1, testutil.CollectAndCount(metrics.EventProcessingFails))
}

// TestExitEventProcessor_FailureOnDBWrite tests the failure scenario when writing to the DataStore fails.
func TestExitEventProcessor_FailureOnDBWrite(t *testing.T) {
	mockDataStore := &MockDataStore{
		AddFieldToHashFunc: func(key, field string, value time.Time) error {
			return errors.New("DB write error")
		},
	}
	mockSummaryPoster := &MockSummaryPoster{}

	processor := ExitEventProcessor{
		DataStore:     mockDataStore,
		SummaryPoster: mockSummaryPoster,
	}

	payload := models.ExitEvent{
		VehiclePlate: "ABC123",
		ExitDateTime: time.Now(),
	}
	msgBody, _ := json.Marshal(payload)

	// Call the ProcessMessage method
	err := processor.ProcessMessage(msgBody)

	// Verify that an error occurred
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DB write error")

	// Verify that the correct metric was incremented
	assert.Equal(t, 1, testutil.CollectAndCount(metrics.EventProcessingFails))
}

// TestExitEventProcessor_FailureOnSummaryGeneration tests the failure scenario when summary generation fails.
func TestExitEventProcessor_FailureOnSummaryGeneration(t *testing.T) {
	mockDataStore := &MockDataStore{
		AddFieldToHashFunc: func(key, field string, value time.Time) error {
			return nil
		},
		GetFieldAsTimeFunc: func(key, field, layout string) (time.Time, error) {
			return time.Time{}, errors.New("entry time retrieval error")
		},
	}
	mockSummaryPoster := &MockSummaryPoster{}

	processor := ExitEventProcessor{
		DataStore:     mockDataStore,
		SummaryPoster: mockSummaryPoster,
	}

	payload := models.ExitEvent{
		VehiclePlate: "ABC123",
		ExitDateTime: time.Now(),
	}
	msgBody, _ := json.Marshal(payload)

	// Call the ProcessMessage method
	err := processor.ProcessMessage(msgBody)

	// Verify that an error occurred
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "entry time retrieval error")

	// Verify that the correct metric was incremented
	assert.Equal(t, 1, testutil.CollectAndCount(metrics.EventProcessingFails))
}

// TestExitEventProcessor_FailureOnPostSummary tests the failure scenario when posting the summary fails.
func TestExitEventProcessor_FailureOnPostSummary(t *testing.T) {
	mockDataStore := &MockDataStore{
		AddFieldToHashFunc: func(key, field string, value time.Time) error {
			return nil
		},
		GetFieldAsTimeFunc: func(key, field, layout string) (time.Time, error) {
			return time.Now().Add(-1 * time.Hour), nil
		},
	}

	mockSummaryPoster := &MockSummaryPoster{
		PostSummaryFunc: func(data interface{}) error {
			return errors.New("post summary error")
		},
	}

	processor := ExitEventProcessor{
		DataStore:     mockDataStore,
		SummaryPoster: mockSummaryPoster,
	}

	payload := models.ExitEvent{
		VehiclePlate: "ABC123",
		ExitDateTime: time.Now(),
	}
	msgBody, _ := json.Marshal(payload)

	// Call the ProcessMessage method
	err := processor.ProcessMessage(msgBody)

	// Verify that an error occurred
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "post summary error")

	// Verify that the correct metric was incremented
	assert.Equal(t, 1, testutil.CollectAndCount(metrics.EventProcessingFails))
}
