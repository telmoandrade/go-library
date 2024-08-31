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

	// GracefulServer is a required interface to be used as a parameter in [WithServers]
	GracefulServer interface {
		Start() error
		Stop(context.Context)
		ForceStop()
	}

	// OptionGracefulServer are parameters used in [NewGracefulServer]
	OptionGracefulServer func(*gracefulServer)
)

// NewGracefulServer allows you to create a [GracefulServer] and can use a combination of [WithStart], [WithStop] and [WithForceStop]
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

// WithStart is an [OptionGracefulServer] used to define the function to start the [GracefulServer]
func WithStart(f func() error) OptionGracefulServer {
	return func(gs *gracefulServer) {
		if f != nil {
			gs.start = f
		}
	}
}

// WithStop is an [OptionGracefulServer] used to define the function to stop the [GracefulServer]
func WithStop(f func(context.Context)) OptionGracefulServer {
	return func(gs *gracefulServer) {
		if f != nil {
			gs.stop = f
		}
	}
}

// WithForceStop is an [OptionGracefulServer] used to define the function to forcefully stop the [GracefulServer]
func WithForceStop(f func()) OptionGracefulServer {
	return func(gs *gracefulServer) {
		if f != nil {
			gs.forceStop = f
		}
	}
}

// Start returns the implemented function to start a [GracefulServer]
func (gs *gracefulServer) Start() error { return gs.start() }

// Stop returns the implemented function to stop a [GracefulServer]
func (gs *gracefulServer) Stop(ctx context.Context) { gs.stop(ctx) }

// ForceStop returns the implemented function to forcibly stop a [GracefulServer]
func (gs *gracefulServer) ForceStop() { gs.forceStop() }
