package graceful

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
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
				t.Errorf("NewGracefulShutdown() = nil, want != nil")
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
			n := NewGracefulShutdown(WithTimeout(tt.args))
			gs, _ := n.(*gracefulShutdown)

			if got := gs.timeout; got != tt.want {
				t.Errorf("WithTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithServers(t *testing.T) {
	type want struct {
		simple []GracefulServer
		double []GracefulServer
	}
	tests := []struct {
		name string
		args []GracefulServer
		want want
	}{
		{
			name: "empty",
			want: want{
				simple: []GracefulServer{},
				double: []GracefulServer{},
			},
		},
		{
			name: "nil",
			args: nil,
			want: want{
				simple: []GracefulServer{},
				double: []GracefulServer{},
			},
		},
		{
			name: "GracefulServer nil",
			args: []GracefulServer{
				nil,
			},
			want: want{
				simple: []GracefulServer{},
				double: []GracefulServer{},
			},
		},
		{
			name: "1 GracefulServer",
			args: []GracefulServer{
				&MockGracefulServer{},
			},
			want: want{
				simple: []GracefulServer{&MockGracefulServer{}},
				double: []GracefulServer{&MockGracefulServer{}, &MockGracefulServer{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v simple", tt.name), func(t *testing.T) {
			n := NewGracefulShutdown(WithServers(tt.args...))
			gs, _ := n.(*gracefulShutdown)

			if got := gs.gracefulServers; !reflect.DeepEqual(got, tt.want.simple) {
				t.Errorf("WithServers() = %v, want %v", got, tt.want.simple)
			}
		})

		t.Run(fmt.Sprintf("%v double", tt.name), func(t *testing.T) {
			n := NewGracefulShutdown(
				WithServers(tt.args...),
				WithServers(tt.args...),
			)
			gs, _ := n.(*gracefulShutdown)

			if got := gs.gracefulServers; !reflect.DeepEqual(got, tt.want.double) {
				t.Errorf("WithServers() = %v, want %v", got, tt.want.double)
			}
		})
	}
}

func TestWithNotifyShutdown(t *testing.T) {
	tests := []struct {
		name string
		args func()
	}{
		{
			name: "default",
		},
		{
			name: "custom",
			args: func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewGracefulShutdown(
				WithNotifyShutdown(tt.args),
			)
			gs, _ := n.(*gracefulShutdown)
			if gs.notifyShutdown == nil {
				t.Fatal("notifyShutdown is null")
			}

			gs.notifyShutdown()
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

		n := NewGracefulShutdown()
		gs, _ := n.(*gracefulShutdown)

		gs.runServer(mock)
		go func() {
			<-time.After(100 * time.Microsecond)
			gs.cancelCtx()
		}()
		gs.wg.Wait()
	})

	t.Run("start error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockGracefulServer(ctrl)
		callStart := mock.EXPECT().Start().Return(errors.New("error")).Times(1)
		mock.EXPECT().Stop(gomock.Any()).Times(1).After(callStart)

		n := NewGracefulShutdown()
		gs, _ := n.(*gracefulShutdown)

		gs.runServer(mock)
		gs.wg.Wait()
	})

	t.Run("with timeout", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockGracefulServer(ctrl)
		callStart := mock.EXPECT().Start().Return(errors.New("error")).Times(1)

		callStop := mock.EXPECT().Stop(gomock.Any()).Times(1).After(callStart)
		callStop.Do(func(ctx any) { <-time.After(1 * time.Second) })

		mock.EXPECT().ForceStop().Times(1).After(callStop)

		n := NewGracefulShutdown(WithTimeout(100 * time.Microsecond))
		gs, _ := n.(*gracefulShutdown)

		gs.runServer(mock)
		gs.wg.Wait()
	})
}

func Test_gracefulShutdown_Run(t *testing.T) {
	t.Run("without servers", func(t *testing.T) {
		callNotifyShutdown := 0

		gs := NewGracefulShutdown(
			WithNotifyShutdown(func() { callNotifyShutdown = callNotifyShutdown + 1 }),
		)
		gs.Run(context.Background())

		if callNotifyShutdown != 0 {
			t.Errorf("NotifyShutdown call %v want 0", callNotifyShutdown)
		}
	})

	t.Run("with servers nil", func(t *testing.T) {
		callNotifyShutdown := 0

		gs := NewGracefulShutdown(
			WithServers(nil),
			WithNotifyShutdown(func() { callNotifyShutdown = callNotifyShutdown + 1 }),
		)
		gs.Run(context.Background())

		if callNotifyShutdown != 0 {
			t.Errorf("NotifyShutdown call %v want 0", callNotifyShutdown)
		}
	})

	t.Run("with servers", func(t *testing.T) {
		callNotifyShutdown := 0
		wgNotifyShutdown := sync.WaitGroup{}
		wgNotifyShutdown.Add(1)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockGracefulServer(ctrl)
		mock.EXPECT().Start().Return(errors.New("error"))
		mock.EXPECT().Stop(gomock.Any())

		gs := NewGracefulShutdown(
			WithServers(mock),
			WithNotifyShutdown(func() {
				callNotifyShutdown = callNotifyShutdown + 1
				wgNotifyShutdown.Done()
			}),
		)
		gs.Run(context.Background())

		wgNotifyShutdown.Wait()

		if callNotifyShutdown != 1 {
			t.Errorf("NotifyShutdown call %v want 1", callNotifyShutdown)
		}
	})
}
