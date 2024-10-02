package logger_test

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
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

func ExampleLogIDFromContext() {
	u := logger.LogIDFromContext(context.Background())
	fmt.Println(u.String())

	// Output: 00000000-0000-0000-0000-000000000000
}

func ExampleWithContextLogID() {
	slog.SetDefault(logger.NewLogger())

	ctx := logger.WithContextLogID(context.Background(), uuid.Must(uuid.NewV7()))
	slog.DebugContext(ctx, "message debug")

	// Output:
}

func ExampleWithContextMinLevel() {
	slog.SetDefault(logger.NewLogger())

	ctx := logger.WithContextMinLevel(context.Background(), "info")
	slog.DebugContext(ctx, "message debug")

	// Output:
}
