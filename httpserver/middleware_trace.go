package httpserver

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/telmoandrade/go-library/logger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// MiddlewareTrace is a middleware that adds attributes to spans and metrics for telemetry purposes.
//
// Adds telemetry attributes for monitoring:
//   - attribute http.route: Indicates the pattern of the HTTP request used in spans and metrics.
//   - attribute log.id: Log identifier associated with the request used in spans.
//   - header Traceparent: Trace span associated with the request used in the response.
//
// Important Note:
//   - MiddlewareTrace should be placed after [MiddlewareLogging] middleware.
//   - MiddlewareTrace should be positioned before any other middleware that may alter the response, such as [MiddlewareRecover].
//
// Example:
//
//	mux := httpserver.NewServeMux()
//	mux.Use(httpserver.MiddlewareLogging)
//	mux.Use(httpserver.MiddlewareTrace)   // <--<< MiddlewareTrace must come before MiddlewareRecover and after MiddlewareLogging
//	mux.Use(httpserver.MiddlewareRecover)
//	mux.Get("/", handler)
func MiddlewareTrace(next http.Handler) http.Handler {
	handler := otelhttp.NewMiddleware("",
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return r.Pattern
		}),
	)

	return handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attr := semconv.HTTPRoute(r.Pattern)

		span := trace.SpanFromContext(r.Context())
		span.SetAttributes(attr)

		u, _ := r.Context().Value(logger.ContextLogID).(uuid.UUID)
		if u != uuid.Nil {
			span.SetAttributes(
				attribute.Key("log.id").String(u.String()),
			)
		}

		labeler, _ := otelhttp.LabelerFromContext(r.Context())
		labeler.Add(attr)

		otel.GetTextMapPropagator().Inject(r.Context(), propagation.HeaderCarrier(w.Header()))

		next.ServeHTTP(w, r)
	}))
}
