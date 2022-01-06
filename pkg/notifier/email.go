package notifier

import (
	"context"
	"github.com/go-playground/validator/v10"
	"gopkg.in/gomail.v2"
)

// Message struct will be used to configure the send operation
type Message struct {
	Content   string   `validate:"required"`
	Receivers []string `validate:"required,gt=0,dive,required"`
}

// EmailSender interface
type EmailSender interface {
	DialAndSend(message ...*gomail.Message) error
}

// EmailClientConfig struct holds all required parameters for the EmailClient
type EmailClientConfig struct {
	Sender           EmailSender `validate:"required"`
	FromEmailAddress string      `validate:"required"`
	AlertSubject     string      `validate:"required"`
	ResolveSubject   string      `validate:"required"`
}

// NewEmailClient inits an EmailClient which can send e-mails for alert and resolve notifications
func NewEmailClient(config EmailClientConfig) (*EmailClient, error) {
	err := validator.New().Struct(config)
	if err != nil {
		return nil, err
	}

	return &EmailClient{
		config: config,
	}, nil
}

// EmailClient struct implements the Notifier interface and can send alert and resolve notifications
type EmailClient struct {
	config EmailClientConfig
}

// Alert sends notification to the given receiver audience.
// Email Subject will contain the noun "alert" as prefix.
func (e *EmailClient) Alert(ctx context.Context, msg Message) error {
	if err := e.sendEmailToReceivers(ctx, e.config.AlertSubject, msg); err != nil {
		return AlertError{
			err:      err,
			notifier: EmailNotifier,
		}
	}
	return nil
}

// Resolve sends notification to the given receiver audience.
// Email Subject will contain the adjective "resolved" as prefix.
func (e *EmailClient) Resolve(ctx context.Context, msg Message) error {
	if err := e.sendEmailToReceivers(ctx, e.config.ResolveSubject, msg); err != nil {
		return ResolveError{
			err:      err,
			notifier: EmailNotifier,
		}
	}
	return nil
}

func (e *EmailClient) sendEmailToReceivers(ctx context.Context, subject string, msg Message) error {
	err := validator.New().Struct(msg)
	if err != nil {
		return err
	}

	sendError := make(chan error)
	defer close(sendError)

	go func(errChan chan<- error) {
		errChan <- e.config.Sender.DialAndSend(e.buildMessages(subject, msg)...)
	}(sendError)

	select {
	case err := <-sendError:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (e *EmailClient) buildMessages(subject string, msg Message) []*gomail.Message {
	messages := make([]*gomail.Message, 0, len(msg.Receivers))
	for _, r := range msg.Receivers {
		gm := gomail.NewMessage()
		gm.SetHeader("From", e.config.FromEmailAddress)
		gm.SetHeader("To", r)
		gm.SetHeader("Subject", subject)
		gm.SetBody("text/plain", msg.Content)
		messages = append(messages, gm)
	}
	return messages
}
