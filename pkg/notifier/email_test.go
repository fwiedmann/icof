package notifier_test

import (
	"context"
	"errors"
	"github.com/fwiedmann/icof/pkg/notifier"
	mock_email "github.com/fwiedmann/icof/pkg/notifier/mock-email"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestEmailClient_Alert(t *testing.T) {
	t.Parallel()
	type inputClient struct {
		config notifier.EmailClientConfig
	}

	type repoResponse struct {
		receivers []notifier.EmailReceiver
		err       error
	}

	tests := []struct {
		name                 string
		inputClient          inputClient
		repoResponse         repoResponse
		wantError            bool
		mockError            error
		mockWantMessageCount int
	}{
		{
			name: "should_send_email_to_all_receivers",
			inputClient: inputClient{
				config: notifier.EmailClientConfig{
					FromEmailAddress: "example@example.com",
				},
			},
			repoResponse: repoResponse{
				receivers: []notifier.EmailReceiver{
					{
						Name:                   "colleagues",
						AlertTemplateMessage:   "Alert occurred, I'm AFK",
						ResolveTemplateMessage: "Back in buiss",
						AlertSubject:           "Alert!",
						ResolveSubject:         "Resolved!",
						Addresses: []notifier.Address{
							{
								Email:   "example@example.com",
								Name:    "Andi",
								Surname: "Developer",
							},
						},
					},
				},
				err: nil,
			},
			wantError:            false,
			mockWantMessageCount: 1,
			mockError:            nil,
		},
		{
			name: "should_return_alert_error",
			inputClient: inputClient{
				config: notifier.EmailClientConfig{
					FromEmailAddress: "example@example.com",
				},
			},
			repoResponse: repoResponse{
				receivers: []notifier.EmailReceiver{
					{
						Name:                   "colleagues",
						AlertTemplateMessage:   "Alert occurred, I'm AFK",
						ResolveTemplateMessage: "Back in buiss",
						AlertSubject:           "Alert!",
						ResolveSubject:         "Resolved!",
						Addresses: []notifier.Address{
							{
								Email:   "example@example.com",
								Name:    "Andi",
								Surname: "Developer",
							},
						},
					},
				},
				err: nil,
			},
			wantError:            true,
			mockError:            errors.New("send error"),
			mockWantMessageCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			repoMock := mock_email.NewMockEmailReceiverRepository(mockCtrl)

			repoMock.EXPECT().GetReceivers(gomock.Any()).Return(tt.repoResponse.receivers, tt.repoResponse.err)
			tt.inputClient.config.Sender = mock_email.NewMockEmailSender(t, tt.mockWantMessageCount, tt.mockError)
			c, err := notifier.NewEmailClient(&tt.inputClient.config, repoMock)
			if err != nil {
				t.Fatalf("NewEmailClient() retuned an error but did not want one: %s", err)
			}

			err = c.Alert(context.Background())
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

	type repoResponse struct {
		receivers []notifier.EmailReceiver
		err       error
	}

	tests := []struct {
		name                 string
		inputClient          inputClient
		repoResponse         repoResponse
		wantError            bool
		mockError            error
		mockWantMessageCount int
	}{
		{
			name: "should_send_email_to_all_receivers",
			inputClient: inputClient{
				config: notifier.EmailClientConfig{
					FromEmailAddress: "example@example.com",
				},
			},
			repoResponse: repoResponse{
				receivers: []notifier.EmailReceiver{
					{
						Name:                   "colleagues",
						AlertTemplateMessage:   "Alert occurred, I'm AFK",
						ResolveTemplateMessage: "Back in buiss",
						AlertSubject:           "Alert!",
						ResolveSubject:         "Resolved!",
						Addresses: []notifier.Address{
							{
								Email:   "example@example.com",
								Name:    "Andi",
								Surname: "Developer",
							},
						},
					},
				},
				err: nil,
			},
			mockWantMessageCount: 1,
			wantError:            false,
		},
		{
			name: "should_return_alert_error",
			inputClient: inputClient{
				config: notifier.EmailClientConfig{
					FromEmailAddress: "example@example.com",
				},
			},
			repoResponse: repoResponse{
				receivers: []notifier.EmailReceiver{
					{
						Name:                   "colleagues",
						AlertTemplateMessage:   "Alert occurred, I'm AFK",
						ResolveTemplateMessage: "Back in buiss",
						AlertSubject:           "Alert!",
						ResolveSubject:         "Resolved!",
						Addresses: []notifier.Address{
							{
								Email:   "example@example.com",
								Name:    "Andi",
								Surname: "Developer",
							},
						},
					},
				},
				err: nil,
			},
			wantError:            true,
			mockError:            errors.New("send error"),
			mockWantMessageCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			repoMock := mock_email.NewMockEmailReceiverRepository(mockCtrl)

			repoMock.EXPECT().GetReceivers(gomock.Any()).Return(tt.repoResponse.receivers, tt.repoResponse.err)
			tt.inputClient.config.Sender = mock_email.NewMockEmailSender(t, tt.mockWantMessageCount, tt.mockError)
			c, err := notifier.NewEmailClient(&tt.inputClient.config, repoMock)
			if (err != nil) && !tt.wantError {
				t.Fatalf("NewEmailClient() retuned an error but did not want one: %s", err)
			}

			err = c.Resolve(context.Background())
			if (err != nil) && !tt.wantError {
				t.Fatalf("Alert() retuned an error but did not want one: %s", err)
			}

			if (err != nil) && tt.wantError && !errors.As(err, &notifier.ResolveError{}) {
				t.Fatalf("Alert() retuned an error with wrong error type. Want an notifier.ResolveError error")
			}
		})
	}
}
