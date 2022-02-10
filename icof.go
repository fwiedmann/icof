package icof

import (
	"context"
)

type ObserverState bool

const (
	Alert    ObserverState = true
	Resolved ObserverState = false
)

// Observer is responsible to send states of an occurring alert with true as the alert started and false when the alert is resolved.
type Observer interface {
	Observe(context.Context, chan<- ObserverState)
}

type Notifier interface {
	Alert(ctx context.Context) error
	Resolve(ctx context.Context) error
}

type StateRepository interface {
	Save(ctx context.Context, state ObserverState) error
	GetLatest(ctx context.Context) (ObserverState, error)
}

type Config struct {
	Observer   Observer
	Notifiers  []Notifier
	Repository StateRepository
}

// Run will handle all incoming alerts.
// It also considers to not send an alert if the last stored state was an alert and the current received state is an alert.
func Run(ctx context.Context, c Config) error {
	lastState, err := c.Repository.GetLatest(ctx)
	if err != nil {
		panic(err)
	}

	alertChan := make(chan ObserverState)
	go c.Observer.Observe(ctx, alertChan)
	for {
		select {
		case <-ctx.Done():
			return nil
		case observedSate := <-alertChan:
			if !shouldSendNotification(observedSate, lastState) {
				break
			}
			// reset after the first resolved alert is sent
			// so that the further alerts will be sent
			lastState = false
			if err := c.Repository.Save(ctx, observedSate); err != nil {
				return err
			}

			if err := handleAlert(ctx, c, observedSate); err != nil {
				return err
			}
		}
	}
}

// when the last stored state was an alert, the observer should not send any new messages
// because the notifier was already triggered
func shouldSendNotification(incomingState ObserverState, stateFromRepo ObserverState) bool {
	return !(incomingState == Alert && stateFromRepo == Alert)
}

func handleAlert(ctx context.Context, c Config, alert ObserverState) error {
	for _, notifier := range c.Notifiers {
		if alert == Alert {
			if err := notifier.Alert(ctx); err != nil {
				return err
			}
			continue
		}
		if err := notifier.Resolve(ctx); err != nil {
			return err
		}
	}
	return nil
}
