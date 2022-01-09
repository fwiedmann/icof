package raspberry_pi

import (
	"context"
	"fmt"
	"github.com/fwiedmann/icof"
	"time"

	"periph.io/x/conn/v3/gpio/gpioreg"

	"periph.io/x/host/v3"

	"periph.io/x/conn/v3/gpio"
)

// PinAlertConfig defines all required parameters to init a GpioPinAlert
type PinAlertConfig struct {
	SleepBetweenReadDuration         time.Duration
	RequiredAlertsBeforeNotification uint
	Pin                              uint
	PinDefaultState                  bool
}

// NewGpioPin check if the underlying host has all relevant drivers and the given pin is present on the board.
func NewGpioPin(config PinAlertConfig) (*GpioPinAlert, error) {
	_, err := host.Init()
	if err != nil {
		return nil, err
	}

	p := gpioreg.ByName(fmt.Sprintf("GPIO%d", config.Pin))
	if p == nil {
		return nil, fmt.Errorf("could not find GPIO pin with number %d", config.Pin)
	}

	return &GpioPinAlert{
		pin:    p,
		config: config,
	}, nil
}

// GpioPinAlert struct holds the state during the GPIO pin observation
type GpioPinAlert struct {
	pin    gpio.PinIO
	config PinAlertConfig

	// state values
	alertCount             uint
	previousReadWasAnAlert bool
}

// Observe continuously checks the current state of the GPIO pin (high or low).
// An alert will be sent if the threshold of required alerts is met. An alert is defined as the opposite of the PinAlertConfig.PinDefaultState.
// Once the pin state equals again with PinAlertConfig.PinDefaultState, a resolved alert will be sent and the GpioPinAlert struct state will be set to its default values.
func (gp *GpioPinAlert) Observe(ctx context.Context, alertChan chan<- icof.ObserverState) {
	for err := ctx.Err(); err == nil; {
		pinState := gp.pin.Read()

		if bool(pinState) != gp.config.PinDefaultState {
			gp.alertCount++
		}

		// resolve alert when pin in his its default state again after the last read was an alert.
		// reset alertCount and previousReadWasAnAlert
		if bool(pinState) == gp.config.PinDefaultState && gp.previousReadWasAnAlert {
			alertChan <- icof.Resolved
			gp.previousReadWasAnAlert = false
			gp.alertCount = 0
		}

		// send alert if the threshold of required alerts is
		if (gp.alertCount >= gp.config.RequiredAlertsBeforeNotification) && !gp.previousReadWasAnAlert {
			alertChan <- icof.Alert
			gp.previousReadWasAnAlert = true
		}

		time.Sleep(gp.config.SleepBetweenReadDuration)
	}
}
