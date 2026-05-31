package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {

	broker := "localhost:1883"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("backend")
	opts.SetCleanSession(true)

	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(3 * time.Second)

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Println("MQTT lost connection:", err)
	})

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Println("MQTT connected / reconnected")
	})

	client := mqtt.NewClient(opts)
	token := client.Connect()

	token.Wait()

	subToken := client.Subscribe("rpi/temperature", 0, func(client mqtt.Client, msg mqtt.Message) {
		temp_a := string(msg.Payload())

		temp, _ := strconv.ParseFloat(temp_a, 64)
		fmt.Println("Temperatura:", temp)

	})
	subToken.Wait()
	select {}
}
