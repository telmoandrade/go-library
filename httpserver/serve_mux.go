package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"
)

type (
	// Handle interface allows the registration of HTTP handlers under the current routing path plus the specified pattern.
	// It includes methods for each standard HTTP method.
	Handle interface {
		// Connect registers a handler for the HTTP CONNECT method, under the current routing path plus the specified pattern.
		Connect(pattern string, handlerFn http.HandlerFunc)
		// Delete registers a handler for the HTTP DELETE method, under the current routing path plus the specified pattern.
		Delete(pattern string, handlerFn http.HandlerFunc)
		// Get registers a handler for the HTTP GET method, under the current routing path plus the specified pattern.
		Get(pattern string, handlerFn http.HandlerFunc)
		// Head registers a handler for the HTTP HEAD method, under the current routing path plus the specified pattern.
		Head(pattern string, handlerFn http.HandlerFunc)
		// Options registers a handler for the HTTP OPTIONS method, under the current routing path plus the specified pattern.
		Options(pattern string, handlerFn http.HandlerFunc)
		// Patch registers a handler for the HTTP PATCH method, under the current routing path plus the specified pattern.
		Patch(pattern string, handlerFn http.HandlerFunc)
		// Post registers a handler for the HTTP POST method, under the current routing path plus the specified pattern.
		Post(pattern string, handlerFn http.HandlerFunc)
		// Put registers a handler for the HTTP PUT method, under the current routing path plus the specified pattern.
		Put(pattern string, handlerFn http.HandlerFunc)
		// Trace registers a handler for the HTTP TRACE method, under the current routing path plus the specified pattern.
		Trace(pattern string, handlerFn http.HandlerFunc)
		// Method registers a handler for the custom HTTP method, under the current routing path plus the specified pattern.
		Method(method, pattern string, handlerFn http.HandlerFunc)
	}

	// Router interface extends the [Handle] interface.
	// It is designed to manage routing paths, middleware registration, and handler registrations.
	Router interface {
		Handle

		// Use appends one or more middlewares to the middleware stack for the Router in the current routing path.
		Use(middlewares ...func(http.Handler) http.Handler)
		// With appends one or more middlewares to the middleware stack in the current routing path to register the inline Handle.
		With(middlewares ...func(http.Handler) http.Handler) Handle
		// Group return a new inline Router under the current routing path plus the specified pattern, inheriting the middleware stack.
		Group(pattern string) Router
		// Route allowing additional routes to be defined within the subrouter under the current routing path plus the specified pattern, inheriting the middleware stack.
		Route(pattern string, fn func(sub Router))
		// Mount attaches another http.Handler under the current routing path plus the specified pattern, useful for integrating external routers or handlers.
		Mount(pattern string, h http.Handler)
	}

	// ServeMux extends [http.Handler] designed to manage routing paths, middleware registration,
	// and handler registrations of the standard library [http.Server].
	// It serves as a versatile routing mechanism that can handle middleware and nested routers efficiently.
	//
	// Behavior:
	//   - The ServeMux stores the routing path to register middlewares and handlers.
	//   - If a pattern using the host conflicts with one that is already registered, it will cause a panic.
	ServeMux struct {
		serveMux    *http.ServeMux
		middlewares []func(http.Handler) http.Handler
		pattern     string
	}

	contextKey struct {
		name string
	}
)

var contextRoutePath = &contextKey{"routePath"}

// NewServeMux creates and returns a new instance of [ServeMux] with enhanced routing and middleware capabilities.
func NewServeMux() *ServeMux {
	return &ServeMux{
		serveMux: http.NewServeMux(),
	}
}

// Use appends one or more middlewares to the middleware stack for the [Router] in the current routing path.
func (mux *ServeMux) Use(middlewares ...func(http.Handler) http.Handler) {
	mux.middlewares = append(mux.middlewares, middlewares...)
}

// With appends one or more middlewares to the middleware stack in the current routing path to register the inline Handle.
func (mux *ServeMux) With(middlewares ...func(http.Handler) http.Handler) Handle {
	return &ServeMux{
		pattern:     mux.pattern,
		middlewares: append(mux.middlewares, middlewares...),
		serveMux:    mux.serveMux,
	}
}

// Group return a new inline [Router] under the current routing path plus the specified pattern, inheriting the middleware stack.
func (mux *ServeMux) Group(pattern string) Router {
	return &ServeMux{
		pattern:     joinPattern(mux.pattern, pattern),
		middlewares: mux.middlewares,
		serveMux:    mux.serveMux,
	}
}

// Route allowing additional routes to be defined within the sub-[Router] under the current routing path plus the
// specified pattern, inheriting the middleware stack.
func (mux *ServeMux) Route(pattern string, fn func(sub Router)) {
	subRouter := &ServeMux{
		pattern:     joinPattern(mux.pattern, pattern),
		middlewares: mux.middlewares,
		serveMux:    mux.serveMux,
	}
	fn(subRouter)

}

// Mount attaches another [http.Handler] under the current routing path plus the specified pattern, useful for
// integrating external routers or handlers.
//
// Important Note:
//   - Avoid using this method to attach another multiplexer as it does not inherit the middleware stack.
//   - It is slower to resolve the multiplexer for HTTP requests compared to using the built-in routing methods.
func (mux *ServeMux) Mount(pattern string, handler http.Handler) {
	targetPattern := strings.TrimSpace(joinPattern(mux.pattern, pattern))
	middlewareMount := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), contextRoutePath, targetPattern)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	stripPattern := ""
	stripPatternSplit := strings.SplitN(targetPattern, "/", 2)
	if len(stripPatternSplit) == 2 {
		stripPattern = fmt.Sprintf("/%s", stripPatternSplit[1])
	}

	mux.serveMux.Handle(targetPattern+"/", http.StripPrefix(stripPattern, middlewareMount(handler)))
}

func (mux *ServeMux) addRoute(method string, pattern string, handlerFn http.Handler) {
	targetPattern := joinPattern(mux.pattern, pattern)
	targetHandler := handlerFn

	middlewareAddRoute := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if targetPattern, ok := r.Context().Value(contextRoutePath).(string); ok {
				patternSplited := strings.SplitN(r.Pattern, " ", 2)
				if len(patternSplited) == 2 {
					r.Pattern = fmt.Sprintf("%s %s", patternSplited[0], joinPattern(targetPattern, patternSplited[1]))
				} else {
					r.Pattern = joinPattern(targetPattern, patternSplited[0])
				}
			}
			next.ServeHTTP(w, r)
		})
	}

	for i := len(mux.middlewares) - 1; i >= 0; i-- {
		targetHandler = mux.middlewares[i](targetHandler)
	}
	targetHandler = middlewareAddRoute(targetHandler)

	mux.serveMux.Handle(strings.TrimSpace(fmt.Sprintf("%s %s", method, targetPattern)), targetHandler)
}

// Connect registers a handler for the HTTP CONNECT method, under the current routing path plus the specified pattern.
func (mux *ServeMux) Connect(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("CONNECT", pattern, handlerFn)
}

// Delete registers a handler for the HTTP DELETE method, under the current routing path plus the specified pattern.
func (mux *ServeMux) Delete(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("DELETE", pattern, handlerFn)
}

// Get registers a handler for the HTTP GET method, under the current routing path plus the specified pattern.
func (mux *ServeMux) Get(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("GET", pattern, handlerFn)
}

// Head registers a handler for the HTTP HEAD method, under the current routing path plus the specified pattern.
func (mux *ServeMux) Head(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("HEAD", pattern, handlerFn)
}

// Options registers a handler for the HTTP OPTIONS method, under the current routing path plus the specified pattern.
func (mux *ServeMux) Options(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("OPTIONS", pattern, handlerFn)
}

// Patch registers a handler for the HTTP PATCH method, under the current routing path plus the specified pattern.
func (mux *ServeMux) Patch(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("PATCH", pattern, handlerFn)
}

// Post registers a handler for the HTTP POST method, under the current routing path plus the specified pattern.
func (mux *ServeMux) Post(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("POST", pattern, handlerFn)
}

// Put registers a handler for the HTTP PUT method, under the current routing path plus the specified pattern.
func (mux *ServeMux) Put(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("PUT", pattern, handlerFn)
}

// Trace registers a handler for the HTTP TRACE method, under the current routing path plus the specified pattern.
func (mux *ServeMux) Trace(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("TRACE", pattern, handlerFn)
}

// Method registers a handler for the custom HTTP method, under the current routing path plus the specified pattern.
func (mux *ServeMux) Method(method, pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute(strings.ToUpper(method), pattern, handlerFn)
}

// ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL.
func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.serveMux.ServeHTTP(w, r)
}

func joinPattern(prefix, suffix string) string {
	prefixSplited := strings.SplitN(prefix, "/", 2)
	suffixSplited := strings.SplitN(suffix, "/", 2)

	prefixHost := prefixSplited[0]
	suffixHost := suffixSplited[0]
	if prefixHost != "" && suffixHost != "" && prefixHost != suffixHost {
		panic(fmt.Sprintf("Hostname conflict %s != %s", prefixHost, suffixHost))
	}

	hostPattern := prefixHost
	if suffixHost != "" {
		hostPattern = suffixHost
	}
	prefixPattern := ""
	if len(prefixSplited) == 2 {
		prefixPattern = fmt.Sprintf("/%s", prefixSplited[1])
	}
	suffixPattern := ""
	if len(suffixSplited) == 2 {
		suffixPattern = fmt.Sprintf("/%s", suffixSplited[1])
	}

	pattern := path.Join(hostPattern, path.Join(prefixPattern, suffixPattern))

	return pattern
}
