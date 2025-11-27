package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DeviceService handles fetching device data from external API
type DeviceService struct {
	BaseURL    string
	HTTPClient *http.Client
}

type DeviceResponse struct {
	ID            int               `json:"id"`
	SerialNumber  string            `json:"serialNumber"`
	DeviceModelID int               `json:"deviceModelId"`
	HouseID       int               `json:"houseId"`
	Name          string            `json:"name"`
	Status        string            `json:"status"`
	Attributes    []DeviceAttribute `json:"attributes"`
	LastUpdated   time.Time         `json:"lastUpdated"`
	CreatedAt     time.Time         `json:"createdAt"`
}

type DeviceAttribute struct {
	ID       int    `json:"id"`
	DeviceID int    `json:"deviceId"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

type DeviceCreate struct {
	SerialNumber  string `json:"serialNumber" binding:"required"`
	DeviceModelID int    `json:"deviceModelId" binding:"required"`
	HouseID       int    `json:"houseId" binding:"required"`
	Name          string `json:"name" binding:"required"`
}

type DeviceUpdate struct {
	SerialNumber  string `json:"serialNumber" binding:"required"`
	DeviceModelID int    `json:"deviceModelId" binding:"required"`
	HouseID       int    `json:"houseId" binding:"required"`
	Name          string `json:"name" binding:"required"`
}

// NewDeviceService creates a new device service
func NewDeviceService(baseURL string) *DeviceService {
	return &DeviceService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *DeviceService) GetDevicesByHouseId(houseId string) ([]DeviceResponse, error) {
	url := fmt.Sprintf("%s/devices?houseId=%s", s.BaseURL, houseId)

	resp, err := s.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching device data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var deviceResp []DeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceResp); err != nil {
		return nil, fmt.Errorf("error decoding device response: %w", err)
	}

	return deviceResp, nil
}

func (s *DeviceService) GetDeviceByID(deviceID string) (*DeviceResponse, error) {
	url := fmt.Sprintf("%s/devices/%s", s.BaseURL, deviceID)

	resp, err := s.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching device data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var deviceResp DeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceResp); err != nil {
		return nil, fmt.Errorf("error decoding device response: %w", err)
	}

	return &deviceResp, nil
}

func (s *DeviceService) CreateDevice(deviceCreate DeviceCreate) (*DeviceResponse, error) {
	url := fmt.Sprintf("%s/devices", s.BaseURL)

	b, err := json.Marshal(deviceCreate)
	if err != nil {
		return nil, fmt.Errorf("error marshaling device data: %w", err)
	}

	resp, err := s.HTTPClient.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("error post device data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var deviceResp DeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceResp); err != nil {
		return nil, fmt.Errorf("error decoding device response: %w", err)
	}

	return &deviceResp, nil
}

func (s *DeviceService) UpdateDevice(deviceID string, deviceUpdate DeviceUpdate) (*DeviceResponse, error) {
	url := fmt.Sprintf("%s/devices/%s", s.BaseURL, deviceID)

	b, err := json.Marshal(deviceUpdate)
	if err != nil {
		return nil, fmt.Errorf("error marshaling device data: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("error creating PUT request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error putting device data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var deviceResp DeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceResp); err != nil {
		return nil, fmt.Errorf("error decoding device response: %w", err)
	}

	return &deviceResp, nil
}

func (s *DeviceService) DeleteDevice(deviceID string) error {
	url := fmt.Sprintf("%s/devices/%s", s.BaseURL, deviceID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("error creating DELETE request: %w", err)
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error delete device: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
