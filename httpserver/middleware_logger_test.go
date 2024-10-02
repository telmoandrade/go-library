package httpserver

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_realIP(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "True-Client-IP",
			args: args{
				key:   "True-Client-IP",
				value: "100.100.100.100",
			},
			want: "100.100.100.100",
		},
		{
			name: "X-Real-IP",
			args: args{
				key:   "X-Real-IP",
				value: "100.100.100.100",
			},
			want: "100.100.100.100",
		},
		{
			name: "X-Forwarded-For 1 ip",
			args: args{
				key:   "X-Forwarded-For",
				value: "100.100.100.100",
			},
			want: "100.100.100.100",
		},
		{
			name: "X-Forwarded-For 2 ip",
			args: args{
				key:   "X-Forwarded-For",
				value: "100.100.100.100, 200.200.200.200",
			},
			want: "100.100.100.100",
		},
		{
			name: "no header",
			want: "127.0.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest("GET", "/", nil)
			r.Header.Add(tt.args.key, tt.args.value)
			r.RemoteAddr = "127.0.0.1:12345"

			if got := realIP(r); got != tt.want {
				t.Errorf("realIP() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("invalid remote addr", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "/", nil)
		r.RemoteAddr = "invalid"

		realIP(r)
		if r.RemoteAddr != "invalid" {
			t.Errorf("realIP() = %v, want %v", r.RemoteAddr, "invalid")
		}
	})
}

func Test_sinceRound(t *testing.T) {
	tests := []struct {
		name string
		args time.Duration
		want time.Duration
	}{
		{
			name: "Nanosecond",
			args: time.Nanosecond * 999,
			want: time.Nanosecond * 999,
		},
		{
			name: "Microsecond - remove Nanosecond",
			args: time.Microsecond + time.Nanosecond,
			want: time.Microsecond,
		},
		{
			name: "Millisecond",
			args: time.Millisecond + time.Nanosecond,
			want: time.Millisecond,
		},
		{
			name: "Second",
			args: time.Second + time.Nanosecond,
			want: time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sinceRound(tt.args); got != tt.want {
				t.Errorf("sinceRound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiddlewareLogging(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	tests := []struct {
		name   string
		status int
		want   int
	}{
		{
			name:   "InternalServerError",
			status: http.StatusInternalServerError,
			want:   http.StatusInternalServerError,
		},
		{
			name:   "NotFound",
			status: http.StatusNotFound,
			want:   http.StatusNotFound,
		},
		{
			name:   "Ok",
			status: http.StatusOK,
			want:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Add("X-Logger-Level", "debug")
			r.Pattern = "GET /"
			w := httptest.NewRecorder()

			m := MiddlewareLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
			}))
			m.ServeHTTP(w, r)

			if w.Code != tt.want {
				t.Errorf("Code() = %v, want %v", w.Code, tt.want)
			}
		})
	}
}
