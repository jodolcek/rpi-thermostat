package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func temperatura() (float64, error) {
	podaci, err := os.ReadFile("/sys/bus/w1/devices/28-062542df2985/w1_slave")
	if err != nil {
		return 0, err
	}

	s := string(podaci)

	if !strings.Contains(s, "YES") {
		return 0, fmt.Errorf("CRC fail")
	}

	temp_string := strings.TrimSpace(strings.Split(s, "t=")[1])

	temp_a, err := strconv.ParseFloat(temp_string, 64)
	if err != nil {
		return 0, err
	}

	temp_b := temp_a / 1000.0
	temp := math.Round(temp_b*2) / 2

	return temp, nil
}

func main() {

	broker := "tcp://server.apps.dj:1883"
	topic := "rpi/temperature"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("rpi-sensor")
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

	/*	for {
			token := client.Connect()
			token.Wait()

			if token.Error() == nil {
				log.Println("Initial MQTT connected")
				break
			}

			log.Println("MQTT connect failed:", token.Error())
			time.Sleep(3 * time.Second)
		}
	*/
	for {
		temp, err := temperatura()
		if err != nil {
			log.Println("Greška:", err)
		} else {
			fmt.Println("Temp:", temp)

			msg := fmt.Sprintf("%.1f", temp)

			token := client.Publish(topic, 0, false, msg)
			token.Wait()

			if token.Error() != nil {
				log.Println("Publish error:", token.Error())
			}
		}

		time.Sleep(3 * time.Second)
	}
}
