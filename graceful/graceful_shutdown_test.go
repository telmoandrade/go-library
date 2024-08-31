package graceful

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	gomock "go.uber.org/mock/gomock"
)

func TestNewGracefulShutdown(t *testing.T) {
	type args struct {
		opts []OptionGracefulShutdown
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
				opts: []OptionGracefulShutdown{
					func(c *gracefulShutdown) {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGracefulShutdown(tt.args.opts...)
			if got == nil {
				t.Errorf("NewControl() = nil, want != nil")
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	tests := []struct {
		name string
		args time.Duration
		want time.Duration
	}{
		{
			name: "timeout 0",
			want: 0,
		},
		{
			name: "timeout 5s",
			args: 5 * time.Second,
			want: 5 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := NewGracefulShutdown(WithTimeout(tt.args))

			if got := gs.timeout; got != tt.want {
				t.Errorf("WithTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithServers(t *testing.T) {
	tests := []struct {
		name string
		args []GracefulServer
		want []GracefulServer
	}{
		{
			name: "empty",
		},
		{
			name: "nil",
			args: nil,
		},
		{
			name: "GracefulServer nil",
			args: []GracefulServer{
				nil,
			},
			want: []GracefulServer{},
		},
		{
			name: "1 GracefulServer",
			args: []GracefulServer{
				&MockGracefulServer{},
			},
			want: []GracefulServer{
				&MockGracefulServer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := NewGracefulShutdown(WithServers(tt.args...))

			if got := gs.gracefulServers; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithServers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gracefulShutdown_runServer(t *testing.T) {
	t.Run("cancel control context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockGracefulServer(ctrl)
		callStart := mock.EXPECT().Start().Return(nil).Times(1)
		mock.EXPECT().Stop(gomock.Any()).Times(1).After(callStart)

		gs := NewGracefulShutdown()
		gs.runServer(mock)
		go func() {
			<-time.After(100 * time.Microsecond)
			gs.cancelCtx()
		}()
		gs.wg.Wait()

		<-time.After(200 * time.Microsecond)
	})

	t.Run("start error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockGracefulServer(ctrl)
		callStart := mock.EXPECT().Start().Return(errors.New("error")).Times(1)
		mock.EXPECT().Stop(gomock.Any()).Times(1).After(callStart)

		gs := NewGracefulShutdown()
		gs.runServer(mock)
		gs.wg.Wait()

		<-time.After(200 * time.Microsecond)
	})

	t.Run("with timeout", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockGracefulServer(ctrl)
		callStart := mock.EXPECT().Start().Return(errors.New("error")).Times(1)

		callStop := mock.EXPECT().Stop(gomock.Any()).Times(1).After(callStart)
		callStop.Do(func(ctx any) { <-time.After(1 * time.Second) })

		mock.EXPECT().ForceStop().Times(1).After(callStop)

		gs := NewGracefulShutdown(WithTimeout(100 * time.Microsecond))
		gs.runServer(mock)
		gs.wg.Wait()

		<-time.After(200 * time.Microsecond)
	})
}

func Test_gracefulShutdown_Run(t *testing.T) {
	t.Run("without servers", func(t *testing.T) {
		gs := NewGracefulShutdown()
		gs.Run(context.Background())
	})

	t.Run("with servers nil", func(t *testing.T) {
		gs := NewGracefulShutdown(WithServers(nil))
		gs.Run(context.Background())
	})

	t.Run("with servers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockGracefulServer(ctrl)
		mock.EXPECT().Start().Return(nil)
		mock.EXPECT().Stop(gomock.Any())

		gs := NewGracefulShutdown(WithServers(mock))

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		gs.Run(ctx)
	})
}
