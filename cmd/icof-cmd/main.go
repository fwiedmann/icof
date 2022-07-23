package main

import (
	"context"
	"github.com/fwiedmann/icof"
	"github.com/fwiedmann/icof/pkg/gocrazy-permanent-data"
	"github.com/fwiedmann/icof/pkg/notifier"
	pi "github.com/fwiedmann/icof/pkg/raspberry-pi"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"net/http"
	"time"
)

const httpServerPort = ":8080"

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
	logger := setupLogger()

	dynamicEmailConfig, err := gocrazy_permanent_data.NewDynamicEmailConfig(logger)
	if err != nil {
		panic(err)
	}

	emailNotifier, err := notifier.NewEmailClient(
		&notifier.EmailClientConfig{
			Sender:           gomail.NewDialer(c.EmailClientConfig.Host, c.EmailClientConfig.Port, c.EmailClientConfig.Username, c.EmailClientConfig.Password),
			FromEmailAddress: c.EmailClientConfig.FromEmailAddress,
		},
		dynamicEmailConfig,
		logger,
	)

	if err != nil {
		panic(err)
	}

	go func() {
		router := mux.NewRouter()
		router.HandleFunc("/email-config", dynamicEmailConfig.GetConfigHandler()).Methods(http.MethodGet)
		router.HandleFunc("/email-config", dynamicEmailConfig.CreatOrUpdateConfig()).Methods(http.MethodPut)
		panic(http.ListenAndServe(httpServerPort, router))
	}()

	panic(icof.Run(context.Background(), icof.Config{
		Observer:   pinAlert,
		Notifiers:  []icof.Notifier{emailNotifier},
		Repository: gocrazy_permanent_data.NewStateRepository(),
		Logger:     logger,
	}))
}

func setupLogger() *log.Logger {
	logger := log.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(&log.JSONFormatter{})
	return logger
}
