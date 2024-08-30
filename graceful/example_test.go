package graceful_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/telmoandrade/go-library/graceful"
)

func ExempleNewGracefulServer() {
	ctx, stop := context.WithCancel(context.Background())

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

	go func() {
		<-time.After(100 * time.Millisecond)
		stop()
	}()

	gs.Run(ctx)

	// Output:
	// Server start
	// Server stop
}

func ExempleNewGracefulServerHttp() {
	ctx, stop := context.WithCancel(context.Background())

	s := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: &http.ServeMux{},
	}

	gs := graceful.NewGracefulShutdown(
		graceful.WithServers(
			graceful.NewGracefulServerHttp(s),
		),
	)

	go func() {
		<-time.After(100 * time.Millisecond)
		stop()
	}()

	gs.Run(ctx)

	// Output:
	// 2006/01/02 15:04:05 INFO [HTTP SERVER] Starting
	// 2006/01/02 15:04:05 INFO [HTTP SERVER] Closing
	// 2006/01/02 15:04:05 INFO [HTTP SERVER] Closed
}
