package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jackc/pgx/v5"
)

type PotHealth struct {
	TS          int64   `json:"ts"`
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	Pressure    float32 `json:"pressure"`
	Moisture    float32 `json:"moisture"`
	Light       float32 `json:"light"`
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var another_f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("ANOTHER_TOPIC: %s\n", msg.Topic())
	fmt.Printf("ANOTHER_MSG: %s\n", msg.Payload())
}

func main() {
	// TimescaleDB
	ctx := context.Background()
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASSWORD")

	connStr := fmt.Sprintf("postgres://%s:%s@localhost:5432", user, password)
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	// MQTT Client
	opts := mqtt.NewClientOptions()
	opts.AddBroker("mqtt://localhost:1883")
	opts.SetClientID("home-matrix")
	opts.SetDefaultPublishHandler(f)
	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer c.Disconnect(0)

	c.Subscribe("pot/idx/health", 0, f)

	// JSON data as a string
	jsonData := `{"ts":1716405092980,"temperature":10.253293,"humidity":0.068169534,"pressure":0.8683571,"moisture":3313.3716,"light":8977.43}`
	ensureTableExists(conn)
	batchInsert(conn)

	// Variable to hold the unmarshaled data
	var data PotHealth

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	token := c.Publish("pot/idx/health", 0, false, jsonData)
	token.Wait()

	time.Sleep(time.Second * 5)

	// Print the struct to verify the data
	fmt.Printf("Timestamp: %d\n", data.TS)
	fmt.Printf("Temperature: %.6f\n", data.Temperature)
	fmt.Printf("Humidity: %.9f\n", data.Humidity)
	fmt.Printf("Pressure: %.7f\n", data.Pressure)
	fmt.Printf("Moisture: %.4f\n", data.Moisture)
	fmt.Printf("Light: %.2f\n", data.Light)

	exit := make(chan os.Signal)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

}

func ensureTableExists(conn *pgx.Conn) {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS sensor_data (
        ts BIGINT,
        temperature DOUBLE PRECISION,
        humidity DOUBLE PRECISION,
        pressure DOUBLE PRECISION,
        moisture DOUBLE PRECISION,
        light DOUBLE PRECISION
    );
    `

	_, err := conn.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Unable to create table: %v\n", err)
	}

	fmt.Println("Ensured table exists or created it if not.")
}

func batchInsert(conn *pgx.Conn) {
	// Start a transaction
	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("Unable to begin transaction: %v\n", err)
	}
	defer tx.Rollback(context.Background())

	// Prepare the batch
	batch := &pgx.Batch{}
	for i := 0; i < 1000; i++ { // Example: Insert 1000 rows
		batch.Queue(`
            INSERT INTO sensor_data (ts, temperature, humidity, pressure, moisture, light)
            VALUES ($1, $2, $3, $4, $5, $6)
        `, 1716405092980+i, 10.253293, 0.068169534, 0.8683571, 3313.3716, 8977.43)
	}

	// Send the batch and get results
	br := tx.SendBatch(context.Background(), batch)
	_, err = br.Exec()
	if err != nil {
		log.Fatalf("Unable to execute batch: %v\n", err)
	}
	err = br.Close()
	if err != nil {
		log.Fatalf("Unable to close batch result: %v\n", err)
	}

	// Commit the transaction
	err = tx.Commit(context.Background())
	if err != nil {
		log.Fatalf("Unable to commit transaction: %v\n", err)
	}

	fmt.Println("Batch inserted successfully!")
}
