package notifier_test

import (
	"context"
	"errors"
	"github.com/fwiedmann/icof/pkg/notifier"
	mock_email "github.com/fwiedmann/icof/pkg/notifier/mock-email"
	"testing"
)

func TestEmailClient_Alert(t *testing.T) {
	t.Parallel()
	type inputClient struct {
		config notifier.EmailClientConfig
	}

	type input struct {
		message notifier.Message
	}

	tests := []struct {
		name                 string
		input                input
		inputClient          inputClient
		wantError            bool
		mockError            error
		mockWantMessageCount int
	}{
		{
			name: "should_send_email_to_all_receivers",
			input: input{
				message: notifier.Message{
					Content:   "ALARM!!!",
					Receivers: []string{"receiver@example.com"},
				},
			},
			inputClient: inputClient{
				config: notifier.EmailClientConfig{
					FromEmailAddress: "example@example.com",
					AlertSubject:     "Alert!",
					ResolveSubject:   "Resolved!",
				},
			},
			wantError:            false,
			mockWantMessageCount: 1,
			mockError:            nil,
		},
		{
			name: "should_return_alert_error",
			input: input{
				message: notifier.Message{
					Content:   "ALARM!!!",
					Receivers: []string{"receiver@example.com"},
				},
			},
			inputClient: inputClient{
				config: notifier.EmailClientConfig{
					FromEmailAddress: "example@example.com",
					AlertSubject:     "Alert!",
					ResolveSubject:   "Resolved!",
				},
			},
			wantError:            true,
			mockError:            errors.New("send error"),
			mockWantMessageCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.inputClient.config.Sender = mock_email.NewMockEmailSender(t, tt.mockWantMessageCount, tt.mockError)
			c, err := notifier.NewEmailClient(&tt.inputClient.config)
			if (err != nil) && !tt.wantError {
				t.Fatalf("NewEmailClient() retuned an error but did not want one: %s", err)
			}

			err = c.Alert(context.Background(), tt.input.message)
			if (err != nil) && !tt.wantError {
				t.Fatalf("Alert() retuned an error but did not want one: %s", err)
			}

			if (err != nil) && tt.wantError && !errors.As(err, &notifier.AlertError{}) {
				t.Fatalf("Alert() retuned an error with wrong error type. Want an notifier.AlertError error")
			}
		})
	}
}

func TestEmailClient_Resolve(t *testing.T) {
	t.Parallel()
	type inputClient struct {
		config notifier.EmailClientConfig
	}

	type input struct {
		message notifier.Message
	}

	tests := []struct {
		name                 string
		input                input
		inputClient          inputClient
		wantError            bool
		mockError            error
		mockWantMessageCount int
	}{
		{
			name: "should_send_email_to_all_receivers",
			input: input{
				message: notifier.Message{
					Content:   "Resolve",
					Receivers: []string{"receiver@example.com"},
				},
			},
			inputClient: inputClient{
				config: notifier.EmailClientConfig{
					FromEmailAddress: "example@example.com",
					AlertSubject:     "Alert!",
					ResolveSubject:   "Resolved!",
				},
			},
			mockWantMessageCount: 1,
			wantError:            false,
		},
		{
			name: "should_return_alert_error",
			input: input{
				message: notifier.Message{
					Content:   "ALARM!!!",
					Receivers: []string{"receiver@example.com"},
				},
			},
			inputClient: inputClient{
				config: notifier.EmailClientConfig{
					FromEmailAddress: "example@example.com",
					AlertSubject:     "Alert!",
					ResolveSubject:   "Resolved!",
				},
			},
			wantError:            true,
			mockError:            errors.New("send error"),
			mockWantMessageCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.inputClient.config.Sender = mock_email.NewMockEmailSender(t, tt.mockWantMessageCount, tt.mockError)
			c, err := notifier.NewEmailClient(&tt.inputClient.config)
			if (err != nil) && !tt.wantError {
				t.Fatalf("NewEmailClient() retuned an error but did not want one: %s", err)
			}

			err = c.Resolve(context.Background(), tt.input.message)
			if (err != nil) && !tt.wantError {
				t.Fatalf("Alert() retuned an error but did not want one: %s", err)
			}

			if (err != nil) && tt.wantError && !errors.As(err, &notifier.ResolveError{}) {
				t.Fatalf("Alert() retuned an error with wrong error type. Want an notifier.ResolveError error")
			}
		})
	}
}
