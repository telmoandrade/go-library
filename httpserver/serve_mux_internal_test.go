package httpserver

import (
	"testing"
)

func Test_joinPattern(t *testing.T) {
	type args struct {
		prefix string
		suffix string
	}
	tests := []struct {
		name     string
		args     args
		want     string
		panicked bool
	}{
		{
			name: "test empty",
			args: args{
				prefix: "",
				suffix: "",
			},
			want:     "",
			panicked: false,
		},
		{
			name: "test empty - hosted prefix",
			args: args{
				prefix: "www.test.com",
				suffix: "",
			},
			want:     "www.test.com",
			panicked: false,
		},
		{
			name: "test empty - hosted suffix",
			args: args{
				prefix: "",
				suffix: "www.test.com",
			},
			want:     "www.test.com",
			panicked: false,
		},
		{
			name: "test empty - hosted prefix and suffix",
			args: args{
				prefix: "www.test.com",
				suffix: "www.test.com",
			},
			want:     "www.test.com",
			panicked: false,
		},
		{
			name: "test empty - hosted prefix and suffix invalid",
			args: args{
				prefix: "www.test1.com",
				suffix: "www.test2.com",
			},
			panicked: true,
		},
		{
			name: "test only prefix",
			args: args{
				prefix: "/admin",
				suffix: "",
			},
			want:     "/admin",
			panicked: false,
		},
		{
			name: "test only prefix hosted",
			args: args{
				prefix: "www.test.com/admin",
				suffix: "",
			},
			want:     "www.test.com/admin",
			panicked: false,
		},
		{
			name: "test only suffix",
			args: args{
				prefix: "",
				suffix: "/user",
			},
			want:     "/user",
			panicked: false,
		},
		{
			name: "test only suffix hosted",
			args: args{
				prefix: "",
				suffix: "www.test.com/user",
			},
			want:     "www.test.com/user",
			panicked: false,
		},
		{
			name: "test prefix + suffix",
			args: args{
				prefix: "/admin",
				suffix: "/user",
			},
			want:     "/admin/user",
			panicked: false,
		},
		{
			name: "test prefix + suffix - hosted prefix",
			args: args{
				prefix: "www.test.com/admin",
				suffix: "/user",
			},
			want:     "www.test.com/admin/user",
			panicked: false,
		},
		{
			name: "test prefix + suffix - hosted suffix",
			args: args{
				prefix: "/admin",
				suffix: "www.test.com/user",
			},
			want:     "www.test.com/admin/user",
			panicked: false,
		},
		{
			name: "test prefix + suffix - hosted",
			args: args{
				prefix: "www.test.com/admin",
				suffix: "www.test.com/user",
			},
			want:     "www.test.com/admin/user",
			panicked: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				e := recover()
				if tt.panicked && e == nil {
					t.Fatal("expected panic")
				}
				if !tt.panicked && e != nil {
					t.Fatal("unexpected panic")
				}
			}()
			if got := joinPattern(tt.args.prefix, tt.args.suffix); got != tt.want {
				t.Errorf("joinPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
