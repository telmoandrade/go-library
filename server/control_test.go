package server

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	gomock "go.uber.org/mock/gomock"
)

func TestNewControl(t *testing.T) {
	type args struct {
		opts []OptionControl
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
				opts: []OptionControl{
					func(c *control) {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewControl(tt.args.opts...)
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
			name: "default",
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
			c := NewControl(WithTimeout(tt.args))

			if got := c.timeout; got != tt.want {
				t.Errorf("WithTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithServers(t *testing.T) {
	tests := []struct {
		name string
		args []Server
		want []Server
	}{
		{
			name: "default",
		},
		{
			name: "nil",
			args: nil,
		},
		{
			name: "server nil",
			args: []Server{
				nil,
			},
			want: []Server{},
		},
		{
			name: "1 server",
			args: []Server{
				&MockServer{},
			},
			want: []Server{
				&MockServer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewControl(WithServers(tt.args...))

			if got := c.servers; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithServers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_control_runServer(t *testing.T) {
	t.Run("cancel control context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockServer(ctrl)
		callStart := mock.EXPECT().Start().Return(nil).Times(1)
		mock.EXPECT().Stop(gomock.Any()).Times(1).After(callStart)

		c := NewControl()
		c.runServer(mock)
		go func() {
			<-time.After(100 * time.Microsecond)
			c.cancelCtx()
		}()
		c.wg.Wait()

		<-time.After(200 * time.Microsecond)
	})

	t.Run("start error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockServer(ctrl)
		callStart := mock.EXPECT().Start().Return(errors.New("error")).Times(1)
		mock.EXPECT().Stop(gomock.Any()).Times(1).After(callStart)

		c := NewControl()
		c.runServer(mock)
		c.wg.Wait()

		<-time.After(200 * time.Microsecond)
	})

	t.Run("with timeout", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockServer(ctrl)
		callStart := mock.EXPECT().Start().Return(errors.New("error")).Times(1)

		callStop := mock.EXPECT().Stop(gomock.Any()).Times(1).After(callStart)
		callStop.Do(func(ctx any) { <-time.After(1 * time.Second) })

		mock.EXPECT().ForceStop().Times(1).After(callStop)

		c := NewControl(WithTimeout(100 * time.Microsecond))
		c.runServer(mock)
		c.wg.Wait()

		<-time.After(200 * time.Microsecond)
	})
}

func Test_control_Run(t *testing.T) {
	t.Run("without servers", func(t *testing.T) {
		c := NewControl()
		c.Run(context.Background())
	})

	t.Run("with servers nil", func(t *testing.T) {
		c := NewControl(WithServers(nil))
		c.Run(context.Background())
	})

	t.Run("with servers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockServer(ctrl)
		mock.EXPECT().Start().Return(nil)
		mock.EXPECT().Stop(gomock.Any())

		c := NewControl(WithServers(mock))

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		c.Run(ctx)
	})
}
