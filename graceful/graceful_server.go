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

	// GracefulServer defines the required methods that any server must implement to participate in the graceful shutdown handler.
	// It is used as a parameter for the [WithServers].
	GracefulServer interface {
		// Start method responsible for starting the server.
		// It should return an error if the startup fails.
		Start() error
		// Stop method responsible for stopping the server gracefully.
		// It has a context parameter to manage timeout signals.
		Stop(context.Context)
		// ForceStop method responsible for forcibly stopping the server if the graceful stop does not complete within the allotted time.
		ForceStop()
	}

	// OptionGracefulServer is used to apply configurations to a [GracefulServer] when creating it with [NewGracefulServer].
	OptionGracefulServer func(*gracefulServer)
)

// NewGracefulServer returns a new [GracefulServer].
// A variadic set of options for configuring the behavior of the server.
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

// WithStart is an [OptionGracefulServer] that defines the function to start [GracefulServer.Start].
// The function that will be invoked to start the server, it should return an error if the startup fails.
func WithStart(fn func() error) OptionGracefulServer {
	return func(gs *gracefulServer) {
		if fn != nil {
			gs.start = fn
		}
	}
}

// WithStop is an [OptionGracefulServer] that defines the function to gracefully stop [GracefulServer.Stop].
// The function responsible for stopping the server, it has a [context.Context] parameter to manage timeout signals.
func WithStop(fn func(context.Context)) OptionGracefulServer {
	return func(gs *gracefulServer) {
		if fn != nil {
			gs.stop = fn
		}
	}
}

// WithForceStop is an [OptionGracefulServer] that defines the function to forcibly stop [GracefulServer.ForceStop].
// The function responsible for forcibly stopping the server if the graceful stop does not complete within the allotted time.
//
// Important Note:
//   - When the timeout is reached, the graceful shutdown handler will call [GracefulServer.ForceStop].
func WithForceStop(fn func()) OptionGracefulServer {
	return func(gs *gracefulServer) {
		if fn != nil {
			gs.forceStop = fn
		}
	}
}

func (gs *gracefulServer) Start() error { return gs.start() }

func (gs *gracefulServer) Stop(ctx context.Context) { gs.stop(ctx) }

func (gs *gracefulServer) ForceStop() { gs.forceStop() }
