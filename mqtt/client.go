package mqtt

import (
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func Connect(broker string) (MQTT.Client, error) {
	opts := MQTT.NewClientOptions().AddBroker(broker)
	opts.SetClientID("go_mqtt_client")
	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}

func Subscribe(client MQTT.Client, topic string, callback MQTT.MessageHandler) {
	if token := client.Subscribe(topic, 0, callback); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %v", token.Error())
	}
}
