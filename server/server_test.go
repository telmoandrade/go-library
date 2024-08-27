package server

import (
	"context"
	"testing"
)

func TestNewServer(t *testing.T) {
	type args struct {
		opts []OptionServer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "default",
		},
		{
			name: "WithStart",
			args: args{
				[]OptionServer{
					WithStart(func() error { return nil }),
				},
			},
		},
		{
			name: "WithStop",
			args: args{
				[]OptionServer{
					WithStop(func(ctx context.Context) {}),
				},
			},
		},
		{
			name: "WithForceStop",
			args: args{
				[]OptionServer{
					WithForceStop(func() {}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServer(tt.args.opts...)
			if got == nil {
				t.Errorf("NewServer() is null")
			}
		})
	}
}

func Test_server_Start(t *testing.T) {
	type args struct {
		opts []OptionServer
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
				[]OptionServer{
					WithStart(func() error { return nil }),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer(tt.args.opts...)
			if got := s.Start(); got != nil {
				t.Errorf("Start() = %v, want %v", got, nil)
			}
		})
	}
}

func Test_server_Stop(t *testing.T) {
	type args struct {
		opts []OptionServer
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
				[]OptionServer{
					WithStop(func(ctx context.Context) {}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer(tt.args.opts...)
			s.Stop(context.Background())
		})
	}
}

func Test_server_ForceStop(t *testing.T) {
	type args struct {
		opts []OptionServer
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
				[]OptionServer{
					WithForceStop(func() {}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer(tt.args.opts...)
			s.ForceStop()
		})
	}
}
