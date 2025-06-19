package security

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// AnalysisResponse matches the JSON structure returned by the Python service.
type AnalysisResponse struct {
	Action              string  `json:"action"`
	Confidence          float64 `json:"confidence"`
	Endpoint            string  `json:"endpoint"`
	IPAddress           string  `json:"ip_address"`
	IsAnomaly           bool    `json:"is_anomaly"`
	ReconstructionError float64 `json:"reconstruction_error"`
	ThreatType          string  `json:"threat_type"`
	Timestamp           string  `json:"timestamp"`
}

// AnalyzeRequest sends a log entry to the threat detection service and returns the analysis.
func AnalyzeRequest(logEntryJSON string) (*AnalysisResponse, error) {
	analysisURL := "http://localhost:5000/analyze"
	req, err := http.NewRequest("POST", analysisURL, bytes.NewBuffer([]byte(logEntryJSON)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to analysis service: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("analysis service returned an error (%d): %s", resp.StatusCode, string(body))
	}

	var analysis AnalysisResponse
	if err := json.Unmarshal(body, &analysis); err != nil {
		return nil, fmt.Errorf("failed to decode analysis response: %w", err)
	}
	return &analysis, nil
} 