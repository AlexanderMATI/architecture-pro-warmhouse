package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TempData struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func main() {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.GET("/temperature", func(c *gin.Context) {
		location := c.Query("location")
		if location == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Требуется указание местоположения"})
			return
		}

		data := generateTemperatureData(location, "")

		c.JSON(http.StatusOK, data)
	})

	router.GET("/temperature/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Требуется указать id сенсора"})
			return
		}

		data := generateTemperatureData("", id)

		c.JSON(http.StatusOK, data)
	})

	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func generateTemperatureData(location, sensorID string) TempData {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	value := 18.0 + r.Float64()*20.0

	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}
	
	return TempData{
		Value:       value,
		Unit:        "°C",
		Timestamp:   time.Now(),
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "Термодатчик",
		Description: location,
	}
}