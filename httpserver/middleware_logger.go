package httpserver

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/telmoandrade/go-library/logger"
)

func realIPExtractHeader(r *http.Request) string {
	ip := r.Header.Get("True-Client-IP")
	if ip != "" {
		return ip
	}

	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		i := strings.IndexByte(xff, ',')
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	}

	return ip
}

func realIP(r *http.Request) string {
	ip := realIPExtractHeader(r)

	if ip != "" && net.ParseIP(ip) != nil {
		return ip
	}

	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}
	return r.RemoteAddr
}

func sinceRound(since time.Duration) time.Duration {
	if since > time.Second {
		since = since.Round(time.Second)
	} else if since > time.Millisecond {
		since = since.Round(time.Millisecond)
	} else if since > time.Microsecond {
		since = since.Round(time.Microsecond)
	}

	return since
}

// MiddlewareLogging is a middleware that logs each incoming request along with useful metadata regarding the request.
//
// Response Status Handling:
//   - Error: For response status < 100 and >= 500
//   - Warn: For response status < 200 and >= 400
//   - Info: Other response status
//
// Log Identifier Handling:
//   - If the X-Logger-ID header is present in the request, its value will be used as the log identifier.
//   - If the header is not present or if the value is invalid, a new log identifier will be generated using UUID v7.
//   - The log identifier is then added to the context [logger.ContextLogID].
//
// Log Level Handling:
//   - If the X-Logger-Level header is present in the request, its value will be used as the minimum log level.
//     Allowing lower priority logs at runtime.
//   - The minimum log level is then added to the context [logger.ContextMinLevel].
//
// Important Note:
//   - MiddlewareLogging should be positioned before any other middleware that may alter the response, such as [MiddlewareRecover].
//   - Must be used with [logger.NewHandler] to register the log handle and allow lower priority logging at runtime.
//
// Example:
//
//	mux := httpserver.NewServeMux()
//	mux.Use(httpserver.MiddlewareLogging)   // <--<< MiddlewareLogging must come before MiddlewareRecover
//	mux.Use(httpserver.MiddlewareRecover)
//	mux.Get("/", handler)
func MiddlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ctx := r.Context()

		ctx, u := logger.LogId(ctx, r.Header.Get("X-Logger-ID"))
		w.Header().Add("X-Logger-ID", u.String())

		loggerLevel := r.Header.Get("X-Logger-Level")
		ctx, err := logger.MinLevel(ctx, loggerLevel)
		if err == nil {
			w.Header().Add("X-Logger-Level", loggerLevel)
		}

		wrw := &wrapResponseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(wrw, r.WithContext(ctx))

		since := time.Since(start)
		slogAny := []any{
			slog.Group("log",
				slog.String("id", u.String()),
			),
			slog.Group("request",
				slog.String("method", r.Method),
				slog.String("pattern", r.Pattern),
				slog.String("path", r.URL.Path),
				slog.Int64("size", r.ContentLength),
			),
			slog.Group("user",
				slog.String("agent", r.Header.Get("User-Agent")),
				slog.String("protocol", r.Proto),
				slog.String("host", r.Host),
				slog.String("ip", realIP(r)),
			),
			slog.Group("response",
				slog.Int("status", wrw.code),
				slog.Int64("size", wrw.bytes),
				slog.Float64("time", since.Seconds()),
			),
		}

		msg := fmt.Sprintf("HTTP Response %03d %dB %v %s %s", wrw.code, wrw.bytes, sinceRound(since), r.Method, r.URL.Path)

		if wrw.code < http.StatusContinue || wrw.code >= http.StatusInternalServerError {
			slog.ErrorContext(ctx, msg, slogAny...)
		} else if wrw.code < http.StatusOK || wrw.code >= http.StatusBadRequest {
			slog.WarnContext(ctx, msg, slogAny...)
		} else {
			slog.InfoContext(ctx, msg, slogAny...)
		}
	})
}
