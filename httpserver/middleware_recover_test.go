package httpserver_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/telmoandrade/go-library/httpserver"
)

func TestMiddlewareRecover(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	tests := []struct {
		name   string
		args   bool
		status int
		want   int
	}{
		{
			name:   "no panic",
			status: http.StatusOK,
			want:   http.StatusOK,
		},
		{
			name: "panic",
			args: true,
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			r.Pattern = "GET /"
			w := httptest.NewRecorder()

			m := httpserver.MiddlewareRecover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.args {
					panic("fail")
				}
				w.WriteHeader(tt.status)
			}))

			m.ServeHTTP(w, r)

			if w.Code != tt.want {
				t.Errorf("Code() = %v, want %v", w.Code, tt.want)
			}
		})
	}
}
