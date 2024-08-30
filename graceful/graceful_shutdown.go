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
		ctx       context.Context
		cancelCtx context.CancelFunc
		wg        sync.WaitGroup
		timeout   time.Duration
		servers   []GracefulServer
		once      sync.Once
	}

	// OptionGracefulShutdown used in [NewGracefulShutdown]
	OptionGracefulShutdown func(*gracefulShutdown)
)

// NewGracefulShutdown return a new graceful shutdown manager
func NewGracefulShutdown(opts ...OptionGracefulShutdown) *gracefulShutdown {
	ctx, cancelCtx := context.WithCancel(context.Background())

	gs := &gracefulShutdown{
		ctx:       ctx,
		cancelCtx: cancelCtx,
	}

	for _, opt := range opts {
		opt(gs)
	}

	return gs
}

// WithTimeout sets a timeout to wait for each server to shut down before calling force shutdown
//
// If not configured, the time required for the server to stop will be waited.
func WithTimeout(t time.Duration) OptionGracefulShutdown {
	return func(gs *gracefulShutdown) { gs.timeout = t }
}

// WithServers register a list of servers that the manager will handle
func WithServers(servers ...GracefulServer) OptionGracefulShutdown {
	return func(gs *gracefulShutdown) {
		gs.servers = slices.Clip(
			slices.DeleteFunc(servers, func(server GracefulServer) bool {
				return server == nil
			}),
		)
	}
}

func (gs *gracefulShutdown) runServer(server GracefulServer) {
	gs.wg.Add(1)

	go func() {
		<-gs.ctx.Done()
		defer gs.wg.Done()

		showdownCtx, cancelShowdownCtx := context.WithCancel(context.Background())
		if gs.timeout > 0 {
			showdownCtx, cancelShowdownCtx = context.WithTimeout(showdownCtx, gs.timeout)
		}
		go func() {
			server.Stop(showdownCtx)
			cancelShowdownCtx()
		}()

		<-showdownCtx.Done()
		if errors.Is(showdownCtx.Err(), context.DeadlineExceeded) {
			go server.ForceStop()
		}
	}()

	go func() {
		if err := server.Start(); err != nil {
			gs.cancelCtx()
		}
	}()
}

// Run start all servers and wait for them to close
func (gs *gracefulShutdown) Run(ctx context.Context) {
	if len(gs.servers) == 0 {
		return
	}

	gs.once.Do(func() {
		for _, server := range gs.servers {
			gs.runServer(server)
		}

		signalCtx, cancelCtx := signal.NotifyContext(ctx, os.Interrupt)
		defer cancelCtx()

		go func() {
			<-gs.ctx.Done()
			cancelCtx()
		}()

		<-signalCtx.Done()

		gs.cancelCtx()
		gs.wg.Wait()
	})
}
