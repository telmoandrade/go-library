package httpserver

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
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
		Route(pattern string, fn func(subMux Router))
	}

	// ServeMux extends [http.Handler] designed to manage routing paths, middleware registration,
	// and handler registrations of the standard library [http.Server].
	// It serves as a versatile routing mechanism that can handle middleware and nested routers efficiently.
	//
	// Behavior:
	//   - The ServeMux stores the routing path to register middlewares and handlers.
	//   - If a pattern using the host conflicts with one that is already registered, it will cause a panic.
	ServeMux interface {
		Router
		http.Handler
	}

	serveMuxConfig struct {
		handlerNotFound         http.HandlerFunc
		handlerMethodNotAllowed http.HandlerFunc
		cors                    *cors
		handlerOptionsMaxAge    int
	}

	serveMux struct {
		serveMux     *http.ServeMux
		middlewares  []func(http.Handler) http.Handler
		patternRoute *patternRoute
		routes       map[string]*serveMuxRoute
		config       *serveMuxConfig
	}

	// OptionServeMux is used to apply configurations to a [ServeMux] when creating it with [NewServeMux].
	OptionServeMux func(*serveMux)
)

// NewServeMux creates and returns a new instance of [ServeMux] with enhanced routing and middleware capabilities.
// A variadic set of [OptionServeMux] used to configure the behavior of the [ServeMux].
func NewServeMux(opts ...OptionServeMux) ServeMux {
	mux := &serveMux{
		serveMux:     http.NewServeMux(),
		middlewares:  []func(http.Handler) http.Handler{},
		patternRoute: newPatternRoute(""),
		routes:       map[string]*serveMuxRoute{},
		config: &serveMuxConfig{
			handlerNotFound:         defaultHandlerNotFound,
			handlerMethodNotAllowed: defaultHandlerMethodNotAllowed,
			handlerOptionsMaxAge:    86400,
		},
	}

	for _, opt := range opts {
		opt(mux)
	}

	return mux
}

// Use appends one or more middlewares to the middleware stack for the [Router] in the current routing path.
func (mux *serveMux) Use(middlewares ...func(http.Handler) http.Handler) {
	mux.middlewares = append(mux.middlewares, middlewares...)
}

// With appends one or more middlewares to the middleware stack in the current routing path to register the inline Handle.
func (mux *serveMux) With(middlewares ...func(http.Handler) http.Handler) Handle {
	return &serveMux{
		serveMux:     mux.serveMux,
		middlewares:  append(mux.middlewares, middlewares...),
		patternRoute: mux.patternRoute,
		routes:       mux.routes,
		config:       mux.config,
	}
}

// Group return a new inline [Router] under the current routing path plus the specified pattern, inheriting the middleware stack.
func (mux *serveMux) Group(pattern string) Router {
	return &serveMux{
		serveMux:     mux.serveMux,
		middlewares:  mux.middlewares,
		patternRoute: mux.patternRoute.join(pattern),
		routes:       mux.routes,
		config:       mux.config,
	}
}

// Route allowing additional routes to be defined within the sub-[Router] under the current routing path plus the
// specified pattern, inheriting the middleware stack.
func (mux *serveMux) Route(pattern string, fn func(sub Router)) {
	subRouter := &serveMux{
		serveMux:     mux.serveMux,
		middlewares:  mux.middlewares,
		patternRoute: mux.patternRoute.join(pattern),
		routes:       mux.routes,
		config:       mux.config,
	}
	fn(subRouter)
}

func validateHandler(handler http.Handler) {
	if handler == nil {
		panic(errors.New("httpserver: nil handler"))
	}
	if f, ok := handler.(http.HandlerFunc); ok && f == nil {
		panic(errors.New("httpserver: nil handler"))
	}
}

func (mux *serveMux) mountMiddlewares(smr *serveMuxRoute, handler http.Handler) http.Handler {
	if smr.cors != nil {
		handler = smr.middlewareCors(handler)
	}
	for i := len(mux.middlewares) - 1; i >= 0; i-- {
		handler = mux.middlewares[i](handler)
	}
	return handler
}

func (mux *serveMux) registerHandle(pattern, handlerKind string, handlerFn http.Handler) {
	slog.Info(fmt.Sprintf("[Register HTTP %s] %s", handlerKind, pattern))
	mux.serveMux.Handle(pattern, handlerFn)
}

func (mux *serveMux) registerServeMuxRoute(pattern string, createFn func(smr *serveMuxRoute)) *serveMuxRoute {
	smr, ok := mux.routes[pattern]
	if !ok {
		smr = &serveMuxRoute{
			allowedMethods:       []string{},
			handlerOptionsMaxAge: mux.config.handlerOptionsMaxAge,
			cors:                 mux.config.cors,
		}
		mux.routes[pattern] = smr

		createFn(smr)
	}
	return smr
}

func (mux *serveMux) addRoute(method, pattern string, handlerFn http.Handler) {
	validateHandler(handlerFn)

	pr := mux.patternRoute.join(pattern)

	mux.registerServeMuxRoute(pr.host+"/", func(smr *serveMuxRoute) {
		smr.addMethod(http.MethodOptions)
		mux.registerHandle(pr.host+"/", "HandlerNotFound", mux.mountMiddlewares(smr, http.HandlerFunc(mux.config.handlerNotFound)))
	})

	patternMethodNotAllowed := pr.mountMethodNotAllowed()
	smr := mux.registerServeMuxRoute(patternMethodNotAllowed, func(smr *serveMuxRoute) {
		smr.addMethod(http.MethodOptions)

		var handlerMethodNotAllowed http.Handler = http.HandlerFunc(mux.config.handlerMethodNotAllowed)
		if smr.cors == nil {
			handlerMethodNotAllowed = smr.middlewareMethodNotAllowed(handlerMethodNotAllowed)
		}

		mux.registerHandle(patternMethodNotAllowed, "HandlerMethodNotAllowed", mux.mountMiddlewares(smr, handlerMethodNotAllowed))
	})
	smr.addMethod(method)

	if smr.cors == nil && method == http.MethodOptions {
		handlerFn = smr.middlewareMethodNotAllowed(handlerFn)
	}

	mux.registerHandle(strings.TrimSpace(fmt.Sprintf("%s %s", method, pr.String())), "HandlerFn", mux.mountMiddlewares(smr, handlerFn))
}

// Connect registers a handler for the HTTP CONNECT method, under the current routing path plus the specified pattern.
func (mux *serveMux) Connect(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute(http.MethodConnect, pattern, handlerFn)
}

// Delete registers a handler for the HTTP DELETE method, under the current routing path plus the specified pattern.
func (mux *serveMux) Delete(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute(http.MethodDelete, pattern, handlerFn)
}

// Get registers a handler for the HTTP GET method, under the current routing path plus the specified pattern.
func (mux *serveMux) Get(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute(http.MethodGet, pattern, handlerFn)
}

// Head registers a handler for the HTTP HEAD method, under the current routing path plus the specified pattern.
func (mux *serveMux) Head(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute(http.MethodHead, pattern, handlerFn)
}

// Options registers a handler for the HTTP OPTIONS method, under the current routing path plus the specified pattern.
func (mux *serveMux) Options(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute(http.MethodOptions, pattern, handlerFn)
}

// Patch registers a handler for the HTTP PATCH method, under the current routing path plus the specified pattern.
func (mux *serveMux) Patch(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute(http.MethodPatch, pattern, handlerFn)
}

// Post registers a handler for the HTTP POST method, under the current routing path plus the specified pattern.
func (mux *serveMux) Post(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute(http.MethodPost, pattern, handlerFn)
}

// Put registers a handler for the HTTP PUT method, under the current routing path plus the specified pattern.
func (mux *serveMux) Put(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute(http.MethodPut, pattern, handlerFn)
}

// Trace registers a handler for the HTTP TRACE method, under the current routing path plus the specified pattern.
func (mux *serveMux) Trace(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute(http.MethodTrace, pattern, handlerFn)
}

// Method registers a handler for the custom HTTP method, under the current routing path plus the specified pattern.
func (mux *serveMux) Method(method, pattern string, handlerFn http.HandlerFunc) {
	if method == "" {
		panic(errors.New("method not specified"))
	}

	mux.addRoute(strings.ToUpper(method), pattern, handlerFn)
}

// ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL.
func (mux *serveMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.serveMux.ServeHTTP(w, r)
}
