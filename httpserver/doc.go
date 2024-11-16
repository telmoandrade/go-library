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
// Default behavior when registering a handler:
//   - If no route was registered in '[HOST]/', register a handler for not found routes, inheriting the middleware stack (log HandlerNotFound).
//   - If no route was registered in '[HOST]/[PATH]', register a handler for not allowed methods routes, inheriting the middleware stack (log HandlerMethodNotAllowed).
//   - Then register the defined handler (log HandlerFn).
//
// A variadic set of [OptionServeMux] used to configure the behavior of the [ServeMux]:
//   - [WithHandlerOptionsMaxAge]: Defines in seconds the maximum age for the Cache-Control header response in the options method handlers.
//   - [WithCors]: Defines information that will be used in the handlers for CORS processing.
//   - [WithHandlerNotFound]: Defines the handler for handling not found routes.
//   - [WithHandlerMethodNotAllowed]: Defines the handler not allowed methods routes.
//
// Methods for adding middleware:
//   - [ServeMux.Use]: Appends one or more middlewares.
//   - [ServeMux.With]: Appends one or more middlewares and register the [Handle] inline.
//
// Methods for managing the routing path:
//   - [ServeMux.Group]: Inline router manager, inheriting the middleware stack.
//   - [ServeMux.Route]: Subrouter manager, inheriting the middleware stack.
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
//   - [MiddlewareTrace]: Adds attributes to spans and metrics for telemetry purposes.
//   - [MiddlewareRecover]: Recovers from panics, logs the panic, and responds with an HTTP status of 500 (Internal Server Error).
package httpserver
