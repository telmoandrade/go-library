package httpserver

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_realIPExtractHeader(t *testing.T) {
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
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, "/", nil)
			r.Header.Add(tt.args.key, tt.args.value)

			if got := realIPExtractHeader(r); got != tt.want {
				t.Errorf("realIPExtractHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_realIP(t *testing.T) {
	type args struct {
		key        string
		value      string
		remoteAddr string
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
			name: "no header",
			args: args{
				remoteAddr: "127.0.0.1:12345",
			},
			want: "127.0.0.1",
		},
		{
			name: "invalid remote addr",
			args: args{
				remoteAddr: "invalid",
			},
			want: "invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, "/", nil)
			r.Header.Add(tt.args.key, tt.args.value)
			r.RemoteAddr = tt.args.remoteAddr

			if got := realIP(r); got != tt.want {
				t.Errorf("realIP() = %v, want %v", got, tt.want)
			}
		})
	}
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

	type args struct {
		headerKey   string
		headerValue string
		status      int
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
			name: "InternalServerError",
			args: args{
				status: http.StatusInternalServerError,
			},
			want: want{statusCode: http.StatusInternalServerError},
		},
		{
			name: "NotFound",
			args: args{
				status: http.StatusNotFound,
			},
			want: want{statusCode: http.StatusNotFound},
		},
		{
			name: "Ok",
			args: args{
				status: http.StatusOK,
			},
			want: want{statusCode: http.StatusOK},
		},
		{
			name: "Ok with header",
			args: args{
				status:      http.StatusOK,
				headerKey:   "X-Logger-Level",
				headerValue: "info",
			},
			want: want{statusCode: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Pattern = "GET /"
			r.Header.Add(tt.args.headerKey, tt.args.headerValue)
			w := httptest.NewRecorder()

			m := MiddlewareLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.status)
			}))
			m.ServeHTTP(w, r)

			if w.Code != tt.want.statusCode {
				t.Errorf("Code() = %v, want %v", w.Code, tt.want.statusCode)
			}
		})
	}
}
