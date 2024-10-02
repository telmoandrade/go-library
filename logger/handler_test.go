package logger

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "without option",
		},
		{
			name: "with option",
			args: args{
				opts: []Option{
					func(c *loggerHandler) {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLogger(tt.args.opts...)
			if got == nil {
				t.Errorf("NewLogger() = nil, want != nil")
			}
		})
	}
}

func TestNewHandler(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "without option",
		},
		{
			name: "with option",
			args: args{
				opts: []Option{
					func(c *loggerHandler) {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHandler(tt.args.opts...)
			if got == nil {
				t.Errorf("NewHandler() = nil, want != nil")
			}
		})
	}
}

func TestWithMinLevel(t *testing.T) {
	type args struct {
		l slog.Level
	}
	tests := []struct {
		name string
		args args
		want slog.Level
	}{
		{name: "Level debug", args: args{l: slog.LevelDebug}, want: slog.LevelDebug},
		{name: "Level info", args: args{l: slog.LevelInfo}, want: slog.LevelInfo},
		{name: "Level warn", args: args{l: slog.LevelWarn}, want: slog.LevelWarn},
		{name: "Level error", args: args{l: slog.LevelError}, want: slog.LevelError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lh := &loggerHandler{}
			WithMinLevel(tt.args.l)(lh)
			if got := lh.minLevel; got != tt.want {
				t.Errorf("WithMinLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithHandler(t *testing.T) {
	type args struct {
		handler slog.Handler
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "nil handler",
		},
		{
			name: "handler",
			args: args{
				handler: slog.NewTextHandler(os.Stdout, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := NewMockSlogHandler(ctrl)

			lh := &loggerHandler{
				handler: mock,
			}
			WithHandler(tt.args.handler)(lh)
			if lh.handler == nil {
				t.Errorf("WithHandler() = %v, want %v", lh.handler, nil)
			}
		})
	}
}

func TestWithMaxLevelAddSource(t *testing.T) {
	type args struct {
		l slog.Level
	}
	tests := []struct {
		name string
		args args
		want slog.Level
	}{
		{name: "Level debug", args: args{l: slog.LevelDebug}, want: slog.LevelDebug},
		{name: "Level info", args: args{l: slog.LevelInfo}, want: slog.LevelInfo},
		{name: "Level warn", args: args{l: slog.LevelWarn}, want: slog.LevelWarn},
		{name: "Level error", args: args{l: slog.LevelError}, want: slog.LevelError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lh := &loggerHandler{}
			WithMaxLevelAddSource(tt.args.l)(lh)
			if got := lh.maxLevelAddSource; got != tt.want {
				t.Errorf("WithMaxLevelAddSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loggerHandler_Enabled(t *testing.T) {
	type fields struct {
		minLevel slog.Level
	}
	type args struct {
		ctx context.Context
		l   slog.Level
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "enable true, without context log level",
			fields: fields{
				minLevel: slog.LevelInfo,
			},
			args: args{
				ctx: context.Background(),
				l:   slog.LevelInfo,
			},
			want: true,
		},
		{
			name: "enable false, without context log level",
			fields: fields{
				minLevel: slog.LevelInfo,
			},
			args: args{
				ctx: context.Background(),
				l:   slog.LevelDebug,
			},
			want: false,
		},
		{
			name: "enable true, with context log level",
			fields: fields{
				minLevel: slog.LevelInfo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), ContextMinLevel, slog.LevelDebug),
				l:   slog.LevelDebug,
			},
			want: true,
		},
		{
			name: "enable false, with context log level",
			fields: fields{
				minLevel: slog.LevelInfo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), ContextMinLevel, slog.LevelInfo),
				l:   slog.LevelDebug,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lh := &loggerHandler{
				minLevel: tt.fields.minLevel,
			}
			if got := lh.Enabled(tt.args.ctx, tt.args.l); got != tt.want {
				t.Errorf("loggerHandler.Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loggerHandler_Handle(t *testing.T) {
	type fields struct {
		maxLevelAddSource slog.Level
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				maxLevelAddSource: slog.LevelDebug,
			},
			args: args{ctx: context.Background()},
		},
		{
			name: "with add source",
			fields: fields{
				maxLevelAddSource: slog.LevelInfo,
			},
			args: args{ctx: context.Background()},
		},
		{
			name: "with log nil",
			fields: fields{
				maxLevelAddSource: slog.LevelInfo,
			},
			args: args{ctx: context.WithValue(context.Background(), ContextLogID, uuid.Nil)},
		},
		{
			name: "with log",
			fields: fields{
				maxLevelAddSource: slog.LevelInfo,
			},
			args: args{ctx: context.WithValue(context.Background(), ContextLogID, uuid.Must(uuid.NewV7()))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := NewMockSlogHandler(ctrl)
			mock.EXPECT().Handle(gomock.Any(), gomock.Any()).Times(1)

			lh := &loggerHandler{
				handler:           mock,
				maxLevelAddSource: tt.fields.maxLevelAddSource,
			}

			r := slog.Record{}
			r.Level = slog.LevelInfo

			err := lh.Handle(tt.args.ctx, r)
			if err != nil {
				t.Errorf("invald error on loggerHandler.Handle() = %v", err)
			}
		})
	}
}

func Test_loggerHandler_WithAttrs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSlogHandler(ctrl)
	mock.EXPECT().WithAttrs(gomock.Any()).Times(1)

	lh := &loggerHandler{
		handler: mock,
	}
	lh.WithAttrs([]slog.Attr{})
}

func Test_loggerHandler_WithGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSlogHandler(ctrl)
	mock.EXPECT().WithGroup(gomock.Any()).Times(1)

	lh := &loggerHandler{
		handler: mock,
	}
	lh.WithGroup("group")
}
