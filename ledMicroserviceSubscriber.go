package main

import (
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	rpio "github.com/stianeikeland/go-rpio"
)

var ADDRESS string

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("IP address must be provided as a command line argument")
		os.Exit(1)
	}
	ADDRESS = os.Args[1]
	fmt.Println(ADDRESS)

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	fmt.Println("Finished!")
}
