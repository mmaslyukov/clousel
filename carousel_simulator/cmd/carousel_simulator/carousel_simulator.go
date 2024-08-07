package main

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type EventAck struct {
	CarouselId    string `json:CarouselId`
	SequenceNum   int    `json:SequenceNum`
	CorrelationId string `json:CorrelationId`
	Type          string `json:Type`
}

type MessageGeneral struct {
	CarouselId  string `json:CarouselId`
	SequenceNum int    `json:SequenceNum`
	EventId     string `json:EventId`
	Type        string `json:Type`
}

type MessageCommand struct {
	MessageGeneral
	Command string `json:Command`
}

var seq int
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	var cmd MessageCommand
	if err := json.Unmarshal(msg.Payload(), &cmd); err != nil {
		fmt.Printf("Fail to unmarshal:%s", err)
		return
	}
	seq++
	ack := EventAck{
		CarouselId:    cmd.CarouselId,
		SequenceNum:   seq,
		CorrelationId: cmd.EventId,
		Type:          "EventAck",
	}

	if payload, err := json.Marshal(&ack); err == nil {
		publish(fmt.Sprintf("/clousel/carousel/%s", ack.CarouselId), payload, client)
	} else {
		fmt.Printf("Fail to nmarshal:%s", err)

	}

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func sub(client mqtt.Client) {
	topic := "/clousel/cloud/#"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s\n", topic)
}

func publish(topic string, payload []byte, client mqtt.Client) {
	fmt.Printf("Publish to topic:'%s', payload:%s", topic, payload)
	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	time.Sleep(time.Second)
}
func heartbeat() {
}
func main() {
	var broker = "192.168.0.150"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	// opts.SetUsername("emqx")
	// opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())

	}

	sub(client)
	for {
		time.Sleep(time.Hour)
	}
}
