package event

import (
	"fmt"
	"math/rand"
	"time"
)

type EntryEventPayload struct {
	ID            string    `json:"id"`
	VehiclePlate  string    `json:"vehicle_plate"`
	EntryDateTime time.Time `json:"entry_date_time"`
}

func GenerateEntryEvent() EntryEventPayload {

	return EntryEventPayload{
		ID:            fmt.Sprintf("%d", rand.Int()),
		VehiclePlate:  fmt.Sprintf("plate-%d", rand.Intn(1000)),
		EntryDateTime: time.Now().UTC(),
	}
}

type ExitEventPayload struct {
	ID           string    `json:"id"`
	VehiclePlate string    `json:"vehicle_plate"`
	ExitDateTime time.Time `json:"exit_date_time"`
}

func GenerateExitEvent() ExitEventPayload {

	return ExitEventPayload{
		ID:           fmt.Sprintf("%d", rand.Int()),
		VehiclePlate: fmt.Sprintf("plate-%d", rand.Intn(1000)),
		ExitDateTime: time.Now().UTC(),
	}
}
