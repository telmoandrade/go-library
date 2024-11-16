package httpserver

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_wrapResponseWriter_WriteHeader(t *testing.T) {
	type args struct {
		code1 int
		code2 int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "code http.StatusOK, http.StatusOK",
			args: args{
				code1: http.StatusOK,
				code2: http.StatusOK,
			},
			want: http.StatusOK,
		},
		{
			name: "code http.StatusOK, http.StatusBadRequest",
			args: args{
				code1: http.StatusOK,
				code2: http.StatusBadRequest,
			},
			want: http.StatusOK,
		},
		{
			name: "code http.StatusBadRequest, http.StatusOK",
			args: args{
				code1: http.StatusBadRequest,
				code2: http.StatusOK,
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &httptest.ResponseRecorder{
				Body: new(bytes.Buffer),
			}
			wrw := &wrapResponseWriter{
				ResponseWriter: w,
			}

			wrw.WriteHeader(tt.args.code1)
			if wrw.code != tt.want {
				t.Errorf("WriteHeader(%v) = %v, want %v", tt.args.code1, wrw.code, tt.want)
			}

			wrw.WriteHeader(tt.args.code2)
			if wrw.code != tt.want {
				t.Errorf("WriteHeader(%v) = %v, want %v", tt.args.code2, wrw.code, tt.want)
			}
		})
	}
}

func Test_wrapResponseWriter_Write(t *testing.T) {
	tests := []struct {
		name    string
		args    []byte
		want    int
		wantErr bool
	}{
		{
			name:    "no data",
			want:    0,
			wantErr: false,
		},
		{
			name:    "data",
			args:    []byte("data"),
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &httptest.ResponseRecorder{
				Body: new(bytes.Buffer),
			}
			wrw := &wrapResponseWriter{
				ResponseWriter: w,
			}

			var wantBytes int64

			got, err := wrw.Write(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wrapResponseWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("wrapResponseWriter.Write() = %v, want %v", got, tt.want)
			}
			wantBytes += int64(got)

			got, err = wrw.Write(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wrapResponseWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("wrapResponseWriter.Write() = %v, want %v", got, tt.want)
			}
			wantBytes += int64(got)

			if wrw.bytes != wantBytes {
				t.Errorf("wrapResponseWriter.bytes = %v, want %v", wrw.bytes, wantBytes)
			}
		})
	}
}
