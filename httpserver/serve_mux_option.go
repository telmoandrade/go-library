package httpserver

import "net/http"

func defaultHandlerNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 not found"))
}

func defaultHandlerMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 method not allowed"))
	}
}

// WithHandlerNotFound is an [OptionServeMux] that sets up a handler for not found routes.
// Inheriting the middleware stack.
func WithHandlerNotFound(handlerFn http.HandlerFunc) OptionServeMux {
	return func(mux *serveMux) {
		if handlerFn != nil {
			mux.config.handlerNotFound = handlerFn
		} else {
			mux.config.handlerNotFound = defaultHandlerNotFound
		}
	}
}

// WithHandlerMethodNotAllowed is an [OptionServeMux] that sets up a handler for not allowed methods routes.
// Inheriting the middleware stack.
//
// Default Behavior without using Cors:
//   - Automatically add the Allow in the response header.
//   - Automatically add the Cache-Control in the response header. See [WithHandlerOptionsMaxAge].
func WithHandlerMethodNotAllowed(handlerFn http.HandlerFunc) OptionServeMux {
	return func(mux *serveMux) {
		if handlerFn != nil {
			mux.config.handlerMethodNotAllowed = handlerFn
		} else {
			mux.config.handlerMethodNotAllowed = defaultHandlerMethodNotAllowed
		}
	}
}

// TODO refatorar nome
// WithHandlerOptionsMaxAge is an [OptionServeMux] that defines in seconds the maximum age for the Cache-Control header response in the options method handlers.
//
// Default:
//   - The default maximum age for the Cache-Control header response is 86400 seconds.
func WithHandlerOptionsMaxAge(seconds int) OptionServeMux {
	return func(mux *serveMux) {
		mux.config.handlerOptionsMaxAge = seconds
	}
}
