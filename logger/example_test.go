package logger_test

import (
	"context"
	"log/slog"

	"github.com/telmoandrade/go-library/logger"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

func ExampleNewLogger() {
	l := logger.NewLogger(
		logger.WithMinLevel(slog.LevelInfo),
		logger.WithHandler(otelslog.NewHandler("")),
		logger.WithMaxLevelAddSource(slog.LevelInfo),
	)

	slog.SetDefault(l)
	// Output:
}

func ExampleNewHandler() {
	h := logger.NewHandler(
		logger.WithMinLevel(slog.LevelInfo),
		logger.WithHandler(otelslog.NewHandler("")),
		logger.WithMaxLevelAddSource(slog.LevelInfo),
	)

	l := slog.New(h)

	slog.SetDefault(l)
	// Output:
}

func ExampleLogId() {
	slog.SetDefault(logger.NewLogger())

	ctx, u := logger.LogId(context.Background(), "")
	slog.DebugContext(ctx, "message debug", slog.String("id", u.String()))

	// Output:
}

func ExampleMinLevel() {
	slog.SetDefault(logger.NewLogger())

	ctx, _ := logger.MinLevel(context.Background(), "info")
	slog.DebugContext(ctx, "message debug")

	// Output:
}
