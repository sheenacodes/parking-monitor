package summary_logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"parking-record-keeper/models"
)

// PostSummary sends the ParkingLog to the REST API server.
func PostSummary(apiURL string, parkingLog models.ParkingLog) error {
	jsonBody, err := json.Marshal(parkingLog)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to post ParkingLog: %s", resp.Status)
	}

	return nil
}
