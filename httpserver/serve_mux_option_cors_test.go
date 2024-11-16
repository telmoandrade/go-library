package httpserver

import (
	"net/http"
	"reflect"
	"testing"
)

func TestWithCors(t *testing.T) {
	type args struct {
		opts []OptionCors
	}
	type want struct {
		allowedOrigins         []string
		allowedWildcardOrigins []wildcard
		allowedHeaders         []string
		exposedHeaders         []string
		maxAge                 int
		allowedOriginsAll      bool
		allowedHeadersAll      bool
		allowCredentials       bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty",
			want: want{
				allowedOrigins:         []string{},
				allowedWildcardOrigins: []wildcard{},
				allowedHeaders:         []string{},
				exposedHeaders:         []string{},
				maxAge:                 86400,
				allowedOriginsAll:      true,
				allowedHeadersAll:      true,
				allowCredentials:       false,
			},
		},
		{
			name: "with OptionCors",
			args: args{
				opts: []OptionCors{WithAllowCredentials(true)},
			},
			want: want{
				allowedOrigins:         []string{},
				allowedWildcardOrigins: []wildcard{},
				allowedHeaders:         []string{},
				exposedHeaders:         []string{},
				maxAge:                 86400,
				allowedOriginsAll:      true,
				allowedHeadersAll:      true,
				allowCredentials:       true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux(
				WithCors(tt.args.opts...),
			)
			mux := muxInterface.(*serveMux)
			c := mux.config.cors

			if !reflect.DeepEqual(c.allowedOrigins, tt.want.allowedOrigins) {
				t.Errorf("cors.allowedOrigins = %v, want %v", c.allowedOrigins, tt.want.allowedOrigins)
			}
			if !reflect.DeepEqual(c.allowedWildcardOrigins, tt.want.allowedWildcardOrigins) {
				t.Errorf("cors.allowedWildcardOrigins = %v, want %v", c.allowedWildcardOrigins, tt.want.allowedWildcardOrigins)
			}
			if !reflect.DeepEqual(c.allowedHeaders, tt.want.allowedHeaders) {
				t.Errorf("cors.allowedHeaders = %v, want %v", c.allowedHeaders, tt.want.allowedHeaders)
			}
			if !reflect.DeepEqual(c.exposedHeaders, tt.want.exposedHeaders) {
				t.Errorf("cors.exposedHeaders = %v, want %v", c.exposedHeaders, tt.want.exposedHeaders)
			}
			if c.maxAge != tt.want.maxAge {
				t.Errorf("cors.maxAge = %v, want %v", c.maxAge, tt.want.maxAge)
			}
			if c.allowedOriginsAll != tt.want.allowedOriginsAll {
				t.Errorf("cors.allowedOriginsAll = %v, want %v", c.allowedOriginsAll, tt.want.allowedOriginsAll)
			}
			if c.allowedHeadersAll != tt.want.allowedHeadersAll {
				t.Errorf("cors.allowedHeadersAll = %v, want %v", c.allowedHeadersAll, tt.want.allowedHeadersAll)
			}
			if c.allowCredentials != tt.want.allowCredentials {
				t.Errorf("cors.allowCredentials = %v, want %v", c.allowCredentials, tt.want.allowCredentials)
			}
		})
	}
}

func TestWithAllowedOrigins(t *testing.T) {
	type fields struct {
		allowedOrigins []string
	}
	type args struct {
		allowedOrigins []string
	}
	type want struct {
		allowedOrigins []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "empty",
			fields: fields{
				allowedOrigins: []string{},
			},
			want: want{
				allowedOrigins: []string{},
			},
		},
		{
			name: "not found",
			fields: fields{
				allowedOrigins: []string{},
			},
			args: args{
				allowedOrigins: []string{"http://www.foo.com", "http://www.bar.com"},
			},
			want: want{
				allowedOrigins: []string{"http://www.foo.com", "http://www.bar.com"},
			},
		},
		{
			name: "found",
			fields: fields{
				allowedOrigins: []string{"http://www.foo.com", "http://www.bar.com"},
			},
			args: args{
				allowedOrigins: []string{"http://www.foo.com", "http://www.bar.com"},
			},
			want: want{
				allowedOrigins: []string{"http://www.foo.com", "http://www.bar.com"},
			},
		},
		{
			name: "empty value",
			fields: fields{
				allowedOrigins: []string{"http://www.foo.com", "http://www.bar.com"},
			},
			args: args{
				allowedOrigins: []string{"  ", "http://www.test.com"},
			},
			want: want{
				allowedOrigins: []string{"http://www.foo.com", "http://www.bar.com", "http://www.test.com"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cors{
				allowedOrigins: tt.fields.allowedOrigins,
			}
			WithAllowedOrigins(tt.args.allowedOrigins...)(c)

			if !reflect.DeepEqual(c.allowedOrigins, tt.want.allowedOrigins) {
				t.Errorf("cors.allowedOrigins = %v, want %v", c.allowedOrigins, tt.want.allowedOrigins)
			}
		})
	}
}

func TestWithAllowOriginFunc(t *testing.T) {
	type args struct {
		fn func(r *http.Request, origin string) bool
	}
	type want struct {
		hasAllowOriginFunc bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty",
		},
		{
			name: "empty",
			args: args{
				fn: func(r *http.Request, origin string) bool { return true },
			},
			want: want{
				hasAllowOriginFunc: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cors{}
			WithAllowOriginFunc(tt.args.fn)(c)

			if tt.want.hasAllowOriginFunc && c.allowOriginFunc == nil {
				t.Errorf("cors.allowOriginFunc = %p", c.allowOriginFunc)
			}
			if !tt.want.hasAllowOriginFunc && c.allowOriginFunc != nil {
				t.Errorf("cors.allowOriginFunc = %p", c.allowOriginFunc)
			}
		})
	}
}

func TestWithAllowedHeaders(t *testing.T) {
	type fields struct {
		allowedHeaders []string
	}
	type args struct {
		allowedHeaders []string
	}
	type want struct {
		allowedHeaders []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "empty",
			fields: fields{
				allowedHeaders: []string{},
			},
			want: want{
				allowedHeaders: []string{},
			},
		},
		{
			name: "not found",
			fields: fields{
				allowedHeaders: []string{},
			},
			args: args{
				allowedHeaders: []string{"Header1", "Header2"},
			},
			want: want{
				allowedHeaders: []string{"Header1", "Header2"},
			},
		},
		{
			name: "found",
			fields: fields{
				allowedHeaders: []string{"Header1", "Header2"},
			},
			args: args{
				allowedHeaders: []string{"Header1", "Header2"},
			},
			want: want{
				allowedHeaders: []string{"Header1", "Header2"},
			},
		},
		{
			name: "empty value",
			fields: fields{
				allowedHeaders: []string{"Header1", "Header2"},
			},
			args: args{
				allowedHeaders: []string{"  ", "Header3"},
			},
			want: want{
				allowedHeaders: []string{"Header1", "Header2", "Header3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cors{
				allowedHeaders: tt.fields.allowedHeaders,
			}
			WithAllowedHeaders(tt.args.allowedHeaders...)(c)

			if !reflect.DeepEqual(c.allowedHeaders, tt.want.allowedHeaders) {
				t.Errorf("cors.allowedHeaders = %v, want %v", c.allowedHeaders, tt.want.allowedHeaders)
			}
		})
	}
}

func TestWithExposedHeaders(t *testing.T) {
	type fields struct {
		exposedHeaders []string
	}
	type args struct {
		exposedHeaders []string
	}
	type want struct {
		exposedHeaders []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "empty",
			fields: fields{
				exposedHeaders: []string{},
			},
			want: want{
				exposedHeaders: []string{},
			},
		},
		{
			name: "not found",
			fields: fields{
				exposedHeaders: []string{},
			},
			args: args{
				exposedHeaders: []string{"Header1", "Header2"},
			},
			want: want{
				exposedHeaders: []string{"Header1", "Header2"},
			},
		},
		{
			name: "found",
			fields: fields{
				exposedHeaders: []string{"Header1", "Header2"},
			},
			args: args{
				exposedHeaders: []string{"Header1", "Header2"},
			},
			want: want{
				exposedHeaders: []string{"Header1", "Header2"},
			},
		},
		{
			name: "empty value",
			fields: fields{
				exposedHeaders: []string{"Header1", "Header2"},
			},
			args: args{
				exposedHeaders: []string{"  ", "Header3"},
			},
			want: want{
				exposedHeaders: []string{"Header1", "Header2", "Header3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cors{
				exposedHeaders: tt.fields.exposedHeaders,
			}
			WithExposedHeaders(tt.args.exposedHeaders...)(c)

			if !reflect.DeepEqual(c.exposedHeaders, tt.want.exposedHeaders) {
				t.Errorf("cors.exposedHeaders = %v, want %v", c.exposedHeaders, tt.want.exposedHeaders)
			}
		})
	}
}

func TestWithAllowCredentials(t *testing.T) {
	type args struct {
		allowCredentials bool
	}
	type want struct {
		allowCredentials bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty",
			want: want{
				allowCredentials: false,
			},
		},
		{
			name: "empty",
			args: args{
				allowCredentials: true,
			},
			want: want{
				allowCredentials: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cors{}
			WithAllowCredentials(tt.args.allowCredentials)(c)

			if c.allowCredentials != tt.want.allowCredentials {
				t.Errorf("cors.allowCredentials = %v, want %v", c.allowCredentials, tt.want.allowCredentials)
			}
		})
	}
}

func TestWithCorsMaxAge(t *testing.T) {
	type args struct {
		seconds int
	}
	type want struct {
		maxAge int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty",
			want: want{
				maxAge: 0,
			},
		},
		{
			name: "maxAge",
			args: args{
				seconds: 1,
			},
			want: want{
				maxAge: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cors{}
			WithCorsMaxAge(tt.args.seconds)(c)

			if c.maxAge != tt.want.maxAge {
				t.Errorf("cors.maxAge = %v, want %v", c.maxAge, tt.want.maxAge)
			}
		})
	}
}

func Test_cors_computeAllowedOrigins(t *testing.T) {
	type fields struct {
		allowedOrigins         []string
		allowedWildcardOrigins []wildcard
		allowOriginFunc        func(r *http.Request, origin string) bool
		allowedOriginsAll      bool
	}
	type want struct {
		allowedOrigins         []string
		allowedWildcardOrigins []wildcard
		allowedOriginsAll      bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty",
			want: want{
				allowedOriginsAll: true,
			},
		},
		{
			name: "allowOriginFunc",
			fields: fields{
				allowOriginFunc: func(r *http.Request, origin string) bool { return true },
			},
			want: want{
				allowedOriginsAll: false,
			},
		},
		{
			name: "allowedOrigins(*)",
			fields: fields{
				allowedOrigins: []string{"*"},
			},
			want: want{
				allowedOrigins:    []string{"*"},
				allowedOriginsAll: true,
			},
		},
		{
			name: "allowedOrigins",
			fields: fields{
				allowedOrigins: []string{"http://www.test.com"},
			},
			want: want{
				allowedOrigins:    []string{"http://www.test.com"},
				allowedOriginsAll: false,
			},
		},
		{
			name: "wildcard",
			fields: fields{
				allowedOrigins: []string{"http://*.test.com"},
			},
			want: want{
				allowedWildcardOrigins: []wildcard{{prefix: "http://", suffix: ".test.com", len: 16}},
				allowedOrigins:         []string{},
				allowedOriginsAll:      false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cors{
				allowedOrigins:         tt.fields.allowedOrigins,
				allowedWildcardOrigins: tt.fields.allowedWildcardOrigins,
				allowOriginFunc:        tt.fields.allowOriginFunc,
				allowedOriginsAll:      tt.fields.allowedOriginsAll,
			}
			c.computeAllowedOrigins()

			if !reflect.DeepEqual(c.allowedOrigins, tt.want.allowedOrigins) {
				t.Errorf("allowedOrigins = %v, want %v", c.allowedOrigins, tt.want.allowedOrigins)
			}
			if !reflect.DeepEqual(c.allowedWildcardOrigins, tt.want.allowedWildcardOrigins) {
				t.Errorf("allowedWildcardOrigins = %v, want %v", c.allowedWildcardOrigins, tt.want.allowedWildcardOrigins)
			}
			if c.allowedOriginsAll != tt.want.allowedOriginsAll {
				t.Errorf("allowedOriginsAll = %v, want %v", c.allowedOriginsAll, tt.want.allowedOriginsAll)
			}
		})
	}
}

func Test_cors_computeAllowedHeaders(t *testing.T) {
	type fields struct {
		allowedHeaders    []string
		allowedHeadersAll bool
	}
	type want struct {
		allowedHeaders    []string
		allowedHeadersAll bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty",
			want: want{
				allowedHeadersAll: true,
			},
		},
		{
			name: "allowedHeaders(*)",
			fields: fields{
				allowedHeaders: []string{"Header1", "*"},
			},
			want: want{
				allowedHeaders:    []string{},
				allowedHeadersAll: true,
			},
		},
		{
			name: "allowedHeaders(*)",
			fields: fields{
				allowedHeaders: []string{"Header1", "Header2"},
			},
			want: want{
				allowedHeaders:    []string{"Header1", "Header2"},
				allowedHeadersAll: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cors{
				allowedHeaders:    tt.fields.allowedHeaders,
				allowedHeadersAll: tt.fields.allowedHeadersAll,
			}
			c.computeAllowedHeaders()

			if !reflect.DeepEqual(c.allowedHeaders, tt.want.allowedHeaders) {
				t.Errorf("allowedHeaders = %v, want %v", c.allowedHeaders, tt.want.allowedHeaders)
			}
			if c.allowedHeadersAll != tt.want.allowedHeadersAll {
				t.Errorf("allowedHeadersAll = %v, want %v", c.allowedHeadersAll, tt.want.allowedHeadersAll)
			}
		})
	}
}

func Test_cors_computeExposedHeaders(t *testing.T) {
	type fields struct {
		exposedHeaders   []string
		allowCredentials bool
	}
	type want struct {
		exposedHeaders []string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty",
		},
		{
			name:   "allowCredentials=false",
			fields: fields{exposedHeaders: []string{"Header1", "Header2", "*"}},
			want:   want{exposedHeaders: []string{"Header1", "Header2", "*"}},
		},
		{
			name: "allowCredentials=true",
			fields: fields{
				exposedHeaders:   []string{"Header1", "Header2", "*"},
				allowCredentials: true,
			},
			want: want{exposedHeaders: []string{"Header1", "Header2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cors{
				exposedHeaders:   tt.fields.exposedHeaders,
				allowCredentials: tt.fields.allowCredentials,
			}
			c.computeExposedHeaders()

			if !reflect.DeepEqual(c.exposedHeaders, tt.want.exposedHeaders) {
				t.Errorf("exposedHeaders = %v, want %v", c.exposedHeaders, tt.want.exposedHeaders)
			}
		})
	}

}
func Test_cors_isOriginAllowed(t *testing.T) {
	type fields struct {
		allowedOrigins         []string
		allowedWildcardOrigins []wildcard
		allowOriginFunc        func(r *http.Request, origin string) bool
		allowedHeaders         []string
		exposedHeaders         []string
		maxAge                 int
		allowedOriginsAll      bool
		allowedHeadersAll      bool
		allowCredentials       bool
	}
	type args struct {
		r      *http.Request
		origin string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "empty",
			want: false,
		},
		{
			name: "allowOriginFunc(false)",
			fields: fields{
				allowOriginFunc: func(r *http.Request, origin string) bool { return false },
			},
			want: false,
		},
		{
			name: "allowOriginFunc(true)",
			fields: fields{
				allowOriginFunc: func(r *http.Request, origin string) bool { return true },
			},
			want: true,
		},
		{
			name: "allowedOriginsAll=true",
			fields: fields{
				allowedOriginsAll: true,
			},
			want: true,
		},
		{
			name: "allowedOrigins no match",
			args: args{
				origin: "http://www.test.com",
			},
			fields: fields{
				allowedOrigins: []string{"http://www.foo.com"},
			},
			want: false,
		},
		{
			name: "allowedOrigins match",
			args: args{
				origin: "http://www.test.com",
			},
			fields: fields{
				allowedOrigins: []string{"http://www.test.com"},
			},
			want: true,
		},
		{
			name: "allowedWildcardOrigins no match",
			args: args{
				origin: "http://www.foo.com",
			},
			fields: fields{
				allowedWildcardOrigins: []wildcard{{prefix: "http://", suffix: ".bar.com", len: 15}},
			},
			want: false,
		},
		{
			name: "allowedWildcardOrigins match",
			args: args{
				origin: "http://www.bar.com",
			},
			fields: fields{
				allowedWildcardOrigins: []wildcard{{prefix: "http://", suffix: ".bar.com", len: 15}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cors{
				allowedOrigins:         tt.fields.allowedOrigins,
				allowedWildcardOrigins: tt.fields.allowedWildcardOrigins,
				allowOriginFunc:        tt.fields.allowOriginFunc,
				allowedHeaders:         tt.fields.allowedHeaders,
				exposedHeaders:         tt.fields.exposedHeaders,
				maxAge:                 tt.fields.maxAge,
				allowedOriginsAll:      tt.fields.allowedOriginsAll,
				allowedHeadersAll:      tt.fields.allowedHeadersAll,
				allowCredentials:       tt.fields.allowCredentials,
			}
			if got := c.isOriginAllowed(tt.args.r, tt.args.origin); got != tt.want {
				t.Errorf("cors.isOriginAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wildcard_match(t *testing.T) {
	type fields struct {
		prefix string
		suffix string
		len    int
	}
	type args struct {
		origin string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "empty",
			want: false,
		},
		{
			name:   "no match",
			fields: fields{prefix: "http://", suffix: ".bar.com", len: 15},
			args:   args{origin: "http://www.foo.com"},
			want:   false,
		},
		{
			name:   "match",
			fields: fields{prefix: "http://", suffix: ".bar.com", len: 15},
			args:   args{origin: "http://www.bar.com"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := wildcard{
				prefix: tt.fields.prefix,
				suffix: tt.fields.suffix,
				len:    tt.fields.len,
			}
			if got := w.match(tt.args.origin); got != tt.want {
				t.Errorf("wildcard.match() = %v, want %v", got, tt.want)
			}
		})
	}
}
