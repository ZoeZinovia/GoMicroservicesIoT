# GoMicroservicesIoT

A microservices IoT application was developed for Unix devices.

## Description

This project implements a microservices IoT application that that can run on Linux devices such as Raspberry Pi. The code was developed as part of a study into the comparative performance of 3 languages (Go, C++ and Python) and 2 architectures (microservices vs. monolithic). There are 3 principle IoT microservices: temperature and humidity sensing, passive infrared sensing and led actuation. The microservices communicate via MQTT, a lightweight asynchronous messaging protocol. The following figure provides an example of such a microservices setup from another study:

![image showing microservices architecture with MQTT](https://devblog.axway.com/wp-content/uploads/MQTT_2.png)

## Getting Started

### Dependencies

There are a number of requirements for this project. All requirements must be installed on the embedded device, e.g. Raspberry Pi:

* Go must be [installed](https://golang.org/dl/)
* Go ROUTE and Go PATH must be configured according to your preferences. See [here](https://www.geeksforgeeks.org/golang-gopath-and-goroot/) for more information.
* Eclipse Paho MQTT is [required](https://github.com/eclipse/paho.mqtt.golang)
* Go-rpio is [required](https://github.com/stianeikeland/go-rpio) to interact with the GPIO pins.
* WiringPi is [required]() since Go-rpio doesn't interface with a dht11 sensor.
* There are a few other Go library imports that should not require additional installations. 

Additionally, MQTT Mosquitto or another MQTT broker must be installed on the device that will receive messages from the embedded device. See [here](https://mosquitto.org/) for more information.

### Installing and using

* Simply clone the code from this repository. Visual Studio Code or Go Land are recommended IDEs but the code can be run in the terminal as well.
* If running in terminal, you can compile the files by running the ```compileGosScripts``` bash script.
* Since each microservice runs independently, the compiled code also need to be executed independently. All 3 microservices can be run in parallel with the following bash command ``` bash runGoScripts <embedded device ip address>```
* In order to receive messages with temperature, humidity and PIR information, you will need to subscribe to the MQTT topic. You can do this by running the Python code in this [repository](https://github.com/ZoeZinovia/PythonMicroserviceSubscriber).
* Alternatively you can run the following command to receive messages in terminal:
```mosquitto_sub -h <embedded device ip address> -t Temperature```
* The possible topics are: Temperature, Humidity or PIR
* You can also publish messages to the Led actuator with the following command:
```mosquitto_pub -h <embedded device ip address> -m <value> -t LED```
  where value is True if you want the led to turn on or false otherwise
