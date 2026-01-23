package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Temperature struct {
	Value       float64   `json:"value"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	DeviceID    string    `json:"Device_id"`
	DeviceType  string    `json:"Device_type"`
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
			c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Требуется указание локации"})
			return
		}

		data := generateTemperature(location, "")

		c.JSON(http.StatusOK, data)
	})

	router.GET("/temperature/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Требуется ID устройства"})
			return
		}

		data := generateTemperature("", id)

		c.JSON(http.StatusOK, data)
	})

	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func generateTemperature(location, DeviceID strin {
	rand.Seed(time.Now().UnixNano())
	value := 18.0 + rand.Float64()*10.0

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

	retu{
		Value:       value,
		Timestamp:   time.Now(),
		Location:    location,
		Status:      "активный",
		DeviceID:    DeviceID,
		DeviceType:  "Термодатчик",
		Description: location,
	}
}

