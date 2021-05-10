package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/d2r2/go-dht"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	rpio "github.com/stianeikeland/go-rpio"
)

var sessionStatus bool = true
var counter int = 0
var start = time.Now()
var TOPIC_H string = "Humidity"
var TOPIC_T string = "Temperature"

type humStruct struct {
	Humidity float64
	Unit     string
}

type tempStruct struct {
	Temperature float64
	Unit        string
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("Message received")
}

func publish(client mqtt.Client) {
	sensorType := dht.DHT11

	pin := 1
	temperature, humidity, retried, _ :=
		dht.ReadDHTxxWithRetry(sensorType, pin, false, 10)
	// if err != nil {
	// 	lg.Fatal(err)
	// }
	// print temperature and humidity
	fmt.Sprintf("Sensor = %v: Temperature = %v*C, Humidity = %v%% (retried %d times)",
		sensorType, temperature, humidity, retried)

	// num := 10
	// for i := 0; i < num; i++ {
	// 	currentTemp := tempStruct{
	// 		Temperature: 19,
	// 		Unit:        "%",
	// 	}
	// 	jsonTemp, jsonErr := json.Marshal(currentTemp)
	// 	if jsonErr != nil {
	// 		log.Fatal(jsonErr)
	// 	}
	// 	token := client.Publish(TOPIC_T, 0, false, string(jsonTemp))
	// 	token.Wait()
	// 	time.Sleep(time.Second)
	// }
	sessionStatus = false
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
	publish(client)

	// Stay in loop to receive message
	for sessionStatus {
		//Do nothing
	}

	// Disconnect
	client.Disconnect(100)

	fmt.Println("Ending run!")
}
