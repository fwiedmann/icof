package main

import (
	"context"
	"fmt"
	"time"

	pi "github.com/fwiedmann/icof/pkg/raspberry-pi"
)

func main() {
	pinAlert, err := pi.NewGpioPin(
		pi.PinAlertConfig{
			SleepBetweenReadDuration:         1 * time.Second,
			RequiredAlertsBeforeNotification: 5,
			Pin:                              17,
			PinDefaultState:                  true,
		})

	if err != nil {
		panic(err)
	}

	run(pinAlert)
}

// Observer is responsible to send states of an occurring alert with true as the alert started and false when the alert is resolved.
type Observer interface {
	Observe(context.Context, chan<- bool)
}

func run(o Observer) {
	alertChan := make(chan bool)
	go o.Observe(context.Background(), alertChan)

	for {
		select {
		case alert := <-alertChan:
			fmt.Println(alert)
		}
	}
}
