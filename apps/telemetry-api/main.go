package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Telemetry struct {
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

func main() {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Get telemetry by device ID
	router.GET("/telemetry/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "device ID is required"})
			return
		}

		// Generate random telemetry data based on device ID
		data := generateTelemetryData(id)

		c.JSON(http.StatusOK, data)
	})

	// Start server
	log.Println("Telemetry API starting on :8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func generateTelemetryData(deviceID string) []Telemetry {

	telemetry := make([]Telemetry, 0, 2)

	for i := 0; i < 2; i++ {
		metrics := make([]TelemetryMetric, 0, 2)

		temperatureValue := 18.0 + float64(time.Now().UnixNano()%10) + float64(time.Now().UnixNano()%100)/100.0
		metric1 := TelemetryMetric{
			ID:    "1",
			Key:   "temperature",
			Value: strconv.FormatFloat(temperatureValue, 'f', 2, 64),
			Unit:  "С",
		}

		humidityValue := 20.0 + float64(time.Now().UnixNano()%20) + float64(time.Now().UnixNano()%100)/100.0
		metric2 := TelemetryMetric{
			ID:    "2",
			Key:   "humidity",
			Value: strconv.FormatFloat(humidityValue, 'f', 2, 64),
			Unit:  "%",
		}

		metrics = append(metrics, metric1, metric2)

		t := Telemetry{
			ID:        "1",
			DeviceID:  deviceID,
			CreatedAt: time.Now(),
			Metric:    metrics,
			HouseID:   "2",
		}

		telemetry = append(telemetry, t)
	}

	return telemetry
}
