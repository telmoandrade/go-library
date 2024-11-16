package logger

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/google/uuid"
)

func TestLogId(t *testing.T) {
	type args struct {
		propagation string
		ctx         context.Context
	}
	tests := []struct {
		name string
		args args
		want uuid.UUID
	}{
		{
			name: "empty",
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "propagation",
			args: args{
				propagation: "11111111-1111-1111-1111-111111111111",
				ctx:         context.Background(),
			},
			want: uuid.Must(uuid.Parse("11111111-1111-1111-1111-111111111111")),
		},
		{
			name: "log exists",
			args: args{
				ctx: context.WithValue(context.Background(), ContextLogID, uuid.Must(uuid.Parse("11111111-1111-1111-1111-111111111111"))),
			},
			want: uuid.Must(uuid.Parse("11111111-1111-1111-1111-111111111111")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, u := LogId(tt.args.ctx, tt.args.propagation)

			u2, ok := ctx.Value(ContextLogID).(uuid.UUID)
			if !ok {
				t.Fatal("invalid logger.ContextLogID")
			}

			if u != u2 {
				t.Errorf("LogId() = %v, want %v", u, u2)
			}
			if u == uuid.Nil {
				t.Errorf("LogId() = %v, want %v", u, tt.want)
			}
			if tt.args.propagation != "" && u != tt.want {
				t.Errorf("LogId() = %v, want %v", u, tt.want)
			}
		})
	}
}

func TestMinLevel(t *testing.T) {
	type args struct {
		l string
	}
	type want struct {
		level slog.Level
		err   error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty",
			want: want{
				err: ErrInvalidMinLevel,
			},
		},
		{
			name: "debug",
			args: args{
				l: "debug",
			},
			want: want{
				level: slog.LevelDebug,
			},
		},
		{
			name: "info",
			args: args{
				l: "info",
			},
			want: want{
				level: slog.LevelInfo,
			},
		},
		{
			name: "warn",
			args: args{
				l: "warn",
			},
			want: want{
				level: slog.LevelWarn,
			},
		},
		{
			name: "error",
			args: args{
				l: "error",
			},
			want: want{
				level: slog.LevelError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := MinLevel(context.Background(), tt.args.l)
			if tt.want.err != nil {
				if !errors.Is(err, tt.want.err) {
					t.Errorf("MinLevel() = %v, want %v", err, tt.want.err)
				}
			} else {
				got, ok := ctx.Value(ContextMinLevel).(slog.Level)
				if !ok {
					t.Fatal("not ok")
				}
				if got != tt.want.level {
					t.Errorf("MinLevel() = %v, want %v", got, tt.want.level)
				}
			}
		})
	}
}
