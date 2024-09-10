package models

import "time"

// EntryEvent represents an event payload when a vehicle enters the parking area.
type EntryEvent struct {
	ID            string    `json:"id"`
	VehiclePlate  string    `json:"vehicle_plate"`
	EntryDateTime time.Time `json:"entry_date_time"`
}

// ExitEvent represents an event payload when a vehicle exits the parking area.
type ExitEvent struct {
	ID           string    `json:"id"`
	VehiclePlate string    `json:"vehicle_plate"`
	ExitDateTime time.Time `json:"exit_date_time"`
}

// ParkingLog represents the log of parking duration to be used as postbody in api calls.
type ParkingLog struct {
	VehiclePlate  string    `json:"vehicle_plate"`
	ExitDateTime  time.Time `json:"exit_date_time"`
	EntryDateTime time.Time `json:"entry_date_time"`
	Duration      string    `json:"duration"`
}
