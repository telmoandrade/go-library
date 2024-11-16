package httpserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_defaultHandlerNotFound(t *testing.T) {
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "call",
			want: want{
				statusCode: http.StatusNotFound,
				body:       "404 not found",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			defaultHandlerNotFound(w, r)

			if w.Code != tt.want.statusCode {
				t.Errorf("Code() = %v, want %v", w.Code, tt.want.statusCode)
			}

			respBody, err := io.ReadAll(w.Body)
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			if string(respBody) != tt.want.body {
				t.Errorf("Body = %v, want %v", string(respBody), tt.want.body)
			}
		})
	}
}

func Test_defaultHandlerMethodNotAllowed(t *testing.T) {
	type args struct {
		method string
	}
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "http.MethodGet",
			args: args{
				method: http.MethodGet,
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				body:       "405 method not allowed",
			},
		},
		{
			name: "http.MethodOptions",
			args: args{
				method: http.MethodOptions,
			},
			want: want{
				statusCode: http.StatusNoContent,
				body:       "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.args.method, "/", nil)
			w := httptest.NewRecorder()

			defaultHandlerMethodNotAllowed(w, r)

			if w.Code != tt.want.statusCode {
				t.Errorf("Code() = %v, want %v", w.Code, tt.want.statusCode)
			}

			respBody, err := io.ReadAll(w.Body)
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			if string(respBody) != tt.want.body {
				t.Errorf("Body = %v, want %v", string(respBody), tt.want.body)
			}
		})
	}
}

func TestWithHandlerNotFound(t *testing.T) {
	type args struct {
		handlerFn http.HandlerFunc
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "empty",
		},
		{
			name: "custom",
			args: args{
				handlerFn: func(w http.ResponseWriter, r *http.Request) {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux(
				WithHandlerNotFound(tt.args.handlerFn),
			)
			mux := muxInterface.(*serveMux)
			h := mux.config.handlerNotFound

			if h == nil {
				t.Errorf("WithHandlerNotFound invalid")
			}
		})
	}
}

func TestWithHandlerMethodNotAllowed(t *testing.T) {
	type args struct {
		handlerFn http.HandlerFunc
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "empty",
		},
		{
			name: "custom",
			args: args{
				handlerFn: func(w http.ResponseWriter, r *http.Request) {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux(
				WithHandlerMethodNotAllowed(tt.args.handlerFn),
			)
			mux := muxInterface.(*serveMux)
			h := mux.config.handlerMethodNotAllowed

			if h == nil {
				t.Errorf("WithHandlerMethodNotAllowed invalid")
			}
		})
	}
}

func TestWithHandlerOptionsMaxAge(t *testing.T) {
	type args struct {
		seconds int
	}
	type want struct {
		seconds int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty",
		},
		{
			name: "86400 seconds",
			args: args{seconds: 86400},
			want: want{seconds: 86400},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux(
				WithHandlerOptionsMaxAge(tt.args.seconds),
			)
			mux := muxInterface.(*serveMux)
			if got := mux.config.handlerOptionsMaxAge; got != tt.want.seconds {
				t.Errorf("seconds = %v, want %v", got, tt.want.seconds)
			}
		})
	}
}
