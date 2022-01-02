package notifier

import (
	"fmt"
)

// Notifier enum
type Notifier int

// String will return the name of the notifier
func (n Notifier) String() string {
	return [...]string{"E-Mail"}[n]
}

const (
	// EmailNotifier enum for EmailClient
	EmailNotifier Notifier = iota
)

// AlertError for errors which occurred during the alert send operation
type AlertError struct {
	err      error
	notifier Notifier
}

// Error
func (ae AlertError) Error() string {
	return fmt.Sprintf("could not send alert message with \"%s\" notifier: %s", ae.notifier, ae.err)
}

// Unwrap can be used by the errors package to access the underlying error
func (ae AlertError) Unwrap() error {
	return ae.err
}

// ResolveError for errors which occurred during the alert send operation
type ResolveError struct {
	err      error
	notifier Notifier
}

// Error
func (re ResolveError) Error() string {
	return fmt.Sprintf("could not send resolve message with \"%s\" notifier: %s", re.notifier, re.err)
}

// Unwrap can be used by the errors package to access the underlying error
func (re ResolveError) Unwrap() error {
	return re.err
}
