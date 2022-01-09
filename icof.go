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

func Run(c Config) {
	alertChan := make(chan ObserverState)

	// TODO diff last stored state with current observer state. If they diff, send alert or resolve message (this will be needed when icof was shutdown)
	go c.Observer.Observe(context.Background(), alertChan)
	for {
		select {
		case observedSate := <-alertChan:
			for _, notifier := range c.Notifiers {
				if observedSate == Alert {
					notifier.Alert(context.Background())
					continue
				}
				notifier.Resolve(context.Background())
			}
		}
	}
}
