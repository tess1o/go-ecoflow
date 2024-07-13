package main

import (
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tess1o/go-ecoflow"
	"log"
	"time"
)

func main() {
	// device serial number. MQTT doesn't have a way to get all linked devices (client via REST API has such possibility).
	var deviceSn = "Ecoflow device serial number "

	mqttClientConfig := ecoflow.MqttClientConfiguration{
		Email:            "ecoflow_email_address",
		Password:         "ecoflow_password",
		OnConnect:        connectHandler,     //can be nil
		OnConnectionLost: connectLostHandler, // can be nil
		OnReconnect:      nil,                //can be nil
	}

	//an error can be returned if wrong login/password is provided, ecoflow api is not available, network issue, etc
	client, err := ecoflow.NewMqttClient(context.Background(), mqttClientConfig)
	if err != nil {
		log.Fatalf("Unable to create mqtt client: %+v\n", err)
	}

	//connect to MQTT broker and subscribe to the device's topic where its parameters are published
	// It's not described in documentation,
	// however looks like it sends to the topic only parameters that are changed, not all list of current values
	err = client.SubscribeForParameters(deviceSn, messagePubHandler)
	if err != nil {
		log.Fatalf("Unable to subscribe: %+v\n", err)
	}
	// keep receiving parameters for 1 hour
	time.Sleep(1 * time.Hour)
}

// handle payload - device's parameters
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var params ecoflow.MqttDeviceParams
	err := json.Unmarshal(msg.Payload(), &params)
	if err != nil {
		fmt.Printf("Unable to parse message %s from topic %s due to error: %+v\n", msg.Payload(), msg.Topic(), err)
	} else {
		fmt.Printf("Received message: %+v from topic: %s\n", params, msg.Topic())
	}
}

// executes when we're successfully connected to the mqtt broker
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	optionsReader := client.OptionsReader()
	fmt.Printf("Connected to the broker: %s\n", optionsReader.Servers()[0].String())
}

// executes when the mqtt connection is lost
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}
