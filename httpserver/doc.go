// Package httpserver implements functions to extend the use of the http.Server standard library.
//
// # ServeMux
//
// [ServeMux] is an extension with new methods for [http.ServeMux] is a multiplexer for HTTP requests.
// use [NewServeMux] allocates and returns a new ServeMux.
//
// Stores routing path to register middlewares and handlers.
//
// Methods for helping routing path and registering middlewares:
//   - [ServeMux.Use] appends one or more middlewares onto the Router stack.
//   - [ServeMux.With] adds inline middlewares for registers the handler.
//   - [ServeMux.Group] adds a new inline-Router along the current routing path + /pattern, with middleware stack.
//   - [ServeMux.Route] mounts a sub-Router along the current routing path + /pattern, with middleware stack.
//   - [ServeMux.Mount] attaches another http.ServeMux along a /pattern/*.
//
// WARNING: Avoid using [ServeMux.Mount], because it is slower to resolve the multiplexer for HTTP requests.
//
// Methods for registering handler  along the current routing path + /pattern:
//   - [ServeMux.Connect] with http method CONNECT
//   - [ServeMux.Delete] with http method DELETE
//   - [ServeMux.Get] with http method GET
//   - [ServeMux.Head] with http method HEAD
//   - [ServeMux.Options] with http method OPTIONS
//   - [ServeMux.Patch] with http method PATCH
//   - [ServeMux.Post] with http method POST
//   - [ServeMux.Put] with http method PUT
//   - [ServeMux.Trace] with http method TRACE
//   - [ServeMux.Method] with custom http method by parameter
package httpserver
