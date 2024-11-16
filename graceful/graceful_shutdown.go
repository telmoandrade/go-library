package graceful

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"slices"
	"sync"
	"time"
)

type (
	gracefulShutdown struct {
		ctx             context.Context
		cancelCtx       context.CancelFunc
		wg              sync.WaitGroup
		timeout         time.Duration
		gracefulServers []GracefulServer
		once            sync.Once
		notifyShutdown  func()
	}

	// GracefulShutdown is responsible for managing the lifecycle of the graceful shutdown handler, overseeing the startup, shutdown,
	// and completion of multiple servers in an organized and predictable manner.
	GracefulShutdown interface {
		// Runs executes the server startup and shutdown process, handling a life cycle of all servers.
		// If the context used to control the shutdown process signals a timeout or cancellation, GracefulShutdown will initiate a graceful shutdown.
		//
		// Starts all registered servers and waits for them to close gracefully.
		Run(ctx context.Context)
	}

	// OptionGracefulShutdown is used to apply configurations to a [GracefulShutdown] when creating it with [NewGracefulShutdown].
	OptionGracefulShutdown func(*gracefulShutdown)
)

// NewGracefulShutdown returns a new [GracefulShutdown] handler to ensure graceful shutdown of the application.
// A variadic set of [OptionGracefulShutdown] that can configure the behavior of the shutdown handler.
func NewGracefulShutdown(opts ...OptionGracefulShutdown) GracefulShutdown {
	ctx, cancelCtx := context.WithCancel(context.Background())

	gs := &gracefulShutdown{
		ctx:             ctx,
		cancelCtx:       cancelCtx,
		gracefulServers: []GracefulServer{},
		notifyShutdown:  func() {},
	}

	for _, opt := range opts {
		opt(gs)
	}

	return gs
}

// WithTimeout is an [OptionGracefulShutdown] that this option sets the timeout period the graceful
// shutdown handler will wait before forcibly terminating the servers.
//
// Default Behavior:
//   - If the timeout is set to 0, the handler will wait indefinitely for the [GracefulServer.Stop] method to complete.
func WithTimeout(t time.Duration) OptionGracefulShutdown {
	return func(gs *gracefulShutdown) { gs.timeout = t }
}

// WithServers is an [OptionGracefulShutdown] that adds a variable list of servers that the graceful shutdown handler will manage.
// Servers must implement the [GracefulServer] interface.
func WithServers(servers ...GracefulServer) OptionGracefulShutdown {
	return func(gs *gracefulShutdown) {
		gs.gracefulServers = slices.Clip(
			slices.DeleteFunc(append(gs.gracefulServers, servers...), func(s GracefulServer) bool {
				return s == nil
			}),
		)
	}
}

// WithNotifyShutdown is an [OptionGracefulShutdown] that defines the function to notify the shutdown process begins.
// The function that will be invoked to notify the shutdown process begins.
func WithNotifyShutdown(fn func()) OptionGracefulShutdown {
	return func(gs *gracefulShutdown) {
		if fn != nil {
			gs.notifyShutdown = fn
		}
	}
}

func (gs *gracefulShutdown) runServer(s GracefulServer) {
	gs.wg.Add(1)

	go func() {
		<-gs.ctx.Done()
		defer gs.wg.Done()

		showdownCtx, cancelShowdownCtx := context.WithCancel(context.Background())
		if gs.timeout > 0 {
			showdownCtx, cancelShowdownCtx = context.WithTimeout(showdownCtx, gs.timeout)
		}
		go func() {
			<-showdownCtx.Done()
			if errors.Is(showdownCtx.Err(), context.DeadlineExceeded) {
				s.ForceStop()
			}
		}()

		s.Stop(showdownCtx)
		cancelShowdownCtx()
	}()

	go func() {
		if err := s.Start(); err != nil {
			gs.cancelCtx()
		}
	}()
}

func (gs *gracefulShutdown) Run(ctx context.Context) {
	if len(gs.gracefulServers) == 0 {
		return
	}

	gs.once.Do(func() {
		go func() {
			<-gs.ctx.Done()
			gs.notifyShutdown()
		}()

		for _, s := range gs.gracefulServers {
			gs.runServer(s)
		}

		signalCtx, cancelCtx := signal.NotifyContext(ctx, os.Interrupt)

		go func() {
			<-gs.ctx.Done()
			cancelCtx()
		}()

		<-signalCtx.Done()

		gs.cancelCtx()
		gs.wg.Wait()
	})
}
