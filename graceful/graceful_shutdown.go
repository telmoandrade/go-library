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
	}

	// OptionGracefulShutdown are parameters used in [NewGracefulShutdown]
	OptionGracefulShutdown func(*gracefulShutdown)
)

// NewGracefulShutdown returns a handler to ensure graceful shutdown of the application
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

// WithTimeout is an [OptionGracefulShutdown] used to set a timeout that the graceful shutdown handler will wait before calling GracefulServer.ForceStop
//
// If not set, the graceful shutdown handler will wait for GracefulServer.Stop to complete, regardless of how long it takes.
//
// If the graceful shutdown handler calls GracefulServer.ForceStop, it will NOT wait for its call to complete
func WithTimeout(t time.Duration) OptionGracefulShutdown {
	return func(gs *gracefulShutdown) { gs.timeout = t }
}

// WithServers is an [OptionGracefulShutdown] used to register the various [GracefulServer]
func WithServers(servers ...GracefulServer) OptionGracefulShutdown {
	return func(gs *gracefulShutdown) {
		gs.gracefulServers = slices.Clip(
			slices.DeleteFunc(servers, func(s GracefulServer) bool {
				return s == nil
			}),
		)
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
			s.Stop(showdownCtx)
			cancelShowdownCtx()
		}()

		<-showdownCtx.Done()
		if errors.Is(showdownCtx.Err(), context.DeadlineExceeded) {
			go s.ForceStop()
		}
	}()

	go func() {
		if err := s.Start(); err != nil {
			gs.cancelCtx()
		}
	}()
}

// Run start all servers and wait for them to close
func (gs *gracefulShutdown) Run(ctx context.Context) {
	if len(gs.gracefulServers) == 0 {
		return
	}

	gs.once.Do(func() {
		for _, s := range gs.gracefulServers {
			gs.runServer(s)
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
