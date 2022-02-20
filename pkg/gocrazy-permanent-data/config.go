package gocrazy_permanent_data

import (
	"context"
	"encoding/json"
	"github.com/fwiedmann/icof/pkg/notifier"
	"io/ioutil"
)

// ConfigLocation is used to load the config required for icof to start.
// gokrazy stores the permanent data under the /perm directory
var ConfigLocation = "/perm/icof/start-config.json"

// New loads the startup config from disk
func New() (StartUpConfig, error) {
	content, err := ioutil.ReadFile(ConfigLocation)
	if err != nil {
		return StartUpConfig{}, err
	}

	var c StartUpConfig
	if err := json.Unmarshal(content, &c); err != nil {
		return StartUpConfig{}, err
	}

	return c, nil
}

// StartUpConfig holds all required properties for icof to start
type StartUpConfig struct {
	EmailClientConfig   EmailClientConfig   `json:"email_config"`
	EmailReceiverConfig EmailReceiverConfig `json:"email_receiver_config"`
}

type EmailClientConfig struct {
	Host             string `json:"host"`
	Port             int    `json:"port"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	FromEmailAddress string `json:"from_email_address"`
}

type EmailReceiverConfig struct {
	AlertSubject   string          `json:"alert_subject"`
	ResolveSubject string          `json:"resolve_subject"`
	Receivers      []EmailReceiver `json:"receivers"`
}

type EmailReceiver struct {
	Name                   string         `json:"name"`
	AlertTemplateMessage   string         `json:"alert_template_message"`
	ResolveTemplateMessage string         `json:"resolve_template_message"`
	Addresses              []EmailAddress `json:"addresses"`
}

type EmailAddress struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

// GetEmailReceivers implements the notifier.EmailReceiverRepository
func (c StartUpConfig) GetEmailReceivers(ctx context.Context) ([]notifier.EmailReceiver, error) {
	config, err := New()
	if err != nil {
		return nil, err
	}

	receivers := make([]notifier.EmailReceiver, 0)
	for _, receiver := range config.EmailReceiverConfig.Receivers {
		addresses := make([]notifier.Address, 0)
		for _, address := range receiver.Addresses {
			addresses = append(addresses, notifier.Address{
				Email:   address.Email,
				Name:    address.Name,
				Surname: address.Surname,
			})
		}
		receivers = append(receivers, notifier.EmailReceiver{
			Name:                   receiver.Name,
			AlertSubject:           config.EmailReceiverConfig.AlertSubject,
			ResolveSubject:         config.EmailReceiverConfig.ResolveSubject,
			AlertTemplateMessage:   receiver.AlertTemplateMessage,
			ResolveTemplateMessage: receiver.ResolveTemplateMessage,
			Addresses:              addresses,
		})
	}
	return receivers, nil
}
