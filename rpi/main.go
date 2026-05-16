package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func temperatura() (string, error) {
	podaci, err := os.ReadFile("/sys/bus/w1/devices/28-062542df2985/w1_slave")
	if err != nil {
		return "0", err
	}

	s := string(podaci)

	if !strings.Contains(s, "YES") {
		return "0", fmt.Errorf("CRC fail")
	}

	temp := strings.TrimSpace((strings.Split(s, "t="))[1])

	return temp, nil
}
func main() {

	for {
		temp, err := temperatura()
		if err != nil {
			log.Println("Greška", err)
		} else {
			fmt.Println(temp)
		}

		time.Sleep(30 * time.Second)
	}
}
