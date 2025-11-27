package handlers

import (
	"fmt"
	"net/http"
	"smarthome/services"

	"github.com/gin-gonic/gin"
)

// DeviceHandler handles device-related requests
type DeviceHandler struct {
	DeviceService    *services.DeviceService
	TelemetryService *services.TelemetryService
}

// NewDeviceHandler creates a new DeviceHandler
func NewDeviceHandler(DeviceService *services.DeviceService, telemetryService *services.TelemetryService) *DeviceHandler {
	return &DeviceHandler{
		DeviceService:    DeviceService,
		TelemetryService: telemetryService,
	}
}

// RegisterRoutes registers the device routes
func (h *DeviceHandler) RegisterRoutes(router *gin.RouterGroup) {
	devices := router.Group("/devices")
	{
		devices.GET("", h.GetDevicesByHouseId)
		devices.GET("/:id", h.GetDeviceByID)
		devices.POST("", h.CreateDevice)
		devices.PUT("/:id", h.UpdateDevice)
		devices.DELETE("/:id", h.DeleteDevice)
		devices.GET("/:id/telemetry", h.GetTelemetry)
	}
}

// GetDevicesByHouseId handles GET /api/v1/devices
func (h *DeviceHandler) GetDevicesByHouseId(c *gin.Context) {
	houseId := c.Query("houseId")
	if houseId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "houseId is required"})
		return
	}

	// Fetch data from the external API
	resp, err := h.DeviceService.GetDevicesByHouseId(houseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch data: %v", err),
		})
		return
	}

	// Return the data
	c.JSON(http.StatusOK, resp)
}

// GetDeviceByID handles GET /api/v1/devices/:id
func (h *DeviceHandler) GetDeviceByID(c *gin.Context) {
	id := c.Param("id")

	// Fetch data from the external API
	resp, err := h.DeviceService.GetDeviceByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch data: %v", err),
		})
		return
	}

	// Return the data
	c.JSON(http.StatusOK, resp)
}

// GetTelemetry handles GET /api/v1/devices/:id/telemetry
func (h *DeviceHandler) GetTelemetry(c *gin.Context) {
	id := c.Param("id")

	// Fetch telemetry data from the external API
	resp, err := h.TelemetryService.GetTelemetryByDeviceID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch telemetry data: %v", err),
		})
		return
	}

	// Return the telemetry data
	c.JSON(http.StatusOK, resp)
}

// CreateDevice handles POST /api/v1/devices
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var deviceCreate services.DeviceCreate
	if err := c.ShouldBindJSON(&deviceCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch data from the external API
	resp, err := h.DeviceService.CreateDevice(deviceCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to post data: %v", err),
		})
		return
	}

	// Return the data
	c.JSON(http.StatusOK, resp)
}

// UpdateDevice handles PUT /api/v1/devices/:id
func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	id := c.Param("id")

	var deviceUpdate services.DeviceUpdate
	if err := c.ShouldBindJSON(&deviceUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch data from the external API
	resp, err := h.DeviceService.UpdateDevice(id, deviceUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to update data: %v", err),
		})
		return
	}

	// Return the data
	c.JSON(http.StatusOK, resp)
}

// DeleteDevice handles DELETE /api/v1/devices/:id
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id := c.Param("id")

	// Fetch data from the external API
	err := h.DeviceService.DeleteDevice(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to delete data: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}
