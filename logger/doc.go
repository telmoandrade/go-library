// Package logger provides a structured logging handler built on top of [slog.Logger].
//
// # Key Features:
//   - Minimum level control for processing a log, This can be dynamically overridden via context with [ContextMinLevel].
//   - Log identifier is dynamically added to log entry if context [ContextLogID] is set.
//   - Source code tracing if the log level is less than or equal to the configured maximum level [WithMaxLevelAddSource].
//   - Delegating the external handler to forward log entries to be processed [WithHandler].
package logger
