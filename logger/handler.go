package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"

	"github.com/google/uuid"
)

type (
	loggerHandler struct {
		handler           slog.Handler
		minLevel          slog.Level
		maxLevelAddSource slog.Level
	}

	contextKey struct {
		name string
	}

	// Option is used to apply configurations to a loggerHandler when creating it with [NewHandler].
	Option func(*loggerHandler)
)

var (
	// ContextLogID is used to record the log identifier in the context.
	// The identifier helps to track and distinguish different logs throughout the system.
	ContextLogID = &contextKey{"logID"}
	// ContextMinLevel is used to set the minimum log level in the context.
	// The minimum log level defines the threshold for which logs are processed, allowing lower priority logs at runtime.
	ContextMinLevel = &contextKey{"minLevel"}

	_ slog.Handler = &loggerHandler{}
)

// NewLogger returns a new [slog.Logger].
// It uses the handler [NewHandler].
// A variadic set of [Option] used to configure the behavior of the handler.
func NewLogger(opts ...Option) *slog.Logger {
	return slog.New(NewHandler(opts...))
}

// NewHandler returns a new handler from [slog.Handler].
// A variadic set of [Option] used to configure the behavior of the handler.
//
// Behavior:
//   - Minimum level control for processing a log, This can be dynamically overridden via context with [ContextMinLevel].
//   - Log identifier is dynamically added to log entry if context [ContextLogID] is set.
//   - Source code tracing if the log level is less than or equal to the configured maximum level [WithMaxLevelAddSource].
//   - Delegating the external handler to forward log entries to be processed [WithHandler].
func NewHandler(opts ...Option) slog.Handler {
	lh := &loggerHandler{
		handler:           slog.NewTextHandler(os.Stdout, nil),
		minLevel:          slog.LevelInfo,
		maxLevelAddSource: slog.LevelDebug,
	}

	slog.Default()

	for _, opt := range opts {
		opt(lh)
	}

	return lh
}

// WithMinLevel is an [Option] that defines the minimum log level required for log entries to be processed.
// Any logs below this level will be ignored.
//
// Default:
//   - The default minimum log level is [slog.LevelInfo], meaning only informational messages and above (warnings, errors, etc.) will be logged.
func WithMinLevel(level slog.Level) Option {
	return func(lh *loggerHandler) {
		lh.minLevel = level
	}
}

// WithHandler is an [Option] that defines an external [slog.Handler] to which the processed log entries will be forwarded.
// This allows for chaining log handlers or customizing how logs are written or formatted.
//
// Default:
//   - If no custom handler is provided, the default handler is [slog.NewTextHandler], which outputs plain-text logs to the console [os.Stdout].
func WithHandler(handler slog.Handler) Option {
	return func(lh *loggerHandler) {
		if handler != nil {
			lh.handler = handler
		}
	}
}

// WithMaxLevelAddSource is an [Option] that defines the maximum log level at which source code information is included in the log entry.
//
// Default Behavior:
//   - If not set, the default maximum level is [slog.LevelDebug].
//   - If the log level is less than or equal to the maximum level, source code information will be added to the log entry.
//
// Source code information:
//   - source.function: The function where the log was generated.
//   - source.file: The file where the log was generated.
//   - source.line: The line number of the log statement.
func WithMaxLevelAddSource(level slog.Level) Option {
	return func(lh *loggerHandler) {
		lh.maxLevelAddSource = level
	}
}

func (lh *loggerHandler) Enabled(ctx context.Context, l slog.Level) bool {
	minLevel := lh.minLevel

	if l < minLevel {
		if l2, ok := ctx.Value(ContextMinLevel).(slog.Level); ok {
			minLevel = l2
		}
	}

	return l >= minLevel
}

func (lh *loggerHandler) Handle(ctx context.Context, r slog.Record) error {
	if u, ok := ctx.Value(ContextLogID).(uuid.UUID); ok {
		if u != uuid.Nil {
			r.AddAttrs(slog.Group("log",
				slog.String("id", u.String()),
			))
		}
	}

	if r.Level <= lh.maxLevelAddSource {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		r.AddAttrs(slog.Group("source",
			slog.String("function", f.Function),
			slog.String("file", f.File),
			slog.Int("line", f.Line),
		))
	}

	return lh.handler.Handle(ctx, r)
}

func (lh *loggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return lh.handler.WithAttrs(attrs)
}

func (lh *loggerHandler) WithGroup(name string) slog.Handler {
	return lh.handler.WithGroup(name)
}
