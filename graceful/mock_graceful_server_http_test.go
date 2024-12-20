// Code generated by MockGen. DO NOT EDIT.
// Source: graceful_server_http.go
//
// Generated by this command:
//
//	mockgen -source=graceful_server_http.go -destination=mock_graceful_server_http_test.go -package graceful
//

// Package graceful is a generated GoMock package.
package graceful

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockhttpServer is a mock of httpServer interface.
type MockhttpServer struct {
	ctrl     *gomock.Controller
	recorder *MockhttpServerMockRecorder
	isgomock struct{}
}

// MockhttpServerMockRecorder is the mock recorder for MockhttpServer.
type MockhttpServerMockRecorder struct {
	mock *MockhttpServer
}

// NewMockhttpServer creates a new mock instance.
func NewMockhttpServer(ctrl *gomock.Controller) *MockhttpServer {
	mock := &MockhttpServer{ctrl: ctrl}
	mock.recorder = &MockhttpServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockhttpServer) EXPECT() *MockhttpServerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockhttpServer) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockhttpServerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockhttpServer)(nil).Close))
}

// ListenAndServe mocks base method.
func (m *MockhttpServer) ListenAndServe() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListenAndServe")
	ret0, _ := ret[0].(error)
	return ret0
}

// ListenAndServe indicates an expected call of ListenAndServe.
func (mr *MockhttpServerMockRecorder) ListenAndServe() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenAndServe", reflect.TypeOf((*MockhttpServer)(nil).ListenAndServe))
}

// ListenAndServeTLS mocks base method.
func (m *MockhttpServer) ListenAndServeTLS(certFile, keyFile string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListenAndServeTLS", certFile, keyFile)
	ret0, _ := ret[0].(error)
	return ret0
}

// ListenAndServeTLS indicates an expected call of ListenAndServeTLS.
func (mr *MockhttpServerMockRecorder) ListenAndServeTLS(certFile, keyFile any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenAndServeTLS", reflect.TypeOf((*MockhttpServer)(nil).ListenAndServeTLS), certFile, keyFile)
}

// Shutdown mocks base method.
func (m *MockhttpServer) Shutdown(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockhttpServerMockRecorder) Shutdown(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockhttpServer)(nil).Shutdown), ctx)
}
