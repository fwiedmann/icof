package notifier

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"html/template"
	"strings"
)

type notifierKind int

func (n notifierKind) String() string {
	return []string{"alert", "resolve"}[n]
}

const (
	alert notifierKind = iota
	resolve
)

// EmailReceiver struct
type EmailReceiver struct {
	Name                   string    `json:"name" validate:"required"`
	AlertSubject           string    `json:"alert_subject" validate:"required"`
	ResolveSubject         string    `json:"resolve_subject" validate:"required"`
	AlertTemplateMessage   string    `json:"alert_template_message" validate:"required"`
	ResolveTemplateMessage string    `json:"resolve_template_message" validate:"required"`
	Addresses              []Address `json:"addresses" validate:"required,dive,required"`
}

// Address struct will be used to send the notification to the given Address.Email and also as the input for the message template
type Address struct {
	Email   string `json:"email" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Surname string `json:"surname" validate:"required"`
}

// EmailSender interface
type EmailSender interface {
	DialAndSend(message ...*gomail.Message) error
}

// EmailClientConfig struct holds all required parameters for the EmailClient
type EmailClientConfig struct {
	Sender           EmailSender `validate:"required"`
	FromEmailAddress string      `validate:"required"`
}

// NewEmailClient inits an EmailClient which can send e-mails for alert and resolve notifications
func NewEmailClient(config *EmailClientConfig, repo EmailReceiverRepository, logger *log.Logger) (*EmailClient, error) {
	err := validator.New().Struct(config)
	if err != nil {
		return nil, err
	}

	return &EmailClient{
		config: config,
		repo:   repo,
		logger: logger.WithFields(map[string]interface{}{"notifier": "email"}),
	}, nil
}

type EmailReceiverRepository interface {
	GetEmailReceivers(ctx context.Context) ([]EmailReceiver, error)
}

// EmailClient struct implements the Notifier interface and can send alert and resolve notifications
type EmailClient struct {
	config *EmailClientConfig
	repo   EmailReceiverRepository
	logger *log.Entry
}

// Name gives the name of the email notifier
func (e *EmailClient) Name() string {
	return "email-notifier"
}

// Alert sends notification to the given receiver audience.
// Email Subject will contain the noun "alert" as prefix.
func (e *EmailClient) Alert(ctx context.Context) error {
	e.logger.Infof("starting to send alert for email receivers")
	if err := e.sendEmailToReceivers(ctx, alert); err != nil {
		return AlertError{
			err:      err,
			notifier: EmailNotifier,
		}
	}
	e.logger.Infof("successfully send alert emails for receivers")
	return nil
}

// Resolve sends notification to the given receiver audience.
// Email Subject will contain the adjective "resolved" as prefix.
func (e *EmailClient) Resolve(ctx context.Context) error {
	e.logger.Infof("starting to send resolve for email receivers")
	if err := e.sendEmailToReceivers(ctx, resolve); err != nil {
		return ResolveError{
			err:      err,
			notifier: EmailNotifier,
		}
	}
	e.logger.Infof("successfully send resolve emails for receivers")
	return nil
}

func (e *EmailClient) sendEmailToReceivers(ctx context.Context, kind notifierKind) error {
	sendError := make(chan error)
	defer close(sendError)

	messages, messageBuildErr := e.buildMessages(kind)
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
			if err == nil {
				e.logger.Infof("successfully send all emails for notification kind %q", kind)
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

func (e *EmailClient) buildMessages(kind notifierKind) ([]*gomail.Message, error) {
	messages := make([]*gomail.Message, 0)
	failedAddressTemplates := make([]failedMessageBuild, 0)

	receivers, err := e.repo.GetEmailReceivers(context.Background())
	if err != nil {
		return nil, err
	}

	for _, receiver := range receivers {
		var tmpl *template.Template
		var subject string
		switch kind {
		case alert:
			subject = receiver.AlertSubject
			tmpl, err = template.New("alert-template").Parse(receiver.AlertTemplateMessage)
			if err != nil {
				return nil, err
			}
		case resolve:
			subject = receiver.ResolveSubject
			tmpl, err = template.New("alert-template").Parse(receiver.ResolveTemplateMessage)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("could not create messages for notifier kind %q", kind)
		}

		for _, address := range receiver.Addresses {
			buf := &strings.Builder{}
			err := tmpl.Execute(buf, address)
			if err != nil || buf.Len() == 0 {
				e.logger.Errorf("could not create template for notification kind %q and address %+v", kind, address)

				failedAddressTemplates = append(failedAddressTemplates, failedMessageBuild{
					address:      address.Email,
					receiverName: receiver.Name,
					templateName: tmpl.Name(),
					errorMessage: err.Error(),
				})
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
