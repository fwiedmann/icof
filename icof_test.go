package icof_test

import (
	"context"
	"github.com/fwiedmann/icof"
	mock_icof "github.com/fwiedmann/icof/mock"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

type ObserverMock struct {
	stateToSend icof.ObserverState
}

func (o ObserverMock) Observe(ctx context.Context, states chan<- icof.ObserverState) {
	states <- o.stateToSend
}

func TestRun_Should_Send_No_Alert_Due_To_Stored_State(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Second * 2)
		cancel()
	}()

	repoMock := mock_icof.NewMockStateRepository(ctrl)
	repoMock.EXPECT().GetLatest(gomock.Any()).Return(icof.Alert, nil).AnyTimes()
	repoMock.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)

	notifierMock := mock_icof.NewMockNotifier(ctrl)

	notifierMock.EXPECT().Resolve(gomock.Any()).Times(0)
	notifierMock.EXPECT().Alert(gomock.Any()).Times(0)
	err := icof.Run(ctx, icof.Config{
		Observer:   ObserverMock{stateToSend: icof.Alert},
		Notifiers:  []icof.Notifier{notifierMock},
		Repository: repoMock,
	})
	if err != nil {
		t.Fatalf("Run() error: %s", err)
	}
}

func TestRun_Should_Send_Alert_Due_To_Stored_State(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Second * 2)
		cancel()
	}()

	repoMock := mock_icof.NewMockStateRepository(ctrl)
	repoMock.EXPECT().GetLatest(gomock.Any()).Return(icof.Resolved, nil).AnyTimes()
	repoMock.EXPECT().Save(gomock.Any(), icof.Alert).Times(1)

	notifierMock := mock_icof.NewMockNotifier(ctrl)

	notifierMock.EXPECT().Resolve(gomock.Any()).Times(0)
	notifierMock.EXPECT().Alert(gomock.Any()).Times(1)
	err := icof.Run(ctx, icof.Config{
		Observer:   ObserverMock{stateToSend: icof.Alert},
		Notifiers:  []icof.Notifier{notifierMock},
		Repository: repoMock,
	})
	if err != nil {
		t.Fatalf("Run() error: %s", err)
	}
}

func TestRun_Should_Send_Resolved(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Second * 2)
		cancel()
	}()

	repoMock := mock_icof.NewMockStateRepository(ctrl)
	repoMock.EXPECT().GetLatest(gomock.Any()).Return(icof.Resolved, nil).AnyTimes()
	repoMock.EXPECT().Save(gomock.Any(), icof.Resolved).Times(1)

	notifierMock := mock_icof.NewMockNotifier(ctrl)

	notifierMock.EXPECT().Resolve(gomock.Any()).Times(1)
	notifierMock.EXPECT().Alert(gomock.Any()).Times(0)
	err := icof.Run(ctx, icof.Config{
		Observer:   ObserverMock{stateToSend: icof.Resolved},
		Notifiers:  []icof.Notifier{notifierMock},
		Repository: repoMock,
	})
	if err != nil {
		t.Fatalf("Run() error: %s", err)
	}
}
