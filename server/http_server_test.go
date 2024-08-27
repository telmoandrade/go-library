//go:generate mockgen -source=http_server.go -destination=mock_http_server_test.go -package server

package server

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	gomock "go.uber.org/mock/gomock"
)

func Test_httpStop(t *testing.T) {
	slog.SetLogLoggerLevel(slog.Level(16))

	t.Run("no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockHttpServer(ctrl)
		mock.EXPECT().Shutdown(gomock.Any()).Return(nil).Times(1)

		httpStop(mock)(context.Background())
	})

	t.Run("with error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockHttpServer(ctrl)
		mock.EXPECT().Shutdown(gomock.Any()).Return(errors.New("error")).Times(1)

		httpStop(mock)(context.Background())
	})
}

func Test_httpForceStop(t *testing.T) {
	slog.SetLogLoggerLevel(slog.Level(16))

	t.Run("no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockHttpServer(ctrl)
		mock.EXPECT().Close().Return(nil).Times(1)

		httpForceStop(mock)()
	})

	t.Run("with error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockHttpServer(ctrl)
		mock.EXPECT().Close().Return(errors.New("error")).Times(1)

		httpForceStop(mock)()
	})
}

func TestNewServerHttp(t *testing.T) {
	slog.SetLogLoggerLevel(slog.Level(16))

	t.Run("no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockHttpServer(ctrl)
		mock.EXPECT().ListenAndServe().Return(nil).Times(1)

		if got := NewServerHttp(mock).Start(); got != nil {
			t.Errorf("Start() = %v, want %v", got, nil)
		}
	})

	t.Run("with error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockHttpServer(ctrl)
		mock.EXPECT().ListenAndServe().Return(errors.New("error")).Times(1)

		if got := NewServerHttp(mock).Start(); got == nil {
			t.Errorf("Start() = %v, want %v", got, "error")
		}
	})
}

func TestNewServerHttpWithTLS(t *testing.T) {
	slog.SetLogLoggerLevel(slog.Level(16))

	t.Run("no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockHttpServer(ctrl)
		mock.EXPECT().ListenAndServeTLS(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		if got := NewServerHttpWithTLS(mock, "certFile", "keyFile").Start(); got != nil {
			t.Errorf("Start() = %v, want %v", got, nil)
		}
	})

	t.Run("with error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockHttpServer(ctrl)
		mock.EXPECT().ListenAndServeTLS(gomock.Any(), gomock.Any()).Return(errors.New("error")).Times(1)

		if got := NewServerHttpWithTLS(mock, "certFile", "keyFile").Start(); got == nil {
			t.Errorf("Start() = %v, want %v", got, "error")
		}
	})
}
