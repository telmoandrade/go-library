package graceful

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testing"

	gomock "go.uber.org/mock/gomock"
)

func Test_gracefulServerHttpStart(t *testing.T) {
	slog.SetLogLoggerLevel(slog.Level(16))

	errMock := errors.New("error")

	type args struct {
		error                 error
		callListenAndServe    int
		callListenAndServeTLS int
		certFile              string
		keyFile               string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "start without error",
			args: args{
				error:                 nil,
				callListenAndServe:    1,
				callListenAndServeTLS: 0,
			},
			want: nil,
		},
		{
			name: "start with error",
			args: args{
				error:                 errMock,
				callListenAndServe:    1,
				callListenAndServeTLS: 0,
			},
			want: errMock,
		},
		{
			name: "start tls without error",
			args: args{
				error:                 nil,
				callListenAndServe:    0,
				callListenAndServeTLS: 1,
				certFile:              "certFile",
				keyFile:               "keyFile",
			},
			want: nil,
		},
		{
			name: "start tls with error",
			args: args{
				error:                 errMock,
				callListenAndServe:    0,
				callListenAndServeTLS: 1,
				certFile:              "certFile",
				keyFile:               "keyFile",
			},
			want: errMock,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := NewMockhttpServer(ctrl)
			mock.EXPECT().ListenAndServe().Return(tt.args.error).Times(tt.args.callListenAndServe)
			mock.EXPECT().ListenAndServeTLS(gomock.Any(), gomock.Any()).Return(tt.args.error).Times(tt.args.callListenAndServeTLS)

			gs := NewGracefulServerHttp(
				mock,
				WithTLS(tt.args.certFile, tt.args.keyFile),
			).(*gracefulServerHttp)
			err := gracefulServerHttpStart(gs, mock)()
			if err != tt.want {
				t.Errorf("gracefulServerHttpStart() = %v, want %v", err, tt.want)
			}
		})
	}
}

func Test_gracefulServerHttpStop(t *testing.T) {
	slog.SetLogLoggerLevel(slog.Level(16))

	tests := []struct {
		name string
		args error
	}{
		{
			name: "stop without error",
			args: nil,
		},
		{
			name: "stop with error",
			args: errors.New("error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := NewMockhttpServer(ctrl)
			mock.EXPECT().Shutdown(gomock.Any()).Return(tt.args).Times(1)

			gs := NewGracefulServerHttp(mock).(*gracefulServerHttp)
			gracefulServerHttpStop(gs, mock)(context.Background())
		})
	}
}

func Test_gracefulServerHttpForceStop(t *testing.T) {
	slog.SetLogLoggerLevel(slog.Level(16))

	tests := []struct {
		name string
		args error
	}{
		{
			name: "force stop without error",
			args: nil,
		},
		{
			name: "force stop with error",
			args: errors.New("error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := NewMockhttpServer(ctrl)
			mock.EXPECT().Close().Return(tt.args).Times(1)

			gs := NewGracefulServerHttp(mock).(*gracefulServerHttp)
			gracefulServerHttpForceStop(gs, mock)()
		})
	}
}

func TestNewGracefulServerHttp(t *testing.T) {
	type args struct {
		opts       []OptionGracefulServerHttp
		httpServer httpServer
	}
	tests := []struct {
		name       string
		args       args
		wantIsNull bool
	}{
		{
			name: "nil http.Server",
			args: args{
				httpServer: nil,
			},
			wantIsNull: true,
		},
		{
			name: "without option",
			args: args{
				httpServer: &MockhttpServer{},
			},
			wantIsNull: false,
		},
		{
			name: "with option",
			args: args{
				httpServer: &MockhttpServer{},
				opts: []OptionGracefulServerHttp{
					func(c *gracefulServerHttp) {},
				},
			},
			wantIsNull: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGracefulServerHttp(tt.args.httpServer, tt.args.opts...)
			if tt.wantIsNull && got != nil {
				t.Errorf("NewControl() = %v, want %v", got, nil)
			}
			if !tt.wantIsNull && got == nil {
				t.Errorf("NewControl() = %v, want %v", got, &gracefulServerHttp{})
			}
		})
	}
}

func TestWithTLS(t *testing.T) {
	type args struct {
		certFile string
		keyFile  string
	}
	tests := []struct {
		name string
		args args
		want args
	}{
		{
			name: "empty values",
			want: args{
				certFile: "",
				keyFile:  "",
			},
		},
		{
			name: "empty certFile",
			args: args{
				certFile: "",
				keyFile:  "keyFile",
			},
			want: args{
				certFile: "",
				keyFile:  "",
			},
		},
		{
			name: "empty keyFile",
			args: args{
				certFile: "certFile",
				keyFile:  "",
			},
			want: args{
				certFile: "",
				keyFile:  "",
			},
		},
		{
			name: "ok",
			args: args{
				certFile: "certFile",
				keyFile:  "keyFile",
			},
			want: args{
				certFile: "certFile",
				keyFile:  "keyFile",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := NewGracefulServerHttp(&MockhttpServer{}, WithTLS(tt.args.certFile, tt.args.keyFile))
			got := gs.(*gracefulServerHttp)

			if got.certFile != tt.want.certFile {
				t.Errorf("certFile = %v, want %v", got.certFile, tt.want.certFile)
			}
			if got.keyFile != tt.want.keyFile {
				t.Errorf("keyFile = %v, want %v", got.keyFile, tt.want.keyFile)
			}
		})
	}
}

func TestWithSlogAttrs(t *testing.T) {
	tests := []struct {
		name string
		args []slog.Attr
		want string
	}{
		{
			name: "empty",
			want: "[]",
		},
		{
			name: "nil",
			args: nil,
			want: "[]",
		},
		{
			name: "1 slog.Attr",
			args: []slog.Attr{slog.Bool("v", true)},
			want: "[v=true]",
		},
		{
			name: "2 slog.Attr",
			args: []slog.Attr{slog.Bool("v", true), slog.Bool("f", true)},
			want: "[v=true f=true]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &gracefulServerHttp{}
			WithSlogAttrs(tt.args...)(gs)

			got := fmt.Sprintf("%v", gs.attrs)
			if got != tt.want {
				t.Errorf("attrs = %v, want %v", gs.attrs, tt.want)
			}
		})
	}
}
