package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mikebrgs/home-matrix/database"
	"github.com/mikebrgs/home-matrix/models/pot"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var db database.Database

// SetDatabase sets the database for the handler
func SetDatabase(database database.Database) {
	db = database
}

func HandlePotHealthMessage(client mqtt.Client, msg mqtt.Message) {
	var data pot.PotHealth
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Printf("Failed to unmarshal MQTT message: %v", err)
		return
	}

	log.Printf("Received data: %+v", data)

	if err := db.InsertPotHealthData(context.Background(), data); err != nil {
		log.Printf("Failed to insert data into PostgreSQL: %v", err)
	}
}

func HandlePotStatusMessage(client mqtt.Client, msg mqtt.Message) {
	var data pot.PotStatus
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Printf("Failed to unmarshal MQTT message: %v", err)
		return
	}

	log.Printf("Received data: %+v", data)

	if err := db.InsertPotStatusData(context.Background(), data); err != nil {
		log.Printf("Failed to insert data into PostgreSQL: %v", err)
	}
}
