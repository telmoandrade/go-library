package server

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
	control struct {
		ctx       context.Context
		cancelCtx context.CancelFunc
		wg        sync.WaitGroup
		timeout   time.Duration
		servers   []Server
		once      sync.Once
	}

	OptionControl func(*control)
)

func WithTimeout(t time.Duration) OptionControl { return func(c *control) { c.timeout = t } }
func WithServers(s ...Server) OptionControl {
	return func(c *control) {
		c.servers = slices.Clip(
			slices.DeleteFunc(s, func(t Server) bool {
				return t == nil
			}),
		)
	}
}

func NewControl(opts ...OptionControl) *control {
	ctx, cancelCtx := context.WithCancel(context.Background())

	c := &control{
		ctx:       ctx,
		cancelCtx: cancelCtx,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *control) runServer(s Server) {
	c.wg.Add(1)

	go func() {
		<-c.ctx.Done()
		defer c.wg.Done()

		showdownCtx, cancelShowdownCtx := context.WithCancel(context.Background())
		if c.timeout > 0 {
			showdownCtx, cancelShowdownCtx = context.WithTimeout(showdownCtx, c.timeout)
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
			c.cancelCtx()
		}
	}()
}

func (c *control) Run(ctx context.Context) {
	c.once.Do(func() {
		if len(c.servers) > 0 {
			for _, s := range c.servers {
				c.runServer(s)
			}

			signalCtx, cancelCtx := signal.NotifyContext(ctx, os.Interrupt)
			defer cancelCtx()
			<-signalCtx.Done()
		}

		c.cancelCtx()
		c.wg.Wait()
	})
}
