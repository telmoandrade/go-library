//go:generate mockgen -source=http_server.go -destination=mock_http_server_test.go -package server

package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type (
	HttpServer interface {
		Shutdown(ctx context.Context) error
		Close() error
		ListenAndServe() error
		ListenAndServeTLS(certFile, keyFile string) error
	}
)

func httpStop(s HttpServer) func(context.Context) {
	return func(ctx context.Context) {
		slog.Info("[HTTP SERVER] Closing")
		err := s.Shutdown(ctx)
		if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			slog.Error(fmt.Sprintf("[HTTP SERVER] %s", err.Error()))
		}
		slog.Info("[HTTP SERVER] Closed")
	}
}

func httpForceStop(s HttpServer) func() {
	return func() {
		slog.Info("[HTTP SERVER] Force close")
		err := s.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("[HTTP SERVER] %s", err.Error()))
		}
	}
}

func NewServerHttp(s HttpServer) Server {
	return NewServer(
		WithStart(func() error {
			slog.Info("[HTTP SERVER] Starting")
			err := s.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error(fmt.Sprintf("[HTTP SERVER] %s", err.Error()))
				return err
			}
			return nil
		}),
		WithStop(httpStop(s)),
		WithForceStop(httpForceStop(s)),
	)
}

func NewServerHttpWithTLS(s HttpServer, certFile, keyFile string) Server {
	return NewServer(
		WithStart(func() error {
			slog.Info("[HTTP SERVER] Starting")
			err := s.ListenAndServeTLS(certFile, keyFile)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error(fmt.Sprintf("[HTTP SERVER] %s", err.Error()))
				return err
			}
			return nil
		}),
		WithStop(httpStop(s)),
		WithForceStop(httpForceStop(s)),
	)
}
