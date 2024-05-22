package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	// MQTT Client
	opts := mqtt.NewClientOptions()
	opts.AddBroker("mqtt://localhost:1883")
	opts.SetClientID("home-matrix")
	opts.SetDefaultPublishHandler(f)
	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	c.Subscribe("pot/idx/health", 0, f)

	defer c.Disconnect(0)

	//

	// JSON data as a string
	jsonData := `{"ts":1716405092980,"temperature":10.253293,"humidity":0.068169534,"pressure":0.8683571,"moisture":3313.3716,"light":8977.43}`

	// Variable to hold the unmarshaled data
	var data PotHealth

	// Unmarshal the JSON data into the struct
	err := json.Unmarshal([]byte(jsonData), &data)
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
