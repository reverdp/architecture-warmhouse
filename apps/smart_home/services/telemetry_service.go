package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TelemetryService handles fetching telemetry data from external API
type TelemetryService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// TelemetryResponse represents the response from the telemetry API
type TelemetryResponse struct {
	ID        string            `json:"id"`
	DeviceID  string            `json:"deviceId"`
	CreatedAt time.Time         `json:"createdAt"`
	Metric    []TelemetryMetric `json:"metricId"`
	HouseID   string            `json:"houseId"`
}

type TelemetryMetric struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

// NewTelemetryService creates a new telemetry service
func NewTelemetryService(baseURL string) *TelemetryService {
	return &TelemetryService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetTelemetryByDeviceID fetches telemetry data for a specific Device ID
func (s *TelemetryService) GetTelemetryByDeviceID(deviceID string) ([]TelemetryResponse, error) {
	url := fmt.Sprintf("%s/telemetry/%s", s.BaseURL, deviceID)

	resp, err := s.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching telemetry data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var telemetryResp []TelemetryResponse
	if err := json.NewDecoder(resp.Body).Decode(&telemetryResp); err != nil {
		return nil, fmt.Errorf("error decoding telemetry response: %w", err)
	}

	return telemetryResp, nil
}
