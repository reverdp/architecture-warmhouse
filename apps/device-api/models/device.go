package models

import "time"

// DeviceStatus represents the status of device
type DeviceStatus string

const (
	DeviceStatusActive   DeviceStatus = "active"
	DeviceStatusInActive DeviceStatus = "inactive"
)

// Device represents a smart home device
type Device struct {
	ID            int               `json:"id"`
	SerialNumber  string            `json:"serialNumber"`
	DeviceModelID int               `json:"deviceModelId"`
	HouseID       int               `json:"houseId"`
	Name          string            `json:"name"`
	Status        DeviceStatus      `json:"status"`
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

// DeviceCreate represents the data needed to create a new device
type DeviceCreate struct {
	SerialNumber  string `json:"serialNumber" binding:"required"`
	DeviceModelID int    `json:"deviceModelId" binding:"required"`
	HouseID       int    `json:"houseId" binding:"required"`
	Name          string `json:"name" binding:"required"`
}

// DeviceUpdate represents the data that can be updated for a device
type DeviceUpdate struct {
	SerialNumber  string `json:"serialNumber" binding:"required"`
	DeviceModelID int    `json:"deviceModelId" binding:"required"`
	HouseID       int    `json:"houseId" binding:"required"`
	Name          string `json:"name" binding:"required"`
}
