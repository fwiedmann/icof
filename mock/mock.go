// Code generated by MockGen. DO NOT EDIT.
// Source: icof.go

// Package mock_icof is a generated GoMock package.
package mock_icof

import (
	context "context"
	reflect "reflect"

	icof "github.com/fwiedmann/icof"
	gomock "github.com/golang/mock/gomock"
)

// MockObserver is a mock of Observer interface.
type MockObserver struct {
	ctrl     *gomock.Controller
	recorder *MockObserverMockRecorder
}

// MockObserverMockRecorder is the mock recorder for MockObserver.
type MockObserverMockRecorder struct {
	mock *MockObserver
}

// NewMockObserver creates a new mock instance.
func NewMockObserver(ctrl *gomock.Controller) *MockObserver {
	mock := &MockObserver{ctrl: ctrl}
	mock.recorder = &MockObserverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockObserver) EXPECT() *MockObserverMockRecorder {
	return m.recorder
}

// Observe mocks base method.
func (m *MockObserver) Observe(arg0 context.Context, arg1 chan<- icof.ObserverState) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Observe", arg0, arg1)
}

// Observe indicates an expected call of Observe.
func (mr *MockObserverMockRecorder) Observe(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Observe", reflect.TypeOf((*MockObserver)(nil).Observe), arg0, arg1)
}

// MockNotifier is a mock of Notifier interface.
type MockNotifier struct {
	ctrl     *gomock.Controller
	recorder *MockNotifierMockRecorder
}

// MockNotifierMockRecorder is the mock recorder for MockNotifier.
type MockNotifierMockRecorder struct {
	mock *MockNotifier
}

// NewMockNotifier creates a new mock instance.
func NewMockNotifier(ctrl *gomock.Controller) *MockNotifier {
	mock := &MockNotifier{ctrl: ctrl}
	mock.recorder = &MockNotifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotifier) EXPECT() *MockNotifierMockRecorder {
	return m.recorder
}

// Alert mocks base method.
func (m *MockNotifier) Alert(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Alert", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Alert indicates an expected call of Alert.
func (mr *MockNotifierMockRecorder) Alert(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Alert", reflect.TypeOf((*MockNotifier)(nil).Alert), ctx)
}

// Name mocks base method.
func (m *MockNotifier) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockNotifierMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockNotifier)(nil).Name))
}

// Resolve mocks base method.
func (m *MockNotifier) Resolve(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resolve", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Resolve indicates an expected call of Resolve.
func (mr *MockNotifierMockRecorder) Resolve(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resolve", reflect.TypeOf((*MockNotifier)(nil).Resolve), ctx)
}

// MockStateRepository is a mock of StateRepository interface.
type MockStateRepository struct {
	ctrl     *gomock.Controller
	recorder *MockStateRepositoryMockRecorder
}

// MockStateRepositoryMockRecorder is the mock recorder for MockStateRepository.
type MockStateRepositoryMockRecorder struct {
	mock *MockStateRepository
}

// NewMockStateRepository creates a new mock instance.
func NewMockStateRepository(ctrl *gomock.Controller) *MockStateRepository {
	mock := &MockStateRepository{ctrl: ctrl}
	mock.recorder = &MockStateRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStateRepository) EXPECT() *MockStateRepositoryMockRecorder {
	return m.recorder
}

// GetLatest mocks base method.
func (m *MockStateRepository) GetLatest(ctx context.Context) (icof.ObserverState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatest", ctx)
	ret0, _ := ret[0].(icof.ObserverState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatest indicates an expected call of GetLatest.
func (mr *MockStateRepositoryMockRecorder) GetLatest(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatest", reflect.TypeOf((*MockStateRepository)(nil).GetLatest), ctx)
}

// Save mocks base method.
func (m *MockStateRepository) Save(ctx context.Context, state icof.ObserverState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, state)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockStateRepositoryMockRecorder) Save(ctx, state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockStateRepository)(nil).Save), ctx, state)
}
