package mock_email

import (
	"gopkg.in/gomail.v2"
	"testing"
)

func NewMockEmailSender(t *testing.T, msgCount int, returnErr error) MockEmailSender {
	return MockEmailSender{
		t:            t,
		messageCount: msgCount,
		returnErr:    returnErr,
	}
}

type MockEmailSender struct {
	t            *testing.T
	messageCount int
	returnErr    error
}

func (m MockEmailSender) DialAndSend(message ...*gomail.Message) error {

	if len(message) != m.messageCount {
		m.t.Fatalf("Inalid input. Given messages have a size of %d, want %d", len(message), m.messageCount)
	}
	return m.returnErr
}
