package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Telemetry struct {
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

var telemetryData []Telemetry

func main() {
	router := gin.Default()

	router.POST("/telemetry", func(c *gin.Context) {
		var data map[string]interface{}
		
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		telemetry := Telemetry{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Data:      data,
		}

		telemetryData = append(telemetryData, telemetry)

		c.JSON(http.StatusCreated, gin.H{
			"message": "Telemetry received",
			"data":    telemetry,
		})
	})

	router.GET("/telemetry", func(c *gin.Context) {
		c.JSON(http.StatusOK, telemetryData)
	})

	log.Println("Telemetry API starting on :5001")
	if err := router.Run(":5001"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
}