package main

// #cgo LDFLAGS: -lwiringPi
// #include <wiringPi.h>
// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
// // Pi dht11 variables
// #define MAXTIMINGS	85
// #define DHTPIN		7
// int dht11_dat[5] = { 0, 0, 0, 0, 0 };
// // Reading of the dht11 is rather complex in C/C++. See this site that explains how readings are made: http://www.uugear.com/portfolio/dht11-humidity-temperature-sensor-module/
// int* read_dht11_dat()
// {
//	   wiringPiSetupGPIO();
//     u_int8_t laststate	= HIGH;
//     u_int8_t counter		= 0;
//     u_int8_t j		= 0, i;
//	   pinMode(32, OUTPUT);
//	   digitalWrite(32, 0);
//     dht11_dat[0] = dht11_dat[1] = dht11_dat[2] = dht11_dat[3] = dht11_dat[4] = 0;
//     // pull pin down for 18 milliseconds. This is called “Start Signal” and it is to ensure DHT11 has detected the signal from MCU.
//     pinMode( DHTPIN, OUTPUT );
//     digitalWrite( DHTPIN, LOW );
//     delay( 18 );
//     // Then MCU will pull up DATA pin for 40us to wait for DHT11’s response.
//     digitalWrite( DHTPIN, HIGH );
//     delayMicroseconds( 40 );
//     // Prepare to read the pin
//     pinMode( DHTPIN, INPUT );
//     // Detect change and read data
//     for ( i = 0; i < MAXTIMINGS; i++ )
//     {
//         counter = 0;
//         while ( digitalRead( DHTPIN ) == laststate )
//         {
//             counter++;
//             delayMicroseconds( 1 );
//             if ( counter == 255 )
//             {
//                 break;
//             }
//         }
//         laststate = digitalRead( DHTPIN );
//         if ( counter == 255 )
//             break;
//         // Ignore first 3 transitions
//         if ( (i >= 4) && (i % 2 == 0) )
//         {
//             // Add each bit into the storage bytes
//             dht11_dat[j / 8] <<= 1;
//             if ( counter > 16 )
//                 dht11_dat[j / 8] |= 1;
//             j++;
//         }
//     }
//     // Check that 40 bits (8bit x 5 ) were read + verify checksum in the last byte
//     if ( (j >= 40) && (dht11_dat[4] == ( (dht11_dat[0] + dht11_dat[1] + dht11_dat[2] + dht11_dat[3]) & 0xFF) ) )
//     {
//		   FILE *f = fopen("file.txt", "w");
// 		   if (f == NULL)
// 		   {
// 		   		printf("Error opening file!\n");
// 		   		exit(1);
// 		   }
//		   fprintf(f, "Temp: %d, %d, Humidity: %d, %d\n", dht11_dat[0], dht11_dat[1], dht11_dat[2], dht11_dat[3]);
//         fclose(f);
//		   return dht11_dat; // If all ok, return pointer to the data array
//     } else  {
//	       FILE *f = fopen("file.txt", "w");
// 		   if (f == NULL)
// 		   {
// 		   		printf("Error opening file!\n");
// 		   		exit(1);
// 		   }
//		   fprintf(f, "Temp: %d, %d, Humidity: %d, %d\n", dht11_dat[0], dht11_dat[1], dht11_dat[2], dht11_dat[3]);
//		   fclose(f);
//         dht11_dat[0] = 255;
//         return dht11_dat; //If there was an error, set first array element to -1 as flag to main function
//     }
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
	"unsafe"

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

	returnedArray := C.read_dht11_dat()

	fmt.Printf("%T", returnedArray)
	byteSlice := C.GoBytes(unsafe.Pointer(&returnedArray), 5)

	counter := 0
	for (byteSlice[0] == 255) && (counter < 5) {
		returnedArray := C.read_dht11_dat()
		byteSlice = C.GoBytes(unsafe.Pointer(&returnedArray), 5)
		counter++
	}
	if counter == 5 {
		fmt.Println("Problem encountered with DHT. Please check.")
		os.Exit(0)
	}
	mySlice := byteSliceToIntSlice(byteSlice)

	fmt.Println(mySlice[0], mySlice[1], mySlice[2], mySlice[3], mySlice[4])
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
	for i := 0; i < 10; i++ {
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
