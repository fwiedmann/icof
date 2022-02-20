// Code generated by MockGen. DO NOT EDIT.
// Source: email.go

// Package mock_notifier is a generated GoMock package.
package mock_email

import (
	context "context"
	reflect "reflect"

	notifier "github.com/fwiedmann/icof/pkg/notifier"
	gomock "github.com/golang/mock/gomock"
)

// MockEmailReceiverRepository is a mock of EmailReceiverRepository interface.
type MockEmailReceiverRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEmailReceiverRepositoryMockRecorder
}

// MockEmailReceiverRepositoryMockRecorder is the mock recorder for MockEmailReceiverRepository.
type MockEmailReceiverRepositoryMockRecorder struct {
	mock *MockEmailReceiverRepository
}

// NewMockEmailReceiverRepository creates a new mock instance.
func NewMockEmailReceiverRepository(ctrl *gomock.Controller) *MockEmailReceiverRepository {
	mock := &MockEmailReceiverRepository{ctrl: ctrl}
	mock.recorder = &MockEmailReceiverRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailReceiverRepository) EXPECT() *MockEmailReceiverRepositoryMockRecorder {
	return m.recorder
}

// GetReceivers mocks base method.
func (m *MockEmailReceiverRepository) GetEmailReceivers(ctx context.Context) ([]notifier.EmailReceiver, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmailReceivers", ctx)
	ret0, _ := ret[0].([]notifier.EmailReceiver)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReceivers indicates an expected call of GetReceivers.
func (mr *MockEmailReceiverRepositoryMockRecorder) GetReceivers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmailReceivers", reflect.TypeOf((*MockEmailReceiverRepository)(nil).GetEmailReceivers), ctx)
}