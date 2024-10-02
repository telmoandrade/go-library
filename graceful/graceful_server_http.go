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
		GracefulServer
		certFile string
		keyFile  string
		attrs    []any
	}

	httpServer interface {
		ListenAndServe() error
		ListenAndServeTLS(certFile, keyFile string) error
		Shutdown(ctx context.Context) error
		Close() error
	}

	// OptionGracefulServerHttp is used to apply configurations to a [GracefulServerHttp] when creating it with [NewGracefulServerHttp].
	OptionGracefulServerHttp func(*gracefulServerHttp)
)

func gracefulServerHttpStart(gs *gracefulServerHttp, s httpServer) func() error {
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
			slog.Error(fmt.Sprintf("[HTTP SERVER] Error starting: %s", err.Error()), gs.attrs...)
			return err
		}
		return nil
	}
}

func gracefulServerHttpStop(gs *gracefulServerHttp, s httpServer) func(ctx context.Context) {
	return func(ctx context.Context) {
		slog.Info("[HTTP SERVER] Closing", gs.attrs...)
		err := s.Shutdown(ctx)
		if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			slog.Error(fmt.Sprintf("[HTTP SERVER] Error closing: %s", err.Error()), gs.attrs...)
		}
		slog.Info("[HTTP SERVER] Closed", gs.attrs...)
	}
}

func gracefulServerHttpForceStop(gs *gracefulServerHttp, s httpServer) func() {
	return func() {
		slog.Info("[HTTP SERVER] Forcing closing", gs.attrs...)
		err := s.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("[HTTP SERVER] Error forcing closing: %s", err.Error()), gs.attrs...)
		}
	}
}

// NewGracefulServerHttp returns a new [GracefulServer] encapsulating an [http.Server].
// This allows the HTTP server to be managed within the graceful shutdown framework.
// A variadic set of options to configure the behavior of the HTTP server.
func NewGracefulServerHttp(s *http.Server, opts ...OptionGracefulServerHttp) GracefulServer {
	if s == nil {
		return nil
	}

	gs := &gracefulServerHttp{}

	gs.GracefulServer = NewGracefulServer(
		WithStart(gracefulServerHttpStart(gs, s)),
		WithStop(gracefulServerHttpStop(gs, s)),
		WithForceStop(gracefulServerHttpForceStop(gs, s)),
	)

	for _, opt := range opts {
		opt(gs)
	}

	return gs
}

// WithTLS is an [OptionGracefulServerHttp] that configures TLS (Transport Layer Security) for the HTTP server.
// This option allows you to specify the certificate and private key files needed to secure the HTTP server.
func WithTLS(certFile, keyFile string) OptionGracefulServerHttp {
	return func(gs *gracefulServerHttp) {
		if certFile != "" && keyFile != "" {
			gs.certFile = certFile
			gs.keyFile = keyFile
		}
	}
}

// WithSlogAttrs is an [OptionGracefulServerHttp] that allows you to add a variadic list of [slog.Attr] to the log
// handler used by the graceful shutdown server.
// This can be useful for enhancing log output with structured attributes during the start, stop or force stop process of the http server.
func WithSlogAttrs(attrs ...slog.Attr) OptionGracefulServerHttp {
	return func(gs *gracefulServerHttp) {
		for _, a := range attrs {
			gs.attrs = append(gs.attrs, a)
		}
	}
}
