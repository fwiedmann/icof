package notifier

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/gomail.v2"
	"html/template"
	"strings"
)

type notifierKind int

const (
	alert notifierKind = iota
	resolve
)

// Receiver struct
type Receiver struct {
	Name                   string    `validate:"required"`
	AlertTemplateMessage   string    `validate:"required"`
	ResolveTemplateMessage string    `validate:"required"`
	Addresses              []Address `validate:"required,dive,required"`
	parsedAlertTemplate    *template.Template
	parsedResolveTemplate  *template.Template
}

// Address struct will be used to send the notification to the given Address.Email and also as the input for the message template
type Address struct {
	Email   string `validate:"required"`
	Name    string `validate:"required"`
	Surname string `validate:"required"`
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
	Receivers        []Receiver  `validate:"required,dive,required"`
}

// NewEmailClient inits an EmailClient which can send e-mails for alert and resolve notifications
func NewEmailClient(config *EmailClientConfig) (*EmailClient, error) {
	err := validator.New().Struct(config)
	if err != nil {
		return nil, err
	}

	for i, receiver := range config.Receivers {
		tmpl, err := template.New("alert-template").Parse(receiver.AlertTemplateMessage)
		if err != nil {
			return nil, err
		}
		config.Receivers[i].parsedAlertTemplate = tmpl

		tmpl, err = template.New("resolve-template").Parse(receiver.ResolveTemplateMessage)
		if err != nil {
			return nil, err
		}
		config.Receivers[i].parsedResolveTemplate = tmpl
	}

	return &EmailClient{
		config: config,
	}, nil

}

// EmailClient struct implements the Notifier interface and can send alert and resolve notifications
type EmailClient struct {
	config *EmailClientConfig
}

// Alert sends notification to the given receiver audience.
// Email Subject will contain the noun "alert" as prefix.
func (e *EmailClient) Alert(ctx context.Context) error {
	if err := e.sendEmailToReceivers(ctx, e.config.AlertSubject, alert); err != nil {
		return AlertError{
			err:      err,
			notifier: EmailNotifier,
		}
	}
	return nil
}

// Resolve sends notification to the given receiver audience.
// Email Subject will contain the adjective "resolved" as prefix.
func (e *EmailClient) Resolve(ctx context.Context) error {
	if err := e.sendEmailToReceivers(ctx, e.config.ResolveSubject, resolve); err != nil {
		return ResolveError{
			err:      err,
			notifier: EmailNotifier,
		}
	}
	return nil
}

func (e *EmailClient) sendEmailToReceivers(ctx context.Context, subject string, kind notifierKind) error {
	sendError := make(chan error)
	defer close(sendError)

	messages, messageBuildErr := e.buildMessages(subject, kind)
	// if no messages could be created, no email should be sent
	if messageBuildErr != nil && len(messages) == 0 {
		return messageBuildErr
	}

	go func(errChan chan<- error) {
		errChan <- e.config.Sender.DialAndSend(messages...)
	}(sendError)

	for {
		select {
		case err := <-sendError:
			// when an error occurred during message creation but also had some emails which should be sent and a send error occurred, user should be notified
			if err != nil && messageBuildErr != nil {
				return fmt.Errorf("could not create all messages and could not send all emails: message error: %w, email send error, %s", messageBuildErr, err)
			}
			if messageBuildErr != nil {
				return messageBuildErr
			}
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

type failedMessageBuild struct {
	templateName string
	address      string
	receiverName string
	errorMessage string
}

func (e *EmailClient) buildMessages(subject string, kind notifierKind) ([]*gomail.Message, error) {
	messages := make([]*gomail.Message, 0)
	failedAddressTemplates := make([]failedMessageBuild, 0)

	for _, receiver := range e.config.Receivers {
		for _, address := range receiver.Addresses {
			buf := &strings.Builder{}

			templateFunc := func(tmpl *template.Template) {
				err := tmpl.Execute(buf, address)
				if err != nil || buf.Len() == 0 {
					failedAddressTemplates = append(failedAddressTemplates, failedMessageBuild{
						address:      address.Email,
						receiverName: receiver.Name,
						templateName: tmpl.Name(),
					})
				}
			}

			if kind == alert {
				templateFunc(receiver.parsedAlertTemplate)
			}

			if kind == resolve {
				templateFunc(receiver.parsedAlertTemplate)
			}

			gm := gomail.NewMessage()
			gm.SetHeader("From", e.config.FromEmailAddress)
			gm.SetHeader("To", address.Email)
			gm.SetHeader("Subject", subject)
			gm.SetBody("text/plain", buf.String())
			messages = append(messages, gm)
		}
	}

	if len(failedAddressTemplates) != 0 {
		return messages, EmailMessageBuildError{failedMessages: failedAddressTemplates}
	}
	return messages, nil
}

// EmailMessageBuildError represents all messages that could not be built with the template and the given address input
type EmailMessageBuildError struct {
	failedMessages []failedMessageBuild
}

// Error formats the EmailMessageBuildError to a nice error message
func (e EmailMessageBuildError) Error() string {
	return fmt.Sprintf("could not build messages for: %+v", e.failedMessages)
}
