package handlers

import (
	"encoding/json"
	"log"

	"github.com/mikebrgs/home-matrix/internal/database"
	"github.com/mikebrgs/home-matrix/internal/models/pot"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewMQTTHandler(db database.TimescaleDB) *MQTTHandler {
	return &MQTTHandler{&db}
}

type MQTTHandler struct {
	db *database.TimescaleDB
}

func (handler *MQTTHandler) HandlePotHealthMessage(client mqtt.Client, msg mqtt.Message) {
	var data pot.PotHealth
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Printf("Failed to unmarshal MQTT message: %v", err)
		return
	}

	log.Printf("Received data: %+v", data)

	if err := handler.db.InsertPotHealthData(data); err != nil {
		log.Printf("Failed to insert data into PostgreSQL: %v", err)
	}
}

func (handler *MQTTHandler) HandlePotStatusMessage(client mqtt.Client, msg mqtt.Message) {
	var data pot.PotStatus
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Printf("Failed to unmarshal MQTT message: %v", err)
		return
	}

	log.Printf("Received data: %+v", data)

	if err := handler.db.InsertPotStatusData(data); err != nil {
		log.Printf("Failed to insert data into PostgreSQL: %v", err)
	}
}
