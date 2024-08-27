//go:generate mockgen -source=server.go -destination=mock_server_test.go -package server

package server

import (
	"context"
)

type (
	server struct {
		start     func() error
		stop      func(context.Context)
		forceStop func()
	}

	Server interface {
		Start() error
		Stop(context.Context)
		ForceStop()
	}

	OptionServer func(*server)
)

func NewServer(opts ...OptionServer) Server {
	s := &server{
		start:     func() error { return nil },
		stop:      func(context.Context) {},
		forceStop: func() {},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func WithStart(f func() error) OptionServer         { return func(s *server) { s.start = f } }
func WithStop(f func(context.Context)) OptionServer { return func(s *server) { s.stop = f } }
func WithForceStop(f func()) OptionServer           { return func(s *server) { s.forceStop = f } }

func (s *server) Start() error             { return s.start() }
func (s *server) Stop(ctx context.Context) { s.stop(ctx) }
func (s *server) ForceStop()               { s.forceStop() }
