package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	rpio "github.com/stianeikeland/go-rpio"
)

var sessionStatus bool = true

type ledStruct struct {
	LED_1 bool
	GPIO  int
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var led ledStruct
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	if strings.Contains(msg.Payload(), "Done") {
		sessionStatus = false
	} else {
		json.Unmarshal([]byte(msg.Payload()), &led)
		ledPin := rpio.Pin(led.GPIO)
		ledPin.Output()
		if led.LED_1 {
			ledPin.High()
		} else {
			ledPin.Low()
		}
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

var ADDRESS string
var PORT = 1883

func main() {

	// Save the IP address
	if len(os.Args) <= 1 {
		fmt.Println("IP address must be provided as a command line argument")
		os.Exit(1)
	}
	ADDRESS = os.Args[1]
	fmt.Println(ADDRESS)

	// Check that RPIO opened correctly
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// End program with ctrl-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

	// Creat MQTT client
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", ADDRESS, PORT))
	opts.SetClientID("go_mqtt_client")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client)
	// publish(client)
	// topic := "led"
	// token := client.Subscribe(topic, 0, nil)
	// token.Wait()
	// token := client.Publish(topic, 0, false, "Hello!")
	// token.Wait()

	// ledPin := rpio.Pin(12)
	// ledPin.Output()
	// ledPin.Low()
	// time.Sleep(10 * time.Second)

	for sessionStatus {
		//Do nothing
	}
	client.Disconnect(100)

	fmt.Println("Subscribed to the topic!")
}

func sub(client mqtt.Client) {
	topic := "led"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}

func publish(client mqtt.Client) {
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish("led", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}
