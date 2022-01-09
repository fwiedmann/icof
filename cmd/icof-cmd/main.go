package main

import (
	"github.com/fwiedmann/icof"
	pi "github.com/fwiedmann/icof/pkg/raspberry-pi"
	"time"
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

	icof.Run(icof.Config{
		Observer:   pinAlert,
		Notifiers:  []icof.Notifier{},
		Repository: nil,
	})
}
