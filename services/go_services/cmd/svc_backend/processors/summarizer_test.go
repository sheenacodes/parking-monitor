package processors

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGenerateParkingSummary_Success tests the successful generation of a parking summary.
func TestGenerateParkingSummary_Success(t *testing.T) {
	vehiclePlate := "ABC123"
	entryDateTime := time.Now().Add(-2 * time.Hour) // Entry time 2 hours ago
	exitDateTime := time.Now()                      // Current time

	parkingLog, err := GenerateParkingSummary(vehiclePlate, exitDateTime, entryDateTime)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the parking log is not nil
	assert.NotNil(t, parkingLog)

	// Assert that the parking log fields are correctly set
	assert.Equal(t, vehiclePlate, parkingLog.VehiclePlate)
	assert.Equal(t, entryDateTime, parkingLog.EntryDateTime)
	assert.Equal(t, exitDateTime, parkingLog.ExitDateTime)
	assert.Equal(t, exitDateTime.Sub(entryDateTime).String(), parkingLog.Duration)
}

// TestGenerateParkingSummary_ExitBeforeEntry tests the scenario where the exit time is before the entry time.
func TestGenerateParkingSummary_ExitBeforeEntry(t *testing.T) {
	vehiclePlate := "XYZ789"
	entryDateTime := time.Now()                    // Current time
	exitDateTime := time.Now().Add(-1 * time.Hour) // Exit time 1 hour ago

	parkingLog, err := GenerateParkingSummary(vehiclePlate, exitDateTime, entryDateTime)

	// Assert that an error occurred
	assert.Error(t, err)

	// Assert that the error message is correct
	assert.Equal(t, fmt.Sprintf("exit time (%v) is before entry time (%v)", exitDateTime, entryDateTime), err.Error())

	// Assert that the parking log is nil
	assert.Nil(t, parkingLog)
}

// TestGenerateParkingSummary_ZeroDuration tests the scenario where the exit time is the same as the entry time.
func TestGenerateParkingSummary_ZeroDuration(t *testing.T) {
	vehiclePlate := "LMN456"
	entryDateTime := time.Now() // Current time
	exitDateTime := entryDateTime

	parkingLog, err := GenerateParkingSummary(vehiclePlate, exitDateTime, entryDateTime)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the parking log is not nil
	assert.NotNil(t, parkingLog)

	// Assert that the parking log fields are correctly set
	assert.Equal(t, vehiclePlate, parkingLog.VehiclePlate)
	assert.Equal(t, entryDateTime, parkingLog.EntryDateTime)
	assert.Equal(t, exitDateTime, parkingLog.ExitDateTime)
	assert.Equal(t, "0s", parkingLog.Duration)
}

// TestGenerateParkingSummary_LongDuration tests a scenario with a very long parking duration.
func TestGenerateParkingSummary_LongDuration(t *testing.T) {
	vehiclePlate := "OPQ123"
	entryDateTime := time.Now().Add(-30 * 24 * time.Hour) // Entry time 30 days ago
	exitDateTime := time.Now()                            // Current time

	parkingLog, err := GenerateParkingSummary(vehiclePlate, exitDateTime, entryDateTime)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the parking log is not nil
	assert.NotNil(t, parkingLog)

	// Assert that the parking log fields are correctly set
	assert.Equal(t, vehiclePlate, parkingLog.VehiclePlate)
	assert.Equal(t, entryDateTime, parkingLog.EntryDateTime)
	assert.Equal(t, exitDateTime, parkingLog.ExitDateTime)
	assert.Equal(t, exitDateTime.Sub(entryDateTime).String(), parkingLog.Duration)
}
