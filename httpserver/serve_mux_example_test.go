package httpserver_test

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/telmoandrade/go-library/httpserver"
)

func handlerHello(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("hello")); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}

type ctxKey struct {
	name string
}

func middlewarePathValue(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ctx := context.WithValue(r.Context(), ctxKey{"id"}, id)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func handlerGetUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(ctxKey{"id"}).(string)

	fmt.Println(id)
}

func handlerPutUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(ctxKey{"id"}).(string)

	fmt.Println(id)
}

func handler(w http.ResponseWriter, r *http.Request) {}

func ExampleNewServeMux() {
	mux := httpserver.NewServeMux()
	mux.Get("/hello", handlerHello)

	s := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	fmt.Print(s != nil)
	// Output: true
}

func ExampleServeMux_Use() {
	mux := httpserver.NewServeMux()
	mux.Use(middlewarePathValue)
	mux.Get("/user/{id}", handlerGetUser)
	mux.Put("/user/{id}", handlerPutUser)
	// Output:
}

func ExampleServeMux_With() {
	mux := httpserver.NewServeMux()
	mux.With(middlewarePathValue).Get("/user/{id}", handlerGetUser)
	// Output:
}

func ExampleServeMux_Group() {
	mux := httpserver.NewServeMux()
	muxUser := mux.Group("/user")

	muxUser.Use(middlewarePathValue)
	muxUser.Get("/{id}", handlerGetUser)
	muxUser.Put("/{id}", handlerPutUser)

	s := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	fmt.Print(s != nil)
	// Output: true
}

func ExampleServeMux_Route() {
	mux := httpserver.NewServeMux()
	mux.Route("/user", func(muxUser httpserver.Router) {
		muxUser.Use(middlewarePathValue)
		muxUser.Get("/{id}", handlerGetUser)
		muxUser.Put("/{id}", handlerPutUser)
	})

	s := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	fmt.Print(s != nil)
	// Output: true
}

func ExampleServeMux_Mount() {
	muxUser := httpserver.NewServeMux()
	muxUser.Use(middlewarePathValue)
	muxUser.Get("/{id}", handlerGetUser)
	muxUser.Put("/{id}", handlerPutUser)

	mux := httpserver.NewServeMux()
	mux.Mount("/user", muxUser)

	s := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	fmt.Print(s != nil)
	// Output: true
}

func ExampleServeMux_Connect() {
	mux := httpserver.NewServeMux()
	mux.Connect("/pattern", handler)
	// Output:
}

func ExampleServeMux_Delete() {
	mux := httpserver.NewServeMux()
	mux.Delete("/pattern", handler)
	// Output:
}

func ExampleServeMux_Get() {
	mux := httpserver.NewServeMux()
	mux.Get("/pattern", handler)
	// Output:
}

func ExampleServeMux_Head() {
	mux := httpserver.NewServeMux()
	mux.Head("/pattern", handler)
	// Output:
}

func ExampleServeMux_Options() {
	mux := httpserver.NewServeMux()
	mux.Options("/pattern", handler)
	// Output:
}

func ExampleServeMux_Patch() {
	mux := httpserver.NewServeMux()
	mux.Patch("/pattern", handler)
	// Output:
}

func ExampleServeMux_Post() {
	mux := httpserver.NewServeMux()
	mux.Post("/pattern", handler)
	// Output:
}

func ExampleServeMux_Put() {
	mux := httpserver.NewServeMux()
	mux.Put("/pattern", handler)
	// Output:
}

func ExampleServeMux_Trace() {
	mux := httpserver.NewServeMux()
	mux.Trace("/pattern", handler)
	// Output:
}

func ExampleServeMux_Method() {
	mux := httpserver.NewServeMux()
	mux.Method("CUSTOM", "/pattern", handler)
	// Output:
}
