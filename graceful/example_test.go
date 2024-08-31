package graceful_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/telmoandrade/go-library/graceful"
)

func ExampleNewGracefulServer() {
	gs := graceful.NewGracefulShutdown(
		graceful.WithServers(
			graceful.NewGracefulServer(
				graceful.WithStart(func() error {
					fmt.Println("Server start")
					return errors.New("error")
				}),
				graceful.WithStop(func(ctx context.Context) {
					fmt.Println("Server stop")
				}),
			),
		),
	)

	gs.Run(context.Background())
	// Output:
	// Server start
	// Server stop
}

func ExampleNewGracefulServerHttp() {
	ctx, stop := context.WithCancel(context.Background())
	stop()

	s := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: &http.ServeMux{},
	}

	gs := graceful.NewGracefulShutdown(
		graceful.WithServers(
			graceful.NewGracefulServerHttp(s),
		),
	)

	gs.Run(ctx)
	// Output:
}
