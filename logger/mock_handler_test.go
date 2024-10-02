package logger

import (
	"context"
	"log/slog"
	"reflect"

	"go.uber.org/mock/gomock"
)

type (
	MockSlogHandler struct {
		ctrl     *gomock.Controller
		recorder *MockSlogHandlerRecord
	}

	MockSlogHandlerRecord struct {
		mock *MockSlogHandler
	}
)

func NewMockSlogHandler(ctrl *gomock.Controller) *MockSlogHandler {
	mock := &MockSlogHandler{ctrl: ctrl}
	mock.recorder = &MockSlogHandlerRecord{mock}
	return mock
}

func (m *MockSlogHandler) EXPECT() *MockSlogHandlerRecord {
	return m.recorder
}

func (m *MockSlogHandler) Enabled(ctx context.Context, l slog.Level) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Enabled", ctx, l)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (mr *MockSlogHandlerRecord) Enabled(ctx context.Context, l slog.Level) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enabled", reflect.TypeOf((*MockSlogHandler)(nil).Enabled), ctx, l)
}

func (m *MockSlogHandler) Handle(ctx context.Context, r slog.Record) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, r)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockSlogHandlerRecord) Handle(ctx any, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockSlogHandler)(nil).Handle), ctx, r)
}

func (m *MockSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithAttrs", attrs)
	ret0, _ := ret[0].(slog.Handler)
	return ret0
}

func (mr *MockSlogHandlerRecord) WithAttrs(attrs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithAttrs", reflect.TypeOf((*MockSlogHandler)(nil).WithAttrs), attrs)
}

func (m *MockSlogHandler) WithGroup(name string) slog.Handler {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithGroup", name)
	ret0, _ := ret[0].(slog.Handler)
	return ret0
}

func (mr *MockSlogHandlerRecord) WithGroup(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithGroup", reflect.TypeOf((*MockSlogHandler)(nil).WithGroup), name)
}
