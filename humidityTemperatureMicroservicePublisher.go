package main

// #cgo LDFLAGS: -lwiringPi
// #include <wiringPi.h>
// #include <stdio.h>
// #include <stdlib.h>
// #include <stdint.h>
// #include <string.h>
// #include <time.h>
// #include <unistd.h>
// #define MAX_TIMINGS	85
// #define DHT_PIN		7	/* GPIO-4 */
// int data[5] = { 0, 0, 0, 0, 0 };
// clock_t timer = 0;
// int read_dht_data()
// {
//	wiringPiSetup();
// 	uint8_t laststate	= HIGH;
// 	uint8_t counter		= 0;
// 	uint8_t j			= 0, i;
// 	data[0] = data[1] = data[2] = data[3] = data[4] = 0;
// 	/* pull pin down for 18 milliseconds */
// 	pinMode( DHT_PIN, OUTPUT );
// 	digitalWrite( DHT_PIN, LOW );
// 	delay( 18 );
// 	/* prepare to read the pin */
// 	digitalWrite( DHT_PIN, HIGH);
//  delayMicroseconds( 40 );
// 	pinMode( DHT_PIN, INPUT );
// 	/* detect change and read data */
// 	for ( i = 0; i < MAX_TIMINGS; i++ )
// 	{
// 		counter = 0;
// 		while ( digitalRead( DHT_PIN ) == laststate )
// 		{
// 			counter++;
// 			delayMicroseconds( 2 );
// 			if ( counter == 255 )
// 				break;
// 		}
// 		laststate = digitalRead( DHT_PIN );
// 		if ( counter == 255 ){
// 			break;
//		}
// 		/* ignore first 3 transitions */
// 		if ( (i >= 4) && (i % 2 == 0) )
// 		{
// 			/* shove each bit into the storage bytes */
// 			data[j / 8] <<= 1;
// 			if ( counter > 16 )
// 				data[j / 8] |= 1;
// 			j++;
// 		}
// 	}
// 	/*
// 	 * check we read 40 bits (8bit x 5 ) + verify checksum in the last byte
// 	 * print it out if data is good
// 	 */
// 	if ( (j >= 40) &&
// 	     (data[4] == ( (data[0] + data[1] + data[2] + data[3]) & 0xFF) ) )
// 	{
//		FILE *f = fopen("reading.txt", "w");
// 	   	if (f == NULL)
// 	  	{
// 	   		printf("Error opening file!\n");
// 	   		exit(1);
// 	   	}
//	   	fprintf(f, "%d,%d,%d,%d,%d", data[0], data[1], data[2], data[3], data[4]);
//     	fclose(f);
//		timer = clock();
//		return data[0];
// 	} else  {
//		FILE *f = fopen("reading.txt", "w");
// 	   	if (f == NULL)
// 	  	{
// 	   		printf("Error opening file!\n");
// 	   		exit(1);
// 	   	}
//	   	fprintf(f, "%d,%d,%d,%d,%d", data[0], data[1], data[2], data[3], data[4]);
// 		fclose(f);
// 		timer = clock();
//		return data[0];
// 	}
// }
import "C"

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var sessionStatus bool = true
var counter int = 0
var start = time.Now()
var TOPIC_H string = "Humidity"
var TOPIC_T string = "Temperature"
var ADDRESS string
var PORT = 1883
var temperatureReading float32 = 0
var humidityReading float32 = 0

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
	if !sessionStatus {
		doneString := "{\"Done\": \"True\"}"
		client.Publish(TOPIC_T, 0, false, doneString)
		client.Publish(TOPIC_H, 0, false, doneString)
		return
	}
	if temperatureReading == 0 && humidityReading == 0 {
		C.read_dht_data()
		byteSlice, readErr := ioutil.ReadFile("reading.txt")
		if readErr != nil {
			log.Fatal(readErr)
		}
		mySlice := byteSliceToIntSlice(byteSlice)
		temperatureReading = float32(mySlice[2] + (mySlice[3] / 10))
		humidityReading = float32(mySlice[0] + (mySlice[1] / 10))
		end := time.Now()
		duration := end.Sub(start).Seconds()
		fmt.Println("Humidity and temperature runtime after readings =", duration)

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
	jsonHumidity := currentHumidity.structToJSON()
	client.Publish(TOPIC_T, 0, false, string(jsonTemperature))
	client.Publish(TOPIC_H, 0, false, string(jsonHumidity))
	return
}

func getJSON(r reading) []byte {
	return r.structToJSON()
}

func byteSliceToIntSlice(bs []byte) []int {
	strings := strings.Split(string(bs), ",")
	result := make([]int, len(strings))
	for i, s := range strings {
		if len(s) == 0 {
			continue
		}
		n, convErr := strconv.Atoi(s)
		if convErr != nil {
			log.Fatal(convErr)
		}
		result[i] = n
	}
	return result
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
	start = time.Now()
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
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Publish to topic
	numIterations := 100
	for i := 0; i < numIterations; i++ {
		if i == numIterations-1 {
			sessionStatus = false
		}
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
