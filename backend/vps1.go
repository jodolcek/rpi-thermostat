package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
)

type Informations struct {
	Temperature float64 `json:"temperature"`
	Heating     string  `json:"heating"`
	Setpoint    float64 `json:"setpoint"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var mqttClient mqtt.Client

var clients = make(map[*websocket.Conn]bool)

var mu sync.Mutex

var informations = Informations{
	Temperature: 0.0,
	Heating:     "off",
	Setpoint:    0.0,
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	broadcastInformations()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

	mu.Lock()
	delete(clients, conn)
	mu.Unlock()

	conn.Close()
}

func setpointHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Setpoint float64 `json:"setpoint"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	informations.Setpoint = data.Setpoint
	setpoint := fmt.Sprintf("%.1f", data.Setpoint)
	token := mqttClient.Publish("rpi/setpoint", 0, true, setpoint)
	token.Wait()

	fmt.Println("REST setpoint:", setpoint)

	broadcastInformations()

	w.WriteHeader(http.StatusOK)
}

func broadcastInformations() {
	data, err := json.Marshal(informations)
	if err != nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}

}

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
	mqttClient = client
	tmp := client.Subscribe("rpi/temperature", 0, func(client mqtt.Client, msg mqtt.Message) {
		temp_a := string(msg.Payload())

		temp, err := strconv.ParseFloat(temp_a, 64)
		if err != nil {
			return
		}
		informations.Temperature = temp
		fmt.Println("Temperatura:", temp)
		broadcastInformations()

	})
	tmp.Wait()
	h := client.Subscribe("rpi/heating", 0, func(client mqtt.Client, msg mqtt.Message) {
		heating := string(msg.Payload())
		informations.Heating = heating
		fmt.Println("Grijanje:", heating)
		broadcastInformations()
	})
	h.Wait()
	point := client.Subscribe("rpi/setpoint", 0, func(client mqtt.Client, msg mqtt.Message) {
		setpoint_a := string(msg.Payload())
		setpoint, err := strconv.ParseFloat(setpoint_a, 64)
		if err != nil {
			return
		}
		informations.Setpoint = setpoint
		fmt.Println("Postavljena temperatura:", setpoint)
		broadcastInformations()
	})
	point.Wait()
	http.HandleFunc("/ws", handleWS)
	http.HandleFunc("/setpoint", setpointHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
