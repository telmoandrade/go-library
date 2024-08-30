package graceful

import (
	"context"
	"testing"
)

func TestNewGracefulServer(t *testing.T) {
	type args struct {
		opts []OptionGracefulServer
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
				opts: []OptionGracefulServer{
					func(c *gracefulServer) {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGracefulServer(tt.args.opts...)
			if got == nil {
				t.Errorf("NewServer() is null")
			}
		})
	}
}

func TestWithStart(t *testing.T) {
	tests := []struct {
		name string
		args func() error
	}{
		{
			name: "nil",
			args: nil,
		},
		{
			name: "not nil",
			args: func() error { return nil },
		},
	}
	for _, tt := range tests {
		n := NewGracefulServer(WithStart(tt.args))
		gs := n.(*gracefulServer)

		if gs.start == nil {
			t.Errorf("start is null")
		}
	}
}

func TestWithStop(t *testing.T) {
	tests := []struct {
		name string
		args func(context.Context)
	}{
		{
			name: "nil",
			args: nil,
		},
		{
			name: "not nil",
			args: func(context.Context) {},
		},
	}
	for _, tt := range tests {
		n := NewGracefulServer(WithStop(tt.args))
		gs := n.(*gracefulServer)

		if gs.stop == nil {
			t.Errorf("start is null")
		}
	}
}

func TestWithForceStop(t *testing.T) {
	tests := []struct {
		name string
		args func()
	}{
		{
			name: "nil",
			args: nil,
		},
		{
			name: "not nil",
			args: func() {},
		},
	}
	for _, tt := range tests {
		n := NewGracefulServer(WithForceStop(tt.args))
		gs := n.(*gracefulServer)

		if gs.forceStop == nil {
			t.Errorf("start is null")
		}
	}
}

func Test_gracefulServer_Start(t *testing.T) {
	type args struct {
		opts []OptionGracefulServer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "default",
		},
		{
			name: "custom",
			args: args{
				[]OptionGracefulServer{
					WithStart(func() error { return nil }),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := NewGracefulServer(tt.args.opts...)
			if got := gs.Start(); got != nil {
				t.Errorf("Start() = %v, want %v", got, nil)
			}
		})
	}
}

func Test_gracefulServer_Stop(t *testing.T) {
	type args struct {
		opts []OptionGracefulServer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "default",
		},
		{
			name: "custom",
			args: args{
				[]OptionGracefulServer{
					WithStop(func(ctx context.Context) {}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := NewGracefulServer(tt.args.opts...)
			gs.Stop(context.Background())
		})
	}
}

func Test_gracefulServer_ForceStop(t *testing.T) {
	type args struct {
		opts []OptionGracefulServer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "default",
		},
		{
			name: "custom",
			args: args{
				[]OptionGracefulServer{
					WithForceStop(func() {}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := NewGracefulServer(tt.args.opts...)
			gs.ForceStop()
		})
	}
}
