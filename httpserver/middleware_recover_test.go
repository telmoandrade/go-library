package httpserver

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareRecover(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	type args struct {
		generatePanic bool
	}
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no panic",
			args: args{generatePanic: false},
			want: want{statusCode: http.StatusOK},
		},
		{
			name: "panic",
			args: args{generatePanic: true},
			want: want{statusCode: http.StatusInternalServerError},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			m := MiddlewareRecover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.args.generatePanic {
					panic("fail")
				}
				w.WriteHeader(http.StatusOK)
			}))

			m.ServeHTTP(w, r)

			if w.Code != tt.want.statusCode {
				t.Errorf("Code() = %v, want %v", w.Code, tt.want.statusCode)
			}
		})
	}
}
