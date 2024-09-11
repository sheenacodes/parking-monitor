package processors

import (
	"fmt"
	"go_services/cmd/svc_backend/models"
	"time"
)

// GenerateParkingSummary creates a ParkingLog based on the exit event.
func GenerateParkingSummary(vehiclePlate string, exitDateTime time.Time, entryDateTime time.Time) (*models.ParkingLog, error) {

	// Check if exit time is before entry time
	if exitDateTime.Before(entryDateTime) {
		return nil, fmt.Errorf("exit time (%v) is before entry time (%v)", exitDateTime, entryDateTime)
	}

	parkingDuration := exitDateTime.Sub(entryDateTime).String()

	return &models.ParkingLog{
		VehiclePlate:  vehiclePlate,
		EntryDateTime: entryDateTime,
		ExitDateTime:  exitDateTime,
		Duration:      parkingDuration,
	}, nil
}
