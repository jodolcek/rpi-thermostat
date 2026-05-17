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

func temperature() (float64, error) {
	data, err := os.ReadFile("/sys/bus/w1/devices/28-062542df2985/w1_slave")
	if err != nil {
		return 0, err
	}

	s := string(data)

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
func on_off(temp float64, point float64) {
	if temp < point {
		fmt.Println("Grijanje upaljeno")
	} else if temp > point {
		fmt.Println("Grijanje ugašeno")
	}
}
func main() {
	point := -273.0
	temp := -273.0
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

	subToken := client.Subscribe("rpi/setpoint", 0, func(client mqtt.Client, msg mqtt.Message) {
		point_a := string(msg.Payload())

		point, _ = strconv.ParseFloat(point_a, 64)
		fmt.Println("Setpoint:", point)
		on_off(temp, point)
	})

	subToken.Wait()

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
	go func() {
		for {
			temp_check, err := temperature()
			if err != nil {
				log.Println("Greška:", err)
			} else {
				fmt.Println("Temp:", temp_check)
				if temp_check != temp {
					fmt.Println("Nova temp:", temp_check)
					temp = temp_check
					on_off(temp, point)
				}
				msg := fmt.Sprintf("%.1f", temp)

				token := client.Publish(topic, 0, false, msg)
				token.Wait()

				if token.Error() != nil {
					log.Println("Publish error:", token.Error())
				}
			}

			time.Sleep(3 * time.Second)
		}
	}()
	select {}
}
