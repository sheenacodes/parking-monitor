package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPClientPoster implements the SummaryPoster interface using http.Client.
type HTTPClientPoster struct {
	Client *http.Client
	APIURL string
}

// PostSummary sends any struct to the REST API server.
func (p *HTTPClientPoster) PostSummary(data interface{}) error {
	// Marshal the data into JSON
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", p.APIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to post data: %s", resp.Status)
	}

	return nil
}
