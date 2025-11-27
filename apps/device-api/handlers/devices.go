package handlers

import (
	"context"
	"net/http"
	"smarthome/db"
	"smarthome/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeviceHandler handles device-related requests
type DeviceHandler struct {
	DB *db.DB
}

// NewDeviceHandler creates a new DeviceHandler
func NewDeviceHandler(db *db.DB) *DeviceHandler {
	return &DeviceHandler{
		DB: db,
	}
}

// RegisterRoutes registers the device routes
func (h *DeviceHandler) RegisterRoutes(router *gin.RouterGroup) {
	devices := router.Group("/devices")
	{
		devices.GET("", h.GetDevices)
		devices.GET("/:id", h.GetDeviceByID)
		devices.POST("", h.CreateDevice)
		devices.PUT("/:id", h.UpdateDevice)
		devices.DELETE("/:id", h.DeleteDevice)
	}
}

// GetDevices handles GET /api/v1/devices
func (h *DeviceHandler) GetDevices(c *gin.Context) {
	houseIdStr := c.Query("houseId")
	if houseIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "houseId is required"})
		return
	}
	houseId, err := strconv.Atoi(houseIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid houseId"})
		return
	}
	devices, err := h.DB.GetDevicesByHouseId(context.Background(), int(houseId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, device := range devices {
		attrs, err := h.DB.GetDevicesAttributes(context.Background(), device.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Device attributes not found"})
		}
		device.Attributes = attrs
	}

	c.JSON(http.StatusOK, devices)
}

// GetDeviceByID handles GET /api/v1/devices/:id
func (h *DeviceHandler) GetDeviceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	device, err := h.DB.GetDeviceByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	attrs, err := h.DB.GetDevicesAttributes(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device attributes not found"})
		return
	}

	device.Attributes = attrs

	c.JSON(http.StatusOK, device)
}

// CreateDevice handles POST /api/v1/devices
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var deviceCreate models.DeviceCreate
	if err := c.ShouldBindJSON(&deviceCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device, err := h.DB.CreateDevice(context.Background(), deviceCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = h.initDeviceAttributes(device.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, device)
}

// UpdateDevice handles PUT /api/v1/devices/:id
func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	var deviceUpdate models.DeviceUpdate
	if err := c.ShouldBindJSON(&deviceUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device, err := h.DB.UpdateDevice(context.Background(), id, deviceUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = h.DB.DeleteDeviceAttributes(context.Background(), device.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if err = h.initDeviceAttributes(device.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, device)
}

// DeleteDevice handles DELETE /api/v1/devices/:id
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	if err = h.DB.DeleteDeviceAttributes(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	err = h.DB.DeleteDevice(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}

func (h *DeviceHandler) initDeviceAttributes(deviceId int) error {
	attrs := make([]models.DeviceAttribute, 0, 2)
	attr1 := models.DeviceAttribute{
		DeviceID: deviceId,
		Key:      "temperature",
		Value:    "18.0",
	}
	attrs = append(attrs, attr1)
	attr2 := models.DeviceAttribute{
		DeviceID: deviceId,
		Key:      "humidity",
		Value:    "30.0",
	}
	attrs = append(attrs, attr2)

	for _, attr := range attrs {
		_, err := h.DB.CreateDeviceAttribute(context.Background(), attr)
		if err != nil {
			return err
		}
	}

	return nil
}
