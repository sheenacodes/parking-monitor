package summary_logger

import (
	"fmt"
	"parking-record-keeper/models"
	"time"

	"github.com/sheenacodes/sharedutils/redis"
)

// GenerateParkingSummary creates a ParkingLog based on the exit event.
func GenerateParkingSummary(vehiclePlate string, exitDateTime time.Time, rClient *redis.RedisClient) (*models.ParkingLog, error) {
	fieldName := "entry_date_time"
	layout := time.RFC3339
	entryDateTime, err := rClient.GetFieldAsTime(vehiclePlate, fieldName, layout)
	if err != nil {
		return nil, fmt.Errorf("error retrieving entry time: %v", err)
	}

	parkingDuration := exitDateTime.Sub(entryDateTime).String()
	return &models.ParkingLog{
		VehiclePlate:  vehiclePlate,
		EntryDateTime: entryDateTime,
		ExitDateTime:  exitDateTime,
		Duration:      parkingDuration,
	}, nil
}
