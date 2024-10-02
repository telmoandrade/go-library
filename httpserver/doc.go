// Package httpserver extends the functionality of the standard library [http.Server].
//
// # Key Features:
//   - Manage routing paths, middleware registration, and handler registrations of the standard library [http.Server].
//   - Some built-in middleware.
//
// # ServeMux
//
// ServeMux extends [http.Handler] designed to manage routing paths, middleware registration,
// and handler registrations of the standard library [http.Server].
// It serves as a versatile routing mechanism that can handle middleware and nested routers efficiently.
//
// Methods for adding middleware:
//   - [ServeMux.Use]: Appends one or more middlewares.
//   - [ServeMux.With]: Appends one or more middlewares and register the [Handle] inline.
//
// Methods for managing the routing path:
//   - [ServeMux.Group]: Inline router manager, inheriting the middleware stack.
//   - [ServeMux.Route]: Subrouter manager, inheriting the middleware stack.
//   - [ServeMux.Mount]: Attaches external routers or handlers.
//
// Methods for registering an http handler:
//   - [ServeMux.Connect]: Registers a handler for the HTTP CONNECT method.
//   - [ServeMux.Delete]: Registers a handler for the HTTP DELETE method.
//   - [ServeMux.Get]: Registers a handler for the HTTP GET method.
//   - [ServeMux.Head]: Registers a handler for the HTTP HEAD method.
//   - [ServeMux.Options]: Registers a handler for the HTTP OPTIONS method.
//   - [ServeMux.Patch]: Registers a handler for the HTTP PATCH method.
//   - [ServeMux.Post]: Registers a handler for the HTTP POST method.
//   - [ServeMux.Put]: Registers a handler for the HTTP PUT method.
//   - [ServeMux.Trace]: Registers a handler for the HTTP TRACE method.
//   - [ServeMux.Method]: Registers a handler for the custom HTTP method.
//
// # Middlewares
//   - [MiddlewareLogging]: Logs each incoming request along with useful metadata regarding the request.
//   - [MiddlewareRecover]: Recovers from panics, logs the panic, and responds with an HTTP status of 500 (Internal Server Error).
//   - [MiddlewareTelemetryTag]: Adds attributes to spans and metrics for telemetry purposes.
package httpserver
