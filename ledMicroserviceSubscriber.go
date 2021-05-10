package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	rpio "github.com/stianeikeland/go-rpio"
)

var sessionStatus bool = true
var counter int = 0
var start = time.Now()
var TOPIC string = "LED"

type ledStruct struct {
	LED_1 bool
	GPIO  int
}

func saveResultToFile(filename string, result string) {
	byteSlice := []byte(result)
	someError := ioutil.WriteFile(filename, byteSlice, 0666)
	if someError != nil {
		fmt.Println(someError)
		os.Exit(1)
	}
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	counter++
	if counter == 1 {
		start = time.Now()
	}
	var led ledStruct
	ledPin := rpio.Pin(12)
	if strings.Contains(string(msg.Payload()), "Done") {
		sessionStatus = false
		ledPin.Output()
		ledPin.Low()
		end := time.Now()
		duration := end.Sub(start).Seconds()
		fmt.Println(duration)
		resultString := fmt.Sprint("LED subsriber runtime =", duration)
		saveResultToFile("piResultsGo.txt", resultString)
	} else {
		json.Unmarshal([]byte(msg.Payload()), &led)
		ledPin = rpio.Pin(led.GPIO)
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
	fmt.Printf("Connection lost: %v", err)
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

	// Subscribe to topic
	sub(client)

	// Stay in loop to receive message
	for sessionStatus {
		//Do nothing
	}

	// Disconnect
	client.Disconnect(100)

	fmt.Println("Ending run!")
}

func sub(client mqtt.Client) {
	token := client.Subscribe(TOPIC, 1, nil)
	token.Wait()
}

// func publish(client mqtt.Client) {
// 	num := 10
// 	for i := 0; i < num; i++ {
// 		text := fmt.Sprintf("Message %d", i)
// 		token := client.Publish("led", 0, false, text)
// 		token.Wait()
// 		time.Sleep(time.Second)
// 	}
// }
