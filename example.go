package main

import (
	"fmt"

	"github.com/MichaelS11/go-dht"
)

const GPIO = "GPIO4"

func main() {
	hosterr := dht.HostInit()
	if hosterr != nil {
		fmt.Println("HostInit error:", hosterr)
		return
	}

	dht, dhterr := dht.NewDHT(GPIO, dht.Fahrenheit, "")
	if dhterr != nil {
		fmt.Println("NewDHT error:", dhterr)
		return
	}

	humidity, temperature, readerr := dht.Read()

	if readerr != nil {
		fmt.Println("Reader error:", readerr)
		return
	}
	fmt.Println("Humidity:", humidity, ", Temperature:", temperature)
}
