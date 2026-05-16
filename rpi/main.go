package main

import (
	"log"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	p := rpi.P1_37

	if err := p.Out(gpio.High); err != nil {
		log.Fatal(err)
	}

	log.Println("Pin 37 ON")

	time.Sleep(5 * time.Second)

	p.Out(gpio.Low)

	log.Println("Pin 37 OFF")
}
