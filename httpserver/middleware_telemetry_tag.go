package httpserver

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/telmoandrade/go-library/logger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// MiddlewareTelemetryTag is a middleware that adds attributes to spans and metrics for telemetry purposes.
//
// Adds telemetry attributes for monitoring:
//   - http.route: Indicates the route of the HTTP request used in spans and metrics.
//   - log.id: Log identifier associated with the request used in spans.
func MiddlewareTelemetryTag(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routePath := r.Pattern[strings.Index(r.Pattern, "/"):]

		attr := semconv.HTTPRoute(routePath)

		span := trace.SpanFromContext(r.Context())
		span.SetAttributes(attr)

		u := logger.LogIDFromContext(r.Context())
		if u != uuid.Nil {
			span.SetAttributes(
				attribute.Key("log.id").String(u.String()),
			)
		}

		labeler, _ := otelhttp.LabelerFromContext(r.Context())
		labeler.Add(attr)

		next.ServeHTTP(w, r)
	})
}
