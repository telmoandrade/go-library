package httpserver

import (
	"net/http"
	"reflect"
	"testing"
)

func TestWithCors(t *testing.T) {
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
		opts []OptionCors
		want want
	}{
		{
			name: "without option",
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
			name: "WithAllowedOrigins(any)",
			opts: []OptionCors{
				WithAllowedOrigins("http://test.com"),
			},
			want: want{
				allowedOrigins:         []string{"http://test.com"},
				allowedWildcardOrigins: []wildcard{},
				allowedHeaders:         []string{},
				exposedHeaders:         []string{},
				maxAge:                 86400,
				allowedOriginsAll:      false,
				allowedHeadersAll:      true,
				allowCredentials:       false,
			},
		},
		{
			name: "WithAllowOriginFunc()",
			opts: []OptionCors{
				WithAllowOriginFunc(func(r *http.Request, origin string) bool { return true }),
			},
			want: want{
				allowedOrigins:         []string{},
				allowedWildcardOrigins: []wildcard{},
				allowedHeaders:         []string{},
				exposedHeaders:         []string{},
				maxAge:                 86400,
				allowedOriginsAll:      false,
				allowedHeadersAll:      true,
				allowCredentials:       false,
			},
		},
		{
			name: "WithAllowedOrigins(Wildcard)",
			opts: []OptionCors{
				WithAllowedOrigins("http://*.bar.com"),
			},
			want: want{
				allowedOrigins:         []string{},
				allowedWildcardOrigins: []wildcard{{prefix: "http://", suffix: ".bar.com", len: 15}},
				allowedHeaders:         []string{},
				exposedHeaders:         []string{},
				maxAge:                 86400,
				allowedOriginsAll:      false,
				allowedHeadersAll:      true,
				allowCredentials:       false,
			},
		},
		{
			name: "WithAllowedOrigins(*)",
			opts: []OptionCors{
				WithAllowedOrigins("*"),
			},
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
			name: "WithAllowedHeaders(any)",
			opts: []OptionCors{
				WithAllowedHeaders("test"),
			},
			want: want{
				allowedOrigins:         []string{},
				allowedWildcardOrigins: []wildcard{},
				allowedHeaders:         []string{"Test"},
				exposedHeaders:         []string{},
				maxAge:                 86400,
				allowedOriginsAll:      true,
				allowedHeadersAll:      false,
				allowCredentials:       false,
			},
		},
		{
			name: "WithAllowedHeaders(*)",
			opts: []OptionCors{
				WithAllowedHeaders("*"),
			},
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
			name: "WithExposedHeaders",
			opts: []OptionCors{
				WithExposedHeaders("test"),
			},
			want: want{
				allowedOrigins:         []string{},
				allowedWildcardOrigins: []wildcard{},
				allowedHeaders:         []string{},
				exposedHeaders:         []string{"Test"},
				maxAge:                 86400,
				allowedOriginsAll:      true,
				allowedHeadersAll:      true,
				allowCredentials:       false,
			},
		},
		{
			name: "WithAllowCredentials(true)",
			opts: []OptionCors{
				WithAllowCredentials(true),
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
		{
			name: "WithCorsMaxAge(0)",
			opts: []OptionCors{
				WithCorsMaxAge(0),
			},
			want: want{
				allowedOrigins:         []string{},
				allowedWildcardOrigins: []wildcard{},
				allowedHeaders:         []string{},
				exposedHeaders:         []string{},
				maxAge:                 0,
				allowedOriginsAll:      true,
				allowedHeadersAll:      true,
				allowCredentials:       false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := NewServeMux(WithCors(tt.opts...))
			got := mux.options.cors
			if !reflect.DeepEqual(got.allowedOrigins, tt.want.allowedOrigins) {
				t.Errorf("allowedOrigins = %v, want %v", got.allowedOrigins, tt.want.allowedOrigins)
			}
			if !reflect.DeepEqual(got.allowedWildcardOrigins, tt.want.allowedWildcardOrigins) {
				t.Errorf("allowedWildcardOrigins = %v, want %v", got.allowedWildcardOrigins, tt.want.allowedWildcardOrigins)
			}
			if !reflect.DeepEqual(got.allowedHeaders, tt.want.allowedHeaders) {
				t.Errorf("allowedHeaders = %v, want %v", got.allowedHeaders, tt.want.allowedHeaders)
			}
			if !reflect.DeepEqual(got.exposedHeaders, tt.want.exposedHeaders) {
				t.Errorf("exposedHeaders = %v, want %v", got.exposedHeaders, tt.want.exposedHeaders)
			}
			if got.maxAge != tt.want.maxAge {
				t.Errorf("maxAge = %v, want %v", got.maxAge, tt.want.maxAge)
			}
			if got.allowedOriginsAll != tt.want.allowedOriginsAll {
				t.Errorf("allowedOriginsAll = %v, want %v", got.allowedOriginsAll, tt.want.allowedOriginsAll)
			}
			if got.allowedHeadersAll != tt.want.allowedHeadersAll {
				t.Errorf("allowedHeadersAll = %v, want %v", got.allowedHeadersAll, tt.want.allowedHeadersAll)
			}
			if got.allowCredentials != tt.want.allowCredentials {
				t.Errorf("allowCredentials = %v, want %v", got.allowCredentials, tt.want.allowCredentials)
			}
		})
	}
}

func Test_cors_isOriginAllowed(t *testing.T) {
	type args struct {
		origin string
	}
	tests := []struct {
		name string
		args args
		opts []OptionCors
		want bool
	}{
		{
			name: "without option",
			args: args{origin: "http://test.com"},
			want: true,
		},
		{
			name: "WithAllowedOrigins(ok)",
			args: args{origin: "http://test.com"},
			opts: []OptionCors{
				WithAllowedOrigins(
					"http://foo.com",
					"http://test.com",
				),
			},
			want: true,
		},
		{
			name: "WithAllowedOrigins(invalid)",
			args: args{origin: "http://test.com"},
			opts: []OptionCors{
				WithAllowedOrigins(
					"http://foo.com",
					"http://bar.com",
				),
			},
			want: false,
		},
		{
			name: "WithAllowedOrigins(wildcard)",
			args: args{origin: "http://www.bar.com"},
			opts: []OptionCors{
				WithAllowedOrigins(
					"http://*.foo.com",
					"http://*.bar.com",
				),
			},
			want: true,
		},
		{
			name: "WithAllowOriginFunc(true)",
			args: args{origin: "http://test.com"},
			opts: []OptionCors{
				WithAllowOriginFunc(func(r *http.Request, origin string) bool { return true }),
			},
			want: true,
		},
		{
			name: "WithAllowOriginFunc(false)",
			args: args{origin: "http://test.com"},
			opts: []OptionCors{
				WithAllowOriginFunc(func(r *http.Request, origin string) bool { return false }),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := NewServeMux(WithCors(tt.opts...))
			cors := mux.options.cors

			got := cors.isOriginAllowed(nil, tt.args.origin)
			if got != tt.want {
				t.Errorf("cors.isOriginAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wildcard_match(t *testing.T) {
	type fields struct {
		prefix string
		suffix string
	}
	tests := []struct {
		name   string
		fields fields
		args   string
		want   bool
	}{
		{
			name: "ok",
			fields: fields{
				prefix: "http://",
				suffix: ".bar.com",
			},
			args: "http://www.bar.com",
			want: true,
		},
		{
			name: "invalid size",
			fields: fields{
				prefix: "http://",
				suffix: ".bar.com",
			},
			args: "http://.bar.com",
			want: false,
		},
		{
			name: "invalid prefix",
			fields: fields{
				prefix: "http://",
				suffix: ".bar.com",
			},
			args: "https://www.bar.com",
			want: false,
		},
		{
			name: "invalid suffix",
			fields: fields{
				prefix: "http://",
				suffix: ".bar.com",
			},
			args: "https://www.foo.com",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := wildcard{
				prefix: tt.fields.prefix,
				suffix: tt.fields.suffix,
				len:    len(tt.fields.prefix) + len(tt.fields.suffix),
			}
			if got := w.match(tt.args); got != tt.want {
				t.Errorf("wildcard.match() = %v, want %v", got, tt.want)
			}
		})
	}
}
