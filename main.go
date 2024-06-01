package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"

	"github.com/mikebrgs/home-matrix/internal/config"
	"github.com/mikebrgs/home-matrix/internal/database"
	"github.com/mikebrgs/home-matrix/internal/handlers"
	"github.com/mikebrgs/home-matrix/internal/mqtt"
)

func main() {
	// Load config file
	cfg, err := config.LoadConfig("configs/home-matrix.yaml")
	if err != nil {
		slog.Error("Unable to open config: %v\n", err)
		os.Exit(1)
	}

	// Connect to DB
	var connStr = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Dbname)
	slog.Info("Connecting to", "db", connStr)
	db_client, err := database.NewTimescaleDBClient(connStr)
	if err != nil {
		slog.Error("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	db_client.Connect()
	defer db_client.Close()

	// Connect to MQTT broker
	mqtt_client, err := mqtt.NewMQTTClient(
		fmt.Sprintf("mqtt://%s:%d", cfg.MQTT.Host, cfg.MQTT.Port),
		fmt.Sprintf("home-matrix-%d", rand.Intn(1000)),
	)
	if err != nil {
		slog.Error("Unable to connect to MQTT broker: %v\n", err)
		os.Exit(1)
	}
	mqtt_client.Connect()
	defer mqtt_client.Disconnect()

	// Create the handler
	handler := handlers.NewMQTTHandler(*db_client)
	if err = mqtt_client.Subscribe("pot/health", 0, handler.HandlePotHealthMessage); err != nil {
		slog.Error("Unable to subscribe: %v\n", err)
		os.Exit(1)
	}
	if err = mqtt_client.Subscribe("pot/status", 0, handler.HandlePotStatusMessage); err != nil {
		slog.Error("Unable to subscribe: %v\n", err)
		os.Exit(1)
	}

	// Keep the application running and receiving messages
	select {}
}
