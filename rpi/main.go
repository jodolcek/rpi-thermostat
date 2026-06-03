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
	"github.com/joho/godotenv"
)

type State struct {
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
	temp := math.Round(temp_b*10) / 10
	return temp, nil
}

func main() {

	var state State
	var m sync.Mutex
	var heating bool
	const hysteresis = 0.5
	state.point = 0.0
	state.temp = 0.0
	stateCh := make(chan State)

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
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env")
	}
	mqttuser := os.Getenv("mqtt_user")
	mqttpasswd := os.Getenv("mqtt_passwd")
	broker := "tcp://server.apps.dj:1883"
	topic := "rpi/temperature"
	topic2 := "rpi/heating"
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetUsername(mqttuser)
	opts.SetPassword(mqttpasswd)
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
	msg := fmt.Sprintf("%.1f", state.point)

	point := client.Publish("rpi/setpoint", 0, true, msg)
	point.Wait()

	if point.Error() != nil {
		log.Println("Publish error:", token.Error())
	}
	h := client.Publish(topic2, 0, true, "off")
	h.Wait()

	if h.Error() != nil {
		log.Println("Publish error:", token.Error())
	}
	subToken := client.Subscribe("rpi/setpoint", 0, func(client mqtt.Client, msg mqtt.Message) {
		point_a := string(msg.Payload())

		point, _ := strconv.ParseFloat(point_a, 64)
		fmt.Println("Setpoint:", point)
		m.Lock()
		state.point = point
		m.Unlock()
		stateCh <- State{}

	})

	subToken.Wait()
	go func() {
		for range stateCh {

			m.Lock()
			t := state.temp
			p := state.point
			m.Unlock()

			if !heating && t < (p-hysteresis) {
				heating = true
				fmt.Println("Grijanje ON")
				pin.Out(gpio.High)
				msg := "on"

				token := client.Publish(topic2, 0, true, msg)
				token.Wait()

				if token.Error() != nil {
					log.Println("Publish error:", token.Error())
				}
			}

			if heating && t > (p+hysteresis) {
				heating = false
				fmt.Println("Grijanje OFF")
				pin.Out(gpio.Low)
				msg := "off"

				token := client.Publish(topic2, 0, true, msg)
				token.Wait()

				if token.Error() != nil {
					log.Println("Publish error:", token.Error())
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
				if temp_check != state.temp {
					fmt.Println("Nova temp:", temp_check)
					m.Lock()
					state.temp = temp_check
					m.Unlock()
					stateCh <- State{}
					msg := fmt.Sprintf("%.1f", state.temp)

					token := client.Publish(topic, 0, true, msg)
					token.Wait()

					if token.Error() != nil {
						log.Println("Publish error:", token.Error())
					}
				}

			}

			time.Sleep(5 * time.Second)
		}
	}()
	select {}
}
