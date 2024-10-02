package logger

import (
	"context"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

// LogIDFromContext retrieves the log identifier stored in a given context under the [ContextLogID] key.
// If no log identifier is found in the context, it may return a zero-valued UUID.
func LogIDFromContext(ctx context.Context) uuid.UUID {
	u, _ := ctx.Value(ContextLogID).(uuid.UUID)
	return u
}

// WithContextLogID returns a new context after embedding the log identifier in the provided context, using the [ContextLogID] key.
func WithContextLogID(ctx context.Context, u uuid.UUID) context.Context {
	return context.WithValue(ctx, ContextLogID, u)
}

// WithContextMinLevel returns a new context after embedding the minimum log level in the provided context, using the [ContextMinLevel] key.
// Allowing lower priority logs at runtime.
//
// The minimum log level can be one of the following options:
//   - debug: Logs at the debug level, used for detailed information useful for debugging.
//   - info: Logs at the informational level, used for general messages about application progress.
//   - warn: Logs at the warning level, indicating potential issues that do not cause immediate errors.
//   - error: Logs at the error level, used for serious issues that need attention.
func WithContextMinLevel(ctx context.Context, level string) context.Context {
	switch strings.ToLower(level) {
	case "debug":
		ctx = context.WithValue(ctx, ContextMinLevel, slog.LevelDebug)
	case "info":
		ctx = context.WithValue(ctx, ContextMinLevel, slog.LevelInfo)
	case "warn":
		ctx = context.WithValue(ctx, ContextMinLevel, slog.LevelWarn)
	case "error":
		ctx = context.WithValue(ctx, ContextMinLevel, slog.LevelError)
	}
	return ctx
}
