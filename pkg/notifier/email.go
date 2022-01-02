package notifier

import (
	"context"
	"github.com/go-playground/validator/v10"
	"net/smtp"
)

// Message struct will be used to configure the send operation
type Message struct {
	Content   string   `validate:"required"`
	Receivers []string `validate:"required,gt=0,dive,required"`
}

// EmailClientConfig struct holds all required parameters for the EmailClient
type EmailClientConfig struct {
	Username         string `validate:"required"`
	Password         string `validate:"required"`
	SmtpHost         string `validate:"required"`
	FromEmailAddress string `validate:"required"`
}

// NewEmailClient inits an EmailClient which can send e-mails for alert and resolve notifications
func NewEmailClient(config EmailClientConfig) (*EmailClient, error) {
	err := validator.New().Struct(config)
	if err != nil {
		return nil, err
	}

	return &EmailClient{
		auth:             smtp.PlainAuth("", config.Username, config.Password, config.SmtpHost),
		smtpHost:         config.SmtpHost,
		fromEmailAddress: config.FromEmailAddress,
	}, nil
}

// EmailClient struct implements the Notifier interface and can send alert and resolve notifications
type EmailClient struct {
	auth             smtp.Auth
	smtpHost         string
	fromEmailAddress string
}

// Alert sends notification to the given receiver audience.
// Email Subject will contain the noun "alert" as prefix.
func (e *EmailClient) Alert(ctx context.Context, msg Message) error {
	if err := e.sendEmail(ctx, msg); err != nil {
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
	if err := e.sendEmail(ctx, msg); err != nil {
		return ResolveError{
			err:      err,
			notifier: EmailNotifier,
		}
	}
	return nil
}

func (e *EmailClient) sendEmail(ctx context.Context, msg Message) error {
	err := validator.New().Struct(msg)
	if err != nil {
		return err
	}

	sendError := make(chan error)
	defer close(sendError)

	go func(errChan chan<- error) {
		// Todo: Implement EmailMessage struct to be RFC 822-style complaint or with https://github.com/go-gomail/gomail
		errChan <- smtp.SendMail(e.smtpHost, e.auth, e.fromEmailAddress, msg.Receivers, []byte(msg.Content))
	}(sendError)

	select {
	case err := <-sendError:
		return err
	case <-ctx.Done():
		return nil
	}
}

type EmailMessage struct {
	subject string
	to      string
	from    string
	message string
}
