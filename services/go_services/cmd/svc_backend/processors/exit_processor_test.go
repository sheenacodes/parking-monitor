package processors

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"go_services/cmd/svc_backend/metrics"
	"go_services/cmd/svc_backend/models"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"

	"go_services/pkg/logger"

	"github.com/rs/zerolog"
)

func init() {
	// Set up zerolog to log to stderr with human-readable console output
	logger.Log = zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel) // Set the logging level to Debug
}
func TestExitEventProcessor(t *testing.T) {
	testCases := []struct {
		name                 string
		mockDataStore        *MockDataStore
		mockSummaryPoster    *MockSummaryPoster
		msgBody              []byte
		expectedError        bool
		expectedErrorMessage string
		expectedSuccessCount float64
		expectedFailCount    float64
		errorStage           string
	}{
		{
			name: "Success",
			mockDataStore: &MockDataStore{
				AddFieldToHashFunc: func(key, field string, value time.Time) error {
					return nil
				},
				GetFieldAsTimeFunc: func(key, field, layout string) (time.Time, error) {
					return time.Now().Add(-1 * time.Hour), nil
				},
			},
			mockSummaryPoster: &MockSummaryPoster{
				PostSummaryFunc: func(data interface{}) error {
					return nil
				},
			},
			msgBody: func() []byte {
				payload := models.ExitEvent{
					VehiclePlate: "ABC123",
					ExitDateTime: time.Now(),
				}
				data, _ := json.Marshal(payload)
				return data
			}(),
			expectedError:        false,
			expectedErrorMessage: "",
			expectedSuccessCount: 1,
			expectedFailCount:    0,
			errorStage:           "",
		},
		{
			name:                 "FailureOnUnmarshal",
			mockDataStore:        &MockDataStore{},
			mockSummaryPoster:    &MockSummaryPoster{},
			msgBody:              []byte(`{invalid json}`),
			expectedError:        true,
			expectedErrorMessage: "",
			expectedSuccessCount: 0,
			expectedFailCount:    1,
			errorStage:           "json_unmarshal",
		},
		{
			name: "FailureOnDBWrite",
			mockDataStore: &MockDataStore{
				AddFieldToHashFunc: func(key, field string, value time.Time) error {
					return errors.New("DB write error")
				},
				GetFieldAsTimeFunc: func(key, field, layout string) (time.Time, error) {
					return time.Time{}, nil
				},
			},
			mockSummaryPoster: &MockSummaryPoster{},
			msgBody: func() []byte {
				payload := models.ExitEvent{
					VehiclePlate: "ABC123",
					ExitDateTime: time.Now(),
				}
				data, _ := json.Marshal(payload)
				return data
			}(),
			expectedError:        true,
			expectedErrorMessage: "DB write error",
			expectedSuccessCount: 0,
			expectedFailCount:    1,
			errorStage:           "db_write_error",
		},
		{
			name: "FailureOnDBRead",
			mockDataStore: &MockDataStore{
				AddFieldToHashFunc: func(key, field string, value time.Time) error {
					return nil
				},
				GetFieldAsTimeFunc: func(key, field, layout string) (time.Time, error) {
					return time.Time{}, errors.New("entry time retrieval error")
				},
			},
			mockSummaryPoster: &MockSummaryPoster{},
			msgBody: func() []byte {
				payload := models.ExitEvent{
					VehiclePlate: "ABC123",
					ExitDateTime: time.Now(),
				}
				data, _ := json.Marshal(payload)
				return data
			}(),
			expectedError:        true,
			expectedErrorMessage: "entry time retrieval error",
			expectedSuccessCount: 0,
			expectedFailCount:    1,
			errorStage:           "db_read_error",
		},
		{
			name: "FailureOnPostSummary",
			mockDataStore: &MockDataStore{
				AddFieldToHashFunc: func(key, field string, value time.Time) error {
					return nil
				},
				GetFieldAsTimeFunc: func(key, field, layout string) (time.Time, error) {
					return time.Now().Add(-1 * time.Hour), nil
				},
			},
			mockSummaryPoster: &MockSummaryPoster{
				PostSummaryFunc: func(data interface{}) error {
					return errors.New("post summary error")
				},
			},
			msgBody: func() []byte {
				payload := models.ExitEvent{
					VehiclePlate: "ABC123",
					ExitDateTime: time.Now(),
				}
				data, _ := json.Marshal(payload)
				return data
			}(),
			expectedError:        true,
			expectedErrorMessage: "post summary error",
			expectedSuccessCount: 0,
			expectedFailCount:    1,
			errorStage:           "api_post_error",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(tContext *testing.T) {
			metrics.EventProcessingFails.Reset()
			metrics.EventProcessingSuccesses.Reset()
			metrics.EventProcessingLatency.Reset()

			processor := ExitEventProcessor{
				DataStore:     testCase.mockDataStore,
				SummaryPoster: testCase.mockSummaryPoster,
			}

			processError := processor.ProcessMessage(testCase.msgBody)

			if testCase.expectedError {
				assert.Error(tContext, processError)
				count := testutil.ToFloat64(metrics.EventProcessingFails.With(prometheus.Labels{"event_type": "exit", "error_stage": testCase.errorStage}))
				assert.Equal(tContext, testCase.expectedFailCount, count)
				if testCase.expectedErrorMessage != "" {
					assert.Contains(tContext, processError.Error(), testCase.expectedErrorMessage)
				}
			} else {
				assert.NoError(tContext, processError)
				count := testutil.ToFloat64(metrics.EventProcessingSuccesses.With(prometheus.Labels{"event_type": "exit"}))
				assert.Equal(tContext, testCase.expectedSuccessCount, count)

			}

		})
	}
}
