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

// MockHttpServer is a mock of HttpServer interface.
type MockHttpServer struct {
	ctrl     *gomock.Controller
	recorder *MockHttpServerMockRecorder
}

// MockHttpServerMockRecorder is the mock recorder for MockHttpServer.
type MockHttpServerMockRecorder struct {
	mock *MockHttpServer
}

// NewMockHttpServer creates a new mock instance.
func NewMockHttpServer(ctrl *gomock.Controller) *MockHttpServer {
	mock := &MockHttpServer{ctrl: ctrl}
	mock.recorder = &MockHttpServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHttpServer) EXPECT() *MockHttpServerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockHttpServer) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockHttpServerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockHttpServer)(nil).Close))
}

// ListenAndServe mocks base method.
func (m *MockHttpServer) ListenAndServe() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListenAndServe")
	ret0, _ := ret[0].(error)
	return ret0
}

// ListenAndServe indicates an expected call of ListenAndServe.
func (mr *MockHttpServerMockRecorder) ListenAndServe() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenAndServe", reflect.TypeOf((*MockHttpServer)(nil).ListenAndServe))
}

// ListenAndServeTLS mocks base method.
func (m *MockHttpServer) ListenAndServeTLS(certFile, keyFile string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListenAndServeTLS", certFile, keyFile)
	ret0, _ := ret[0].(error)
	return ret0
}

// ListenAndServeTLS indicates an expected call of ListenAndServeTLS.
func (mr *MockHttpServerMockRecorder) ListenAndServeTLS(certFile, keyFile any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenAndServeTLS", reflect.TypeOf((*MockHttpServer)(nil).ListenAndServeTLS), certFile, keyFile)
}

// Shutdown mocks base method.
func (m *MockHttpServer) Shutdown(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockHttpServerMockRecorder) Shutdown(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockHttpServer)(nil).Shutdown), ctx)
}