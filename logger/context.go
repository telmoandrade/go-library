package logger

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

// ErrInvalidMinLevel is returned by [MinLevel] when the given level
// is not present in the allowed list (debug, info, warn, error).
var ErrInvalidMinLevel = errors.New("invalid min level")

// LogId returns a new context and log ID after embedding the log identifier in the provided context, using the [ContextLogID] key.
//
// Behavior:
//   - Creates a new log id if it does not already exist in the context.
//   - Try to use propagation to create the new log id.
func LogId(ctx context.Context, propagation string) (context.Context, uuid.UUID) {
	u, _ := ctx.Value(ContextLogID).(uuid.UUID)
	if u != uuid.Nil {
		return ctx, u
	}

	u, err := uuid.Parse(propagation)
	if err != nil {
		u, _ = uuid.NewV7()
	}

	return context.WithValue(ctx, ContextLogID, u), u
}

// MinLevel returns a new context after embedding the minimum log level in the provided context, using the [ContextMinLevel] key.
// Allowing lower priority logs at runtime.
//
// The minimum log level can be one of the following options:
//   - debug: Logs at the debug level, used for detailed information useful for debugging.
//   - info: Logs at the informational level, used for general messages about application progress.
//   - warn: Logs at the warning level, indicating potential issues that do not cause immediate errors.
//   - error: Logs at the error level, used for serious issues that need attention.
func MinLevel(ctx context.Context, level string) (context.Context, error) {
	switch strings.ToLower(level) {
	case "debug":
		ctx = context.WithValue(ctx, ContextMinLevel, slog.LevelDebug)
	case "info":
		ctx = context.WithValue(ctx, ContextMinLevel, slog.LevelInfo)
	case "warn":
		ctx = context.WithValue(ctx, ContextMinLevel, slog.LevelWarn)
	case "error":
		ctx = context.WithValue(ctx, ContextMinLevel, slog.LevelError)
	default:
		return ctx, ErrInvalidMinLevel
	}
	return ctx, nil
}
