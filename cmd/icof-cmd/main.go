package main

import (
	"context"
	"github.com/fwiedmann/icof"
	"github.com/fwiedmann/icof/pkg/gocrazy-permanent-data"
	"github.com/fwiedmann/icof/pkg/notifier"
	pi "github.com/fwiedmann/icof/pkg/raspberry-pi"
	"gopkg.in/gomail.v2"
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

	c, err := gocrazy_permanent_data.New()
	if err != nil {
		panic(err)
	}

	emailNotifier, err := notifier.NewEmailClient(
		&notifier.EmailClientConfig{
			Sender:           gomail.NewDialer(c.EmailClientConfig.Host, c.EmailClientConfig.Port, c.EmailClientConfig.Username, c.EmailClientConfig.Password),
			FromEmailAddress: c.EmailClientConfig.FromEmailAddress,
		},
		c,
	)

	panic(icof.Run(context.Background(), icof.Config{
		Observer:   pinAlert,
		Notifiers:  []icof.Notifier{emailNotifier},
		Repository: gocrazy_permanent_data.NewStateRepository(),
	}))
}
