package httpserver

import (
	"reflect"
	"testing"
)

func Test_newPatternRoute(t *testing.T) {
	type args struct {
		pattern string
	}
	type want struct {
		host      string
		pattern   string
		endSlash  bool
		multiName string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "pattern empty",
			args: args{
				pattern: "",
			},
			want: want{
				host:     "",
				pattern:  "/",
				endSlash: true,
			},
		},
		{
			name: "pattern host",
			args: args{
				pattern: "host",
			},
			want: want{
				host:     "host",
				pattern:  "/",
				endSlash: true,
			},
		},
		{
			name: "pattern host\\user",
			args: args{
				pattern: "host/user",
			},
			want: want{
				host:     "host",
				pattern:  "/user",
				endSlash: false,
			},
		},
		{
			name: "pattern host\\user\\",
			args: args{
				pattern: "host/user/",
			},
			want: want{
				host:     "host",
				pattern:  "/user/",
				endSlash: true,
			},
		},
		{
			name: "pattern host\\user\\{id}",
			args: args{
				pattern: "host/user/{id}",
			},
			want: want{
				host:     "host",
				pattern:  "/user/{id}",
				endSlash: false,
			},
		},
		{
			name: "pattern host\\{$}",
			args: args{
				pattern: "host/{$}",
			},
			want: want{
				host:     "host",
				pattern:  "/{$}",
				endSlash: false,
			},
		},
		{
			name: "pattern host\\{path...}",
			args: args{
				pattern: "host/{path...}",
			},
			want: want{
				host:      "host",
				pattern:   "/{path...}",
				endSlash:  false,
				multiName: "path",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := newPatternRoute(tt.args.pattern)

			if pr.host != tt.want.host {
				t.Errorf("newPatternRoute().host = %v, want %v", pr.host, tt.want.host)
			}
			if pr.pattern != tt.want.pattern {
				t.Errorf("newPatternRoute().pattern = %v, want %v", pr.pattern, tt.want.pattern)
			}
			if pr.endSlash != tt.want.endSlash {
				t.Errorf("newPatternRoute().endSlash = %v, want %v", pr.endSlash, tt.want.endSlash)
			}
			if pr.multiName != tt.want.multiName {
				t.Errorf("newPatternRoute().opened = %v, want %v", pr.multiName, tt.want.multiName)
			}
		})
	}
}

func Test_patternRoute_extractWildcardSegment(t *testing.T) {
	type fields struct {
		wildcard map[string]bool
	}
	type args struct {
		seg string
	}
	type want struct {
		panicked   bool
		panicError string
		wildcard   map[string]bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
		args   args
	}{
		{
			name:   "{$}",
			fields: fields{wildcard: map[string]bool{}},
			args:   args{seg: "{$}"},
			want:   want{wildcard: map[string]bool{}},
		},
		{
			name:   "{id}",
			fields: fields{wildcard: map[string]bool{}},
			args:   args{seg: "{id}"},
			want:   want{wildcard: map[string]bool{"id": true}},
		},
		{
			name:   "{id...}",
			fields: fields{wildcard: map[string]bool{}},
			args:   args{seg: "{id...}"},
			want:   want{wildcard: map[string]bool{"id": true}},
		},
		{
			name:   "duplicate {id}",
			fields: fields{wildcard: map[string]bool{"id": true}},
			args:   args{seg: "{id}"},
			want:   want{panicked: true, panicError: "httpserver: duplicate wildcard name \"id\""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.want.panicked {
					if r == nil {
						t.Fatalf("expected panic = %s", tt.want.panicError)
					}

					err := r.(error).Error()
					if err != tt.want.panicError {
						t.Errorf("error = %v, want %v", err, tt.want.panicError)
					}
				} else if r != nil {
					t.Fatalf("unexpected panic = %s", r)
				}
			}()

			pr := &patternRoute{
				wildcard: tt.fields.wildcard,
			}
			pr.extractWildcardSegment(tt.args.seg)

			if !reflect.DeepEqual(pr.wildcard, tt.want.wildcard) {
				t.Errorf("wildcard = %v, want %v", pr.wildcard, tt.want.wildcard)
			}
		})
	}
}

func Test_patternRoute_extractWildcard(t *testing.T) {
	type fields struct {
		pattern string
	}
	type want struct {
		endSlash   bool
		wildcard   map[string]bool
		panicked   bool
		panicError string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "pattern \\",
			fields: fields{
				pattern: "/",
			},
			want: want{
				endSlash:   false,
				panicked:   false,
				panicError: "",
				wildcard:   map[string]bool{},
			},
		},
		{
			name: "pattern \\{$}",
			fields: fields{
				pattern: "/{$}",
			},
			want: want{
				endSlash:   false,
				panicked:   false,
				panicError: "",
				wildcard:   map[string]bool{},
			},
		},
		{
			name: "pattern \\{id...}",
			fields: fields{
				pattern: "/{id}",
			},
			want: want{
				endSlash:   false,
				panicked:   false,
				panicError: "",
				wildcard: map[string]bool{
					"id": true,
				},
			},
		},
		{
			name: "pattern \\{id}",
			fields: fields{
				pattern: "/{id}",
			},
			want: want{
				endSlash:   false,
				panicked:   false,
				panicError: "",
				wildcard: map[string]bool{
					"id": true,
				},
			},
		},
		{
			name: "pattern \\user\\{id}\\details",
			fields: fields{
				pattern: "/user/{id}/details",
			},
			want: want{
				endSlash:   false,
				panicked:   false,
				panicError: "",
				wildcard: map[string]bool{
					"id": true,
				},
			},
		},
		{
			name: "pattern \\{id1}\\{id2}",
			fields: fields{
				pattern: "/{id1}/{id2}",
			},
			want: want{
				endSlash:   false,
				panicked:   false,
				panicError: "",
				wildcard: map[string]bool{
					"id1": true,
					"id2": true,
				},
			},
		},
		{
			name: "pattern \\{id}\\{id}",
			fields: fields{
				pattern: "/{id}/{id}",
			},
			want: want{
				endSlash:   false,
				panicked:   true,
				panicError: "httpserver: duplicate wildcard name \"id\"",
				wildcard:   map[string]bool{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &patternRoute{
				wildcard: map[string]bool{},
				pattern:  tt.fields.pattern,
			}

			defer func() {
				r := recover()
				if tt.want.panicked {
					if r == nil {
						t.Fatalf("expected panic = %s", tt.want.panicError)
					}

					err := r.(error).Error()
					if err != tt.want.panicError {
						t.Errorf("error = %v, want %v", err, tt.want.panicError)
					}
				} else if r != nil {
					t.Fatalf("unexpected panic = %s", r)
				}
			}()

			pr.extractWildcard()

			if !reflect.DeepEqual(pr.wildcard, tt.want.wildcard) {
				t.Errorf("wildcard = %v, want %v", pr.wildcard, tt.want.wildcard)
			}
		})
	}
}

func Test_patternRouteJoinWildcard(t *testing.T) {
	type args struct {
		prefixPr *patternRoute
		suffixPr *patternRoute
	}
	type want struct {
		wildcard   map[string]bool
		panicked   bool
		panicError string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty join empty",
			args: args{
				prefixPr: &patternRoute{
					wildcard: map[string]bool{},
				},
				suffixPr: &patternRoute{
					wildcard: map[string]bool{},
				},
			},
			want: want{
				wildcard:   map[string]bool{},
				panicked:   false,
				panicError: "",
			},
		},
		{
			name: "1 join empty",
			args: args{
				prefixPr: &patternRoute{
					wildcard: map[string]bool{
						"id": true,
					},
				},
				suffixPr: &patternRoute{
					wildcard: map[string]bool{},
				},
			},
			want: want{
				wildcard: map[string]bool{
					"id": true,
				},
				panicked:   false,
				panicError: "",
			},
		},
		{
			name: "empty join 1",
			args: args{
				prefixPr: &patternRoute{
					wildcard: map[string]bool{},
				},
				suffixPr: &patternRoute{
					wildcard: map[string]bool{
						"id": true,
					},
				},
			},
			want: want{
				wildcard: map[string]bool{
					"id": true,
				},
				panicked:   false,
				panicError: "",
			},
		},
		{
			name: "1 join 1",
			args: args{
				prefixPr: &patternRoute{
					wildcard: map[string]bool{
						"id1": true,
					},
				},
				suffixPr: &patternRoute{
					wildcard: map[string]bool{
						"id2": true,
					},
				},
			},
			want: want{
				wildcard: map[string]bool{
					"id1": true,
					"id2": true,
				},
				panicked:   false,
				panicError: "",
			},
		},
		{
			name: "conflicts",
			args: args{
				prefixPr: &patternRoute{
					wildcard: map[string]bool{
						"id": true,
					},
				},
				suffixPr: &patternRoute{
					wildcard: map[string]bool{
						"id": true,
					},
				},
			},
			want: want{
				panicked:   true,
				panicError: "httpserver: duplicate wildcard name \"id\"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.want.panicked {
					if r == nil {
						t.Fatalf("expected panic = %s", tt.want.panicError)
					}

					err := r.(error).Error()
					if err != tt.want.panicError {
						t.Errorf("error = %v, want %v", err, tt.want.panicError)
					}
				} else if r != nil {
					t.Fatalf("unexpected panic = %s", r)
				}
			}()

			if got := patternRouteJoinWildcard(tt.args.prefixPr, tt.args.suffixPr); !reflect.DeepEqual(got, tt.want.wildcard) {
				t.Errorf("patternRouteJoinWildcard() = %v, want %v", got, tt.want.wildcard)
			}
		})
	}
}

func Test_patternRouteJoinHost(t *testing.T) {
	type args struct {
		prefixHost string
		suffixHost string
	}
	type want struct {
		panicked   bool
		panicError string
		host       string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "host join empty",
			args: args{
				prefixHost: "host",
				suffixHost: "",
			},
			want: want{
				host: "host",
			},
		},
		{
			name: "empty join host",
			args: args{
				prefixHost: "",
				suffixHost: "host",
			},
			want: want{
				host: "host",
			},
		},
		{
			name: "host join host",
			args: args{
				prefixHost: "host",
				suffixHost: "host",
			},
			want: want{
				host: "host",
			},
		},
		{
			name: "host1 join host2",
			args: args{
				prefixHost: "host1",
				suffixHost: "host2",
			},
			want: want{
				panicked:   true,
				panicError: "httpserver: host host1 conflicts with host host2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.want.panicked {
					if r == nil {
						t.Fatalf("expected panic = %s", tt.want.panicError)
					}

					err := r.(error).Error()
					if err != tt.want.panicError {
						t.Errorf("error = %v, want %v", err, tt.want.panicError)
					}
				} else if r != nil {
					t.Fatalf("unexpected panic = %s", r)
				}
			}()

			if got := patternRouteJoinHost(tt.args.prefixHost, tt.args.suffixHost); got != tt.want.host {
				t.Errorf("patternRouteJoinHost() = %v, want %v", got, tt.want.host)
			}
		})
	}
}

func Test_patternRoute_removeEndSlash(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{pattern: ""},
			want: "",
		},
		{
			name: "\\",
			args: args{pattern: "/"},
			want: "",
		},
		{
			name: "\\admin",
			args: args{pattern: "/admin"},
			want: "/admin",
		},
		{
			name: "\\admin\\",
			args: args{pattern: "/admin/"},
			want: "/admin",
		},
		{
			name: "host",
			args: args{pattern: "host"},
			want: "",
		},
		{
			name: "host\\",
			args: args{pattern: "host/"},
			want: "",
		},
		{
			name: "host\\admin",
			args: args{pattern: "host/admin"},
			want: "/admin",
		},
		{
			name: "host\\admin\\",
			args: args{pattern: "host/admin/"},
			want: "/admin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := newPatternRoute(tt.args.pattern)

			if got := pr.removeEndSlash(); got != tt.want {
				t.Errorf("patternRoute.removeEndSlash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_patternRoute_join(t *testing.T) {
	type args struct {
		pattern1 string
		pattern2 string
	}
	type want struct {
		host       string
		pattern    string
		endSlash   bool
		panicked   bool
		panicError string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty join empty",
			args: args{
				pattern1: "",
				pattern2: "",
			},
			want: want{
				host:     "",
				pattern:  "/",
				endSlash: true,
			},
		},
		{
			name: "host join empty",
			args: args{
				pattern1: "host",
				pattern2: "",
			},
			want: want{
				host:     "host",
				pattern:  "/",
				endSlash: true,
			},
		},
		{
			name: "empty join host",
			args: args{
				pattern1: "",
				pattern2: "host",
			},
			want: want{
				host:     "host",
				pattern:  "/",
				endSlash: true,
			},
		},
		{
			name: "host join host",
			args: args{
				pattern1: "host",
				pattern2: "host",
			},
			want: want{
				host:     "host",
				pattern:  "/",
				endSlash: true,
			},
		},
		{
			name: "host1 join host2",
			args: args{
				pattern1: "host1",
				pattern2: "host2",
			},
			want: want{
				panicked:   true,
				panicError: "httpserver: host host1 conflicts with host host2",
			},
		},
		{
			name: "\\ join \\",
			args: args{
				pattern1: "/",
				pattern2: "/",
			},
			want: want{
				host:     "",
				pattern:  "/",
				endSlash: true,
			},
		},
		{
			name: "\\admin join empty",
			args: args{
				pattern1: "/admin",
				pattern2: "",
			},
			want: want{
				host:     "",
				pattern:  "/admin",
				endSlash: false,
			},
		},
		{
			name: "empty join \\admin",
			args: args{
				pattern1: "",
				pattern2: "/admin",
			},
			want: want{
				host:     "",
				pattern:  "/admin",
				endSlash: false,
			},
		},
		{
			name: "empty join \\admin\\",
			args: args{
				pattern1: "",
				pattern2: "/admin/",
			},
			want: want{
				host:     "",
				pattern:  "/admin/",
				endSlash: true,
			},
		},
		{
			name: "empty join \\admin\\{$}",
			args: args{
				pattern1: "",
				pattern2: "/admin/{$}",
			},
			want: want{
				host:     "",
				pattern:  "/admin/{$}",
				endSlash: false,
			},
		},
		{
			name: "\\admin join \\",
			args: args{
				pattern1: "/admin",
				pattern2: "/",
			},
			want: want{
				host:     "",
				pattern:  "/admin",
				endSlash: false,
			},
		},
		{
			name: "\\ join \\admin",
			args: args{
				pattern1: "/",
				pattern2: "/admin",
			},
			want: want{
				host:     "",
				pattern:  "/admin",
				endSlash: false,
			},
		},
		{
			name: "\\user join \\{id}",
			args: args{
				pattern1: "/user",
				pattern2: "/{id}",
			},
			want: want{
				host:     "",
				pattern:  "/user/{id}",
				endSlash: false,
			},
		},
		{
			name: "\\user join \\{id}\\",
			args: args{
				pattern1: "/user",
				pattern2: "/{id}/",
			},
			want: want{
				host:     "",
				pattern:  "/user/{id}/",
				endSlash: true,
			},
		},
		{
			name: "\\user join \\{id}\\{$}",
			args: args{
				pattern1: "/user",
				pattern2: "/{id}/{$}",
			},
			want: want{
				host:     "",
				pattern:  "/user/{id}/{$}",
				endSlash: false,
			},
		},
		{
			name: "host\\user join host\\{id}\\{$}",
			args: args{
				pattern1: "host/user",
				pattern2: "host/{id}/{$}",
			},
			want: want{
				host:     "host",
				pattern:  "/user/{id}/{$}",
				endSlash: false,
			},
		},
		{
			name: "\\\\ join \\\\",
			args: args{
				pattern1: "//",
				pattern2: "//",
			},
			want: want{
				host:     "",
				pattern:  "///",
				endSlash: true,
			},
		},
		{
			name: "\\\\ join \\{id}",
			args: args{
				pattern1: "//",
				pattern2: "/{id}",
			},
			want: want{
				host:     "",
				pattern:  "//{id}",
				endSlash: false,
			},
		},
		{
			name: "\\\\\\ join \\\\\\",
			args: args{
				pattern1: "///",
				pattern2: "///",
			},
			want: want{
				host:     "",
				pattern:  "/////",
				endSlash: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.want.panicked {
					if r == nil {
						t.Fatalf("expected panic = %s", tt.want.panicError)
					}

					err := r.(error).Error()
					if err != tt.want.panicError {
						t.Errorf("error = %v, want %v", err, tt.want.panicError)
					}
				} else if r != nil {
					t.Fatalf("unexpected panic = %s", r)
				}
			}()

			pr := newPatternRoute(tt.args.pattern1).join(tt.args.pattern2)

			if pr.host != tt.want.host {
				t.Errorf("join().host = %v, want %v", pr.host, tt.want.host)
			}
			if pr.pattern != tt.want.pattern {
				t.Errorf("join().pattern = %v, want %v", pr.pattern, tt.want.pattern)
			}
			if pr.endSlash != tt.want.endSlash {
				t.Errorf("join().endSlash = %v, want %v", pr.endSlash, tt.want.endSlash)
			}
		})
	}
}

func Test_patternRoute_String(t *testing.T) {
	type args struct {
		pattern string
	}
	type want struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty",
			args: args{
				pattern: "",
			},
			want: want{
				pattern: "/",
			},
		},
		{
			name: "host",
			args: args{
				pattern: "host",
			},
			want: want{
				pattern: "host/",
			},
		},
		{
			name: "\\",
			args: args{
				pattern: "",
			},
			want: want{
				pattern: "/",
			},
		},
		{
			name: "\\admin",
			args: args{
				pattern: "/admin",
			},
			want: want{
				pattern: "/admin",
			},
		},
		{
			name: "host\\admin",
			args: args{
				pattern: "host/admin",
			},
			want: want{
				pattern: "host/admin",
			},
		},
		{
			name: "host\\admin\\",
			args: args{
				pattern: "host/admin/",
			},
			want: want{
				pattern: "host/admin/",
			},
		},
		{
			name: "host\\admin\\{$}",
			args: args{
				pattern: "host/admin/{$}",
			},
			want: want{
				pattern: "host/admin/{$}",
			},
		},
		{
			name: "host\\admin\\{path...}",
			args: args{
				pattern: "host/admin/{path...}",
			},
			want: want{
				pattern: "host/admin/{path...}",
			},
		},
		{
			name: "host\\admin\\{id}",
			args: args{
				pattern: "host/admin/{id}",
			},
			want: want{
				pattern: "host/admin/{id}",
			},
		},
		{
			name: "host\\admin\\{id}\\",
			args: args{
				pattern: "host/admin/{id}/",
			},
			want: want{
				pattern: "host/admin/{id}/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := newPatternRoute(tt.args.pattern)

			if got := pr.String(); got != tt.want.pattern {
				t.Errorf("patternRoute.String() = %v, want %v", got, tt.want.pattern)
			}
		})
	}
}

func Test_patternRoute_mountMethodNotAllowed(t *testing.T) {
	type args struct {
		pattern string
	}
	type want struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "pattern empty",
			args: args{
				pattern: "",
			},
			want: want{
				pattern: "/",
			},
		},
		{
			name: "pattern host",
			args: args{
				pattern: "host",
			},
			want: want{
				pattern: "host/",
			},
		},
		{
			name: "pattern host\\user",
			args: args{
				pattern: "host/user",
			},
			want: want{
				pattern: "host/user",
			},
		},
		{
			name: "pattern host\\user\\",
			args: args{
				pattern: "host/user/",
			},
			want: want{
				pattern: "host/user/",
			},
		},
		{
			name: "pattern host\\user\\{id}",
			args: args{
				pattern: "host/user/{id}",
			},
			want: want{
				pattern: "host/user/{id}",
			},
		},
		{
			name: "pattern host\\{$}",
			args: args{
				pattern: "host/{$}",
			},
			want: want{
				pattern: "host/{$}",
			},
		},
		{
			name: "pattern host\\{path...}",
			args: args{
				pattern: "host/{path...}",
			},
			want: want{
				pattern: "host/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := newPatternRoute(tt.args.pattern)

			if got := pr.mountMethodNotAllowed(); got != tt.want.pattern {
				t.Errorf("patternRoute.mountMethodNotAllowed() = %v, want %v", got, tt.want.pattern)
			}
		})
	}
}
