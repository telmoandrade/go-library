package httpserver

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

type (
	// Handle allows registers the handler for the given pattern
	Handle interface {
		// Connect registers the handler along the current routing path + /pattern with http method CONNECT. If the given pattern conflicts with one that is already registered, Handle panics.
		Connect(pattern string, handlerFn http.HandlerFunc)
		// Delete registers the handler along the current routing path + /pattern with http method DELETE. If the given pattern conflicts with one that is already registered, Handle panics.
		Delete(pattern string, handlerFn http.HandlerFunc)
		// Get registers the handler along the current routing path + /pattern with http method GET. If the given pattern conflicts with one that is already registered, Handle panics.
		Get(pattern string, handlerFn http.HandlerFunc)
		// Head registers the handler along the current routing path + /pattern with http method HEAD. If the given pattern conflicts with one that is already registered, Handle panics.
		Head(pattern string, handlerFn http.HandlerFunc)
		// Options registers the handler along the current routing path + /pattern with http method OPTIONS. If the given pattern conflicts with one that is already registered, Handle panics.
		Options(pattern string, handlerFn http.HandlerFunc)
		// Patch registers the handler along the current routing path + /pattern with http method PATCH. If the given pattern conflicts with one that is already registered, Handle panics.
		Patch(pattern string, handlerFn http.HandlerFunc)
		// Post registers the handler along the current routing path + /pattern with http method POST. If the given pattern conflicts with one that is already registered, Handle panics.
		Post(pattern string, handlerFn http.HandlerFunc)
		// Put registers the handler along the current routing path + /pattern with http method PUT. If the given pattern conflicts with one that is already registered, Handle panics.
		Put(pattern string, handlerFn http.HandlerFunc)
		// Trace registers the handler along the current routing path + /pattern with http method TRACE. If the given pattern conflicts with one that is already registered, Handle panics.
		Trace(pattern string, handlerFn http.HandlerFunc)
		// Method registers the handler for the given along the current routing path + /pattern with custom http method by parameter. If the given pattern conflicts with one that is already registered, Handle panics.
		Method(method, pattern string, handlerFn http.HandlerFunc)
	}

	// Router stores routing path to register middlewares and handlers.
	Router interface {
		Handle
		http.Handler

		// Use appends one or more middlewares onto the [Router] stack.
		Use(middlewares ...func(http.Handler) http.Handler)
		// With adds inline middlewares for registers the handler.
		With(middlewares ...func(http.Handler) http.Handler) Handle
		// Group adds a new inline-[Router] along the current routing path + /pattern, with middleware stack.
		Group(pattern string) Router
		// Route mounts a sub-[Router] along the current routing path + /pattern, with middleware stack.
		Route(pattern string, fn func(sub Router))
		// Mount attaches another [http.ServeMux] along a /pattern/*.
		//
		// WARNING: Avoid using, because it is slower to resolve the multiplexer for HTTP requests.
		Mount(pattern string, router Router)
	}

	// ServeMux is an extension with new methods for [http.ServeMux] is a multiplexer for HTTP requests
	//
	// Stores routing path to register middlewares and handlers.
	ServeMux struct {
		*http.ServeMux
		middlewares []func(http.Handler) http.Handler
		pattern     string
	}
)

// NewServeMux allocates and returns a new [ServeMux].
func NewServeMux() *ServeMux {
	return &ServeMux{
		ServeMux: http.NewServeMux(),
	}
}

// Use appends one or more middlewares onto the [Router] stack.
func (mux *ServeMux) Use(middlewares ...func(http.Handler) http.Handler) {
	mux.middlewares = append(mux.middlewares, middlewares...)
}

// With adds inline middlewares for registers the handler.
func (mux *ServeMux) With(middlewares ...func(http.Handler) http.Handler) Handle {
	return &ServeMux{
		pattern:     mux.pattern,
		middlewares: append(mux.middlewares, middlewares...),
		ServeMux:    mux.ServeMux,
	}
}

// Group adds a new inline-[Router] along the current routing path + /pattern, with middleware stack.
func (mux *ServeMux) Group(pattern string) Router {
	return &ServeMux{
		pattern:     path.Join(mux.pattern, pattern),
		middlewares: mux.middlewares,
		ServeMux:    mux.ServeMux,
	}
}

// Route mounts a sub-[Router] along the current routing path + /pattern, with middleware stack.
func (mux *ServeMux) Route(pattern string, fn func(sub Router)) {
	subRouter := &ServeMux{
		pattern:     path.Join(mux.pattern, pattern),
		middlewares: mux.middlewares,
		ServeMux:    mux.ServeMux,
	}
	fn(subRouter)

}

// Mount attaches another [http.ServeMux] along a /pattern/*.
//
// WARNING: Avoid using, because it is slower to resolve the multiplexer for HTTP requests.
func (mux *ServeMux) Mount(pattern string, r Router) {
	targetPattern := path.Join(mux.pattern, pattern)

	mux.Handle(targetPattern+"/", http.StripPrefix(targetPattern, r))
}

func (mux *ServeMux) addRoute(method string, pattern string, handlerFn http.Handler) {
	targetPattern := path.Join(mux.pattern, pattern)
	targetHandler := handlerFn

	for i := len(mux.middlewares) - 1; i >= 0; i-- {
		targetHandler = mux.middlewares[i](targetHandler)
	}

	mux.Handle(fmt.Sprintf("%s %s", method, targetPattern), targetHandler)
}

// Connect registers the handler along the current routing path + /pattern with http method CONNECT. If the given pattern conflicts with one that is already registered, Handle panics.
func (mux *ServeMux) Connect(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("CONNECT", pattern, handlerFn)
}

// Delete registers the handler along the current routing path + /pattern with http method DELETE. If the given pattern conflicts with one that is already registered, Handle panics.
func (mux *ServeMux) Delete(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("DELETE", pattern, handlerFn)
}

// Get registers the handler along the current routing path + /pattern with http method GET. If the given pattern conflicts with one that is already registered, Handle panics.
func (mux *ServeMux) Get(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("GET", pattern, handlerFn)
}

// Head registers the handler along the current routing path + /pattern with http method HEAD. If the given pattern conflicts with one that is already registered, Handle panics.
func (mux *ServeMux) Head(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("HEAD", pattern, handlerFn)
}

// Options registers the handler along the current routing path + /pattern with http method OPTIONS. If the given pattern conflicts with one that is already registered, Handle panics.
func (mux *ServeMux) Options(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("OPTIONS", pattern, handlerFn)
}

// Patch registers the handler along the current routing path + /pattern with http method PATCH. If the given pattern conflicts with one that is already registered, Handle panics.
func (mux *ServeMux) Patch(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("PATCH", pattern, handlerFn)
}

// Post registers the handler along the current routing path + /pattern with http method POST. If the given pattern conflicts with one that is already registered, Handle panics.
func (mux *ServeMux) Post(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("POST", pattern, handlerFn)
}

// Put registers the handler along the current routing path + /pattern with http method PUT. If the given pattern conflicts with one that is already registered, Handle panics.
func (mux *ServeMux) Put(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("PUT", pattern, handlerFn)
}

// Trace registers the handler along the current routing path + /pattern with http method TRACE. If the given pattern conflicts with one that is already registered, Handle panics.
func (mux *ServeMux) Trace(pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute("TRACE", pattern, handlerFn)
}

// Method registers the handler for the given along the current routing path + /pattern with custom http method by parameter. If the given pattern conflicts with one that is already registered, Handle panics.
func (mux *ServeMux) Method(method, pattern string, handlerFn http.HandlerFunc) {
	mux.addRoute(strings.ToUpper(method), pattern, handlerFn)
}
