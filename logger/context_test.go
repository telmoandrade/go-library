package logger_test

import (
	"context"
	"log/slog"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/telmoandrade/go-library/logger"
)

func TestLogIDFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want uuid.UUID
	}{
		{
			name: "not define",
			args: args{
				ctx: context.Background(),
			},
			want: uuid.Nil,
		},
		{
			name: "define",
			args: args{
				ctx: context.WithValue(context.Background(), logger.ContextLogID, uuid.Must(uuid.Parse("11111111-1111-1111-1111-111111111111"))),
			},
			want: uuid.Must(uuid.Parse("11111111-1111-1111-1111-111111111111")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logger.LogIDFromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogIDFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithContextLogID(t *testing.T) {
	type args struct {
		u uuid.UUID
	}
	tests := []struct {
		name string
		args args
		want uuid.UUID
	}{
		{
			name: "not define",
			want: uuid.Nil,
		},
		{
			name: "define",
			args: args{
				u: uuid.Must(uuid.Parse("11111111-1111-1111-1111-111111111111")),
			},
			want: uuid.Must(uuid.Parse("11111111-1111-1111-1111-111111111111")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := logger.WithContextLogID(context.Background(), tt.args.u)
			got, ok := ctx.Value(logger.ContextLogID).(uuid.UUID)
			if !ok {
				t.Fatal("not ok")
			}

			if got != tt.want {
				t.Errorf("WithContextLogID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithContextMinLevel(t *testing.T) {
	type args struct {
		l string
	}
	tests := []struct {
		name string
		args args
		want slog.Level
	}{
		{name: "debug", args: args{l: "debug"}, want: slog.LevelDebug},
		{name: "info", args: args{l: "info"}, want: slog.LevelInfo},
		{name: "warn", args: args{l: "warn"}, want: slog.LevelWarn},
		{name: "error", args: args{l: "error"}, want: slog.LevelError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := logger.WithContextMinLevel(context.Background(), tt.args.l)
			got, ok := ctx.Value(logger.ContextMinLevel).(slog.Level)
			if !ok {
				t.Fatal("not ok")
			}
			if got != tt.want {
				t.Errorf("WithContextMinLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
