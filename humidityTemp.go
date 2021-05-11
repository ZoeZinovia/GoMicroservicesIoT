package main

// #cgo LDFLAGS: -lwiringPi
// #include <wiringPi.h>
// #include <stdio.h>
// #include <stdlib.h>
// #include <stdint.h>
// #define MAX_TIMINGS	85
// #define DHT_PIN		7	/* GPIO-4 */
// void read_dht_data()
// {
//  int data[5] = { 0, 0, 0, 0, 0 };
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
//		FILE *f = fopen("comment.txt", "a");
// 	   	if (f == NULL)
// 	  	{
// 	   		printf("Error opening file!\n");
// 	   		exit(1);
// 	   	}
//	   	fprintf(f, "j value: %d\n", j);
//     	fclose(f);
// 	if ( (j >= 40) &&
// 	     (data[4] == ( (data[0] + data[1] + data[2] + data[3]) & 0xFF) ) )
// 	{
// 		float h = (float)((data[0] << 8) + data[1]) / 10;
// 		if ( h > 100 )
// 		{
// 			h = data[0];	// for DHT11
// 		}
// 		float c = (float)(((data[2] & 0x7F) << 8) + data[3]) / 10;
// 		if ( c > 125 )
// 		{
// 			c = data[2];	// for DHT11
// 		}
// 		if ( data[2] & 0x80 )
// 		{
// 			c = -c;
// 		}
// 		float fT = c * 1.8f + 32;
//		FILE *f = fopen("comment.txt", "a");
// 	   	if (f == NULL)
// 	  	{
// 	   		printf("Error opening file!\n");
// 	   		exit(1);
// 	   	}
//	 	fprintf(f, "%s", "worked :))\n");
//	   	fprintf(f, "%d, %d, %d, %d, %d\n", data[0], data[1], data[2], data[3], data[4]);
//     	fclose(f);
// 	}else  {
//		FILE *f = fopen("comment.txt", "a");
// 	   	if (f == NULL)
// 	  	{
// 	   		printf("Error opening file!\n");
// 	   		exit(1);
// 	   	}
//	   	fprintf(f, "%s", "error :(\n");
//	   	fprintf(f, "%d, %d, %d, %d, %d\n", data[0], data[1], data[2], data[3], data[4]);
//     	fclose(f);
// 	}
// }
import "C"

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
	// temperatureReading, humidityReading, _, err :=
	// 	dht.ReadDHTxxWithRetry(dht.DHT11, 4, false, 10)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	C.read_dht_data()

	// fmt.Printf("%T", returnedArray)
	// byteSlice := C.GoBytes(unsafe.Pointer(&returnedArray), 5)

	// counter := 0
	// for (byteSlice[0] == 255) && (counter < 5) {
	// 	returnedArray := C.read_dht11_dat()
	// 	byteSlice = C.GoBytes(unsafe.Pointer(&returnedArray), 5)
	// 	counter++
	// }
	// if counter == 5 {
	// 	fmt.Println("Problem encountered with DHT. Please check.")
	// 	os.Exit(0)
	// }
	// mySlice := byteSliceToIntSlice(byteSlice)

	// fmt.Println(mySlice[0], mySlice[1], mySlice[2], mySlice[3], mySlice[4])
	// temperatureReading := mySlice[0] + (mySlice[1] / 10)
	// humidityReading := mySlice[2] + (mySlice[3] / 10)

	// fmt.Println("temperature:", temperatureReading, ", humidity:", humidityReading)

	// currentTemperature := tempStruct{
	// 	Temp: temperatureReading,
	// 	Unit: "C",
	// }
	// currentHumidity := humStruct{
	// 	Humidity: humidityReading,
	// 	Unit:     "%",
	// }
	// jsonTemperature := currentTemperature.structToJSON()
	// fmt.Println(string(jsonTemperature))
	// jsonHumidity := currentHumidity.structToJSON()
	// client.Publish(TOPIC_T, 0, false, string(jsonTemperature))
	// client.Publish(TOPIC_H, 0, false, string(jsonHumidity))
	// token1.Wait()
	// token2.Wait()
	// time.Sleep(time.Second)
}

func getJSON(r reading) []byte {
	return r.structToJSON()
}

func byteSliceToIntSlice(bs []byte) []int {
	result := make([]int, len(bs))
	for i, b := range bs {
		result[i] = int(b)
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
	// opts.SetClientID("go_mqtt_client")
	// opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Publish to topic
	for i := 0; i < 2; i++ {
		publish(client)
		time.Sleep(1 * time.Second)
	}

	// Disconnect
	client.Disconnect(100)
	end := time.Now()
	duration := end.Sub(start).Seconds()
	resultString := fmt.Sprint("Humidity and temperature runtime = ", duration, "\n")
	saveResultToFile("piResultsGo.txt", resultString)
	fmt.Println("Humidity and temperature runtime =", duration)
}
