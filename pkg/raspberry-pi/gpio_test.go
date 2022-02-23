package raspberry_pi

import (
	"context"
	"github.com/fwiedmann/icof"
	"testing"
	"time"

	"periph.io/x/conn/v3/gpio"

	"periph.io/x/conn/v3/gpio/gpiotest"
)

func TestObserve(t *testing.T) {
	t.Parallel()

	type want struct {
		alertCount uint
		alertState bool
	}

	type given struct {
		alertState bool
		alertCount uint
	}

	tt := []struct {
		name     string
		pinValue bool
		want     want
		given    given
	}{
		{
			name:     "send a alert",
			pinValue: false,
			want: want{
				alertCount: 1,
				alertState: true,
			},
		},

		{
			name:     "send a resolved alert",
			pinValue: true,
			want:     want{alertCount: 0},
			given: given{
				alertState: true,
				alertCount: 0,
			},
		},
	}

	for _, tableTest := range tt {
		t.Run(tableTest.name, func(t *testing.T) {
			pin := GpioPinAlert{
				pin: &gpiotest.Pin{
					L: gpio.Level(tableTest.pinValue),
				},
				config: PinAlertConfig{
					PinDefaultState:          true,
					SleepBetweenReadDuration: time.Second * 10,
				},
				alertCount:             tableTest.given.alertCount,
				previousReadWasAnAlert: tableTest.given.alertState,
			}

			alertChan := make(chan icof.ObserverState)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
			t.Cleanup(func() {
				cancel()
			})
			go pin.Observe(ctx, alertChan)

			select {
			case <-ctx.Done():
				t.Fatal("Observe() did not send any signal through the channel")
			case alert := <-alertChan:
				if alert != icof.ObserverState(tableTest.want.alertState) {
					t.Fatalf("Observe send %t alert state, but want %t", alert, tableTest.want.alertState)
				}
				if pin.alertCount != tableTest.want.alertCount {
					t.Fatalf("Observe() internal state is incorrect. The alertCount should be %d, but is %d", pin.alertCount, tableTest.want.alertCount)
				}
			}
		})
	}
}
