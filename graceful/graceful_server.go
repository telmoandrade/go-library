//go:generate mockgen -source=graceful_server.go -destination=mock_graceful_graceful_server_test.go -package graceful

package graceful

import (
	"context"
)

type (
	gracefulServer struct {
		start     func() error
		stop      func(context.Context)
		forceStop func()
	}

	// GracefulServer used in [WithServers]
	GracefulServer interface {
		Start() error
		Stop(context.Context)
		ForceStop()
	}

	// OptionGracefulServer used in [NewGracefulServer]
	OptionGracefulServer func(*gracefulServer)
)

// NewGracefulServer [GracefulServer] wrapper, create a new server for graceful shutdown
func NewGracefulServer(opts ...OptionGracefulServer) GracefulServer {
	gs := &gracefulServer{
		start:     func() error { return nil },
		stop:      func(context.Context) {},
		forceStop: func() {},
	}

	for _, opt := range opts {
		opt(gs)
	}

	return gs
}

// WithStart used to define a function to start the server
func WithStart(f func() error) OptionGracefulServer {
	return func(gs *gracefulServer) {
		if f != nil {
			gs.start = f
		}
	}
}

// WithStop used to define a function to stop the server
func WithStop(f func(context.Context)) OptionGracefulServer {
	return func(gs *gracefulServer) {
		if f != nil {
			gs.stop = f
		}
	}
}

// WithForceStop used to define a function to forcibly stop the server
func WithForceStop(f func()) OptionGracefulServer {
	return func(gs *gracefulServer) {
		if f != nil {
			gs.forceStop = f
		}
	}
}

// Start returns the implemented function to start a server
func (gs *gracefulServer) Start() error { return gs.start() }

// Stop returns the implemented function to stop a server
func (gs *gracefulServer) Stop(ctx context.Context) { gs.stop(ctx) }

// ForceStop returns the implemented function to forcibly stop a server
func (gs *gracefulServer) ForceStop() { gs.forceStop() }
