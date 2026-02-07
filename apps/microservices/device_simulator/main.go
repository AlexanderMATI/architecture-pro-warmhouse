package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/streadway/amqp"
)

type TelemetryData struct {
	SensorID  int       `json:"sensor_id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	log.Println("Device Simulator starting...")
	time.Sleep(15 * time.Second) 

	rabbitURL := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("telemetry", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue: %s", err)
	}

	for {
		temp := 20.0 + rand.Float64()*5.0 
		data := TelemetryData{
			SensorID:  1, 
			Value:     float64(int(temp*10)) / 10,
			Timestamp: time.Now(),
		}
		body, _ := json.Marshal(data)

		err = ch.Publish("", q.Name, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
		if err != nil {
			log.Printf("Failed to publish a message: %s", err)
		} else {
			log.Printf("Sent telemetry: %s", body)
		}

		time.Sleep(5 * time.Second)
	}
}
