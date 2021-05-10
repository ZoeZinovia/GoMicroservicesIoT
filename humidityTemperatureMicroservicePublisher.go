package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/d2r2/go-dht"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var sessionStatus bool = true
var counter int = 0
var start = time.Now()
var TOPIC_H string = "Humidity"
var TOPIC_T string = "Temperature"
var ADDRESS string
var PORT = 1883

type tempStruct struct {
	Temp float32
	Unit string
}

type humStruct struct {
	Humidity float32
	Unit     string
}

type reading interface {
	structToJSON() []byte
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("Message received")
}

func publish(client mqtt.Client) {
	temperatureReading, humidityReading, _, err :=
		dht.ReadDHTxxWithRetry(dht.DHT11, 4, false, 10)
	if err != nil {
		log.Fatal(err)
	}
	currentTemperature := tempStruct{
		Temp: temperatureReading,
		Unit: "C",
	}
	currentHumidity := humStruct{
		Humidity: humidityReading,
		Unit:     "%",
	}
	jsonTemperature := currentTemperature.structToJSON()
	fmt.Println(string(jsonTemperature))
	jsonHumidity := currentHumidity.structToJSON()
	client.Publish(TOPIC_T, 0, false, string(jsonTemperature))
	client.Publish(TOPIC_H, 0, false, string(jsonHumidity))
	// token1.Wait()
	// token2.Wait()
	// time.Sleep(time.Second)
}

func getJSON(r reading) []byte {
	return r.structToJSON()
}

func (ts tempStruct) structToJSON() []byte {
	jsonReading, jsonErr := json.Marshal(ts)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return jsonReading
}

func (ts humStruct) structToJSON() []byte {
	jsonReading, jsonErr := json.Marshal(ts)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return jsonReading
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

func saveResultToFile(filename string, result string) {
	file, errOpen := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errOpen != nil {
		log.Fatal(errOpen)
	}
	byteSlice := []byte(result)
	_, errWrite := file.Write(byteSlice)
	if errWrite != nil {
		log.Fatal(errWrite)
	}
}

func main() {

	// Save the IP address
	if len(os.Args) <= 1 {
		fmt.Println("IP address must be provided as a command line argument")
		os.Exit(1)
	}
	ADDRESS = os.Args[1]
	fmt.Println(ADDRESS)

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
	for i := 0; i < 100; i++ {
		publish(client)
	}

	// Disconnect
	client.Disconnect(100)
	end := time.Now()
	duration := end.Sub(start).Seconds()
	resultString := fmt.Sprint("Humidity and temperature runtime = ", duration, "\n")
	saveResultToFile("piResultsGo.txt", resultString)
	fmt.Println("Humidity and temperature runtime =", duration)
}
