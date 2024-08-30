//go:generate mockgen -source=graceful_server_http.go -destination=mock_graceful_server_http_test.go -package graceful

package graceful

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type (
	gracefulServerHttp struct {
		*gracefulServer
		certFile string
		keyFile  string
		attrs    []any
	}

	// HttpServer http.Server wrapper
	HttpServer interface {
		// ListenAndServe http.Server.ListenAndServe()
		ListenAndServe() error
		// ListenAndServeTLS http.Server.ListenAndServeTLS()
		ListenAndServeTLS(certFile, keyFile string) error
		// Shutdown http.Server.Shutdown()
		Shutdown(ctx context.Context) error
		// Close http.Server.Close()
		Close() error
	}

	// OptionGracefulServerHttp used in [NewGracefulServerHttp]
	OptionGracefulServerHttp func(*gracefulServerHttp)
)

func gracefulServerHttpStart(gs *gracefulServerHttp, s HttpServer) func() error {
	return func() error {
		var err error
		if gs.certFile != "" {
			slog.Info("[HTTP SERVER] Starting with TLS", slog.Int("port", 8080))
			err = s.ListenAndServeTLS(gs.certFile, gs.keyFile)
		} else {
			slog.Info("[HTTP SERVER] Starting", gs.attrs...)
			err = s.ListenAndServe()
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("[HTTP SERVER] %s", err.Error()), gs.attrs...)
			return err
		}
		return nil
	}
}

func gracefulServerHttpStop(gs *gracefulServerHttp, s HttpServer) func(ctx context.Context) {
	return func(ctx context.Context) {
		slog.Info("[HTTP SERVER] Closing", gs.attrs...)
		err := s.Shutdown(ctx)
		if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			slog.Error(fmt.Sprintf("[HTTP SERVER] %s", err.Error()), gs.attrs...)
		}
		slog.Info("[HTTP SERVER] Closed", gs.attrs...)
	}
}

func gracefulServerHttpForceStop(gs *gracefulServerHttp, s HttpServer) func() {
	return func() {
		slog.Info("[HTTP SERVER] Force close", gs.attrs...)
		err := s.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("[HTTP SERVER] %s", err.Error()), gs.attrs...)
		}
	}
}

// NewGracefulServerHttp create a new graceful server for http.Server
func NewGracefulServerHttp(s HttpServer, opts ...OptionGracefulServerHttp) GracefulServer {
	if s == nil {
		return nil
	}

	gs := &gracefulServerHttp{}

	gs.gracefulServer = &gracefulServer{
		start:     gracefulServerHttpStart(gs, s),
		stop:      gracefulServerHttpStop(gs, s),
		forceStop: gracefulServerHttpForceStop(gs, s),
	}

	for _, opt := range opts {
		opt(gs)
	}

	return gs
}

// WithTLS configure the filenames containing a certificate and matching private key. Look http.Server for more details
func WithTLS(certFile, keyFile string) OptionGracefulServerHttp {
	return func(gs *gracefulServerHttp) {
		if certFile != "" && keyFile != "" {
			gs.certFile = certFile
			gs.keyFile = keyFile
		}
	}
}

// WithSlogAttrs add information to graceful server log when starting and stopping
func WithSlogAttrs(attrs ...slog.Attr) OptionGracefulServerHttp {
	return func(gs *gracefulServerHttp) {
		for _, a := range attrs {
			gs.attrs = append(gs.attrs, a)
		}
	}
}
