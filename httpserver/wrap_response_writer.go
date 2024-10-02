package httpserver

import (
	"net/http"
	"sync"
)

type (
	// wrapResponseWriter wraps a [http.ResponseWriter] to track the number of bytes written and the status code.
	wrapResponseWriter struct {
		http.ResponseWriter
		mu          sync.Mutex
		wroteHeader bool
		code        int
		bytes       int64
	}
)

var _ http.ResponseWriter = &wrapResponseWriter{}

// WriteHeader write status code to [http.ResponseWriter] and store the status code
func (wrw *wrapResponseWriter) WriteHeader(code int) {
	if !wrw.wroteHeader {
		wrw.wroteHeader = true
		wrw.code = code
		wrw.ResponseWriter.WriteHeader(code)
	}
}

// Write data to [http.ResponseWriter] and store the number of bytes written
func (w *wrapResponseWriter) Write(b []byte) (int, error) {
	w.WriteHeader(http.StatusOK)
	n, err := w.ResponseWriter.Write(b)

	w.mu.Lock()
	w.bytes += int64(n)
	w.mu.Unlock()

	return n, err
}
