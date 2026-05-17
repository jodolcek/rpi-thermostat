package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Stanje struct {
	temp  float64
	point float64
}

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

func main() {
	var stanje Stanje
	var m sync.Mutex
	stanje.point = -273.0
	stanje.temp = -273.0
	stanjeCh := make(chan Stanje)

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	pin := rpi.P1_37

	for i := 0; i < 10; i++ {
		err := pin.Out(gpio.Low)
		if err == nil {
			break
		}

		log.Println("GPIO error, retry", i, ":", err)
		time.Sleep(1 * time.Second)
	}
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

		point, _ := strconv.ParseFloat(point_a, 64)
		fmt.Println("Setpoint:", point)
		m.Lock()
		stanje.point = point
		m.Unlock()
		stanjeCh <- Stanje{}

	})

	subToken.Wait()
	go func() {
		for range stanjeCh {

			m.Lock()
			t := stanje.temp
			p := stanje.point
			m.Unlock()

			if t < p {
				fmt.Println("Grijanje ON")
				if err := pin.Out(gpio.High); err != nil {
					log.Println("Greška:", err)
				}
			} else {
				fmt.Println("Grijanje OFF")
				if err := pin.Out(gpio.Low); err != nil {
					log.Println("Greška:", err)
				}
			}
		}
	}()
	go func() {
		for {
			temp_check, err := temperature()
			if err != nil {
				log.Println("Greška:", err)
			} else {
				fmt.Println("Temp:", temp_check)
				if temp_check != stanje.temp {
					fmt.Println("Nova temp:", temp_check)
					m.Lock()
					stanje.temp = temp_check
					m.Unlock()
					stanjeCh <- Stanje{}
					msg := fmt.Sprintf("%.1f", stanje.temp)

					token := client.Publish(topic, 0, false, msg)
					token.Wait()

					if token.Error() != nil {
						log.Println("Publish error:", token.Error())
					}
				}

			}

			time.Sleep(3 * time.Second)
		}
	}()
	select {}
}
