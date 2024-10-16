package httpserver

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_serveMuxRoute_appendMethod(t *testing.T) {
	tests := []struct {
		name           string
		allowedMethods []string
		method         string
		want           int
	}{
		{
			name:           "not exists",
			allowedMethods: []string{},
			method:         "CUSTOM",
			want:           1,
		},
		{
			name:           "exists 1",
			allowedMethods: []string{"OPTIONS"},
			method:         "CUSTOM",
			want:           2,
		},
		{
			name:           "found",
			allowedMethods: []string{"CUSTOM"},
			method:         "CUSTOM",
			want:           1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smr := &serveMuxRoute{allowedMethods: tt.allowedMethods}
			smr.appendMethod(tt.method)

			if got := len(smr.allowedMethods); got != tt.want {
				t.Errorf("len(allowedMethods) = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serveMuxRoute_addMethod(t *testing.T) {
	tests := []struct {
		name           string
		allowedMethods []string
		method         string
		want           []string
	}{
		{
			name:           "CUSTOM",
			allowedMethods: []string{},
			method:         "CUSTOM",
			want:           []string{"CUSTOM", "OPTIONS"},
		},
		{
			name:           "GET",
			allowedMethods: []string{},
			method:         "GET",
			want:           []string{"GET", "HEAD", "OPTIONS"},
		},

		{
			name:           "empty",
			allowedMethods: []string{},
			method:         "",
			want:           []string{"CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT", "TRACE"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smr := &serveMuxRoute{allowedMethods: tt.allowedMethods}
			smr.addMethod(tt.method)

			if got := smr.allowedMethods; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("serveMuxRoute.allowedMethods = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serveMuxRoute_handlerOptions(t *testing.T) {
	type fields struct {
		allowedMethods []string
		maxAge         int
		cors           *cors
	}
	type want struct {
		statusCode         int
		headerAllow        string
		headerCacheControl string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "without cors, max age = 0",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				maxAge:         0,
			},
			want: want{
				statusCode:         http.StatusNoContent,
				headerAllow:        "GET, HEAD, OPTIONS",
				headerCacheControl: "",
			},
		},
		{
			name: "without cors, max age = 86400",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				maxAge:         86400,
			},
			want: want{
				statusCode:         http.StatusNoContent,
				headerAllow:        "GET, HEAD, OPTIONS",
				headerCacheControl: "max-age=86400",
			},
		},
		{
			name: "with cors",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				cors:           &cors{},
			},
			want: want{
				statusCode:         http.StatusNoContent,
				headerAllow:        "",
				headerCacheControl: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/", nil)

			smr := &serveMuxRoute{
				allowedMethods: tt.fields.allowedMethods,
				maxAge:         tt.fields.maxAge,
				cors:           tt.fields.cors,
			}

			smr.handlerOptions(w, r)

			if got := w.Result().StatusCode; got != tt.want.statusCode {
				t.Errorf("StatusCode = %v, want %v", got, tt.want.statusCode)
			}
			if got := w.Header().Get("Allow"); got != tt.want.headerAllow {
				t.Errorf("Allow = %v, want %v", got, tt.want)
			}
			if got := w.Header().Get("Cache-Control"); got != tt.want.headerCacheControl {
				t.Errorf("Cache-Control = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serveMuxRoute_middlewareCors(t *testing.T) {
	type fields struct {
		allowedMethods []string
		maxAge         int
		cors           *cors
	}
	type args struct {
		method                            string
		headerOrigin                      string
		headerAccessControlRequestMethod  string
		headerAccessControlRequestHeaders string
	}
	type want struct {
		statusCode                          int
		headerAllow                         string
		headerCacheControl                  string
		headerAccessControlAllowOrigin      string
		headerAccessControlAllowCredentials string
		headerAccessControlAllowMethods     string
		headerAccessControlMaxAge           string
		headerAccessControlAllowHeaders     string
		headerAccessControlExposeHeaders    string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "without cors",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
			},
			args: args{
				method: "GET",
			},
			want: want{
				statusCode:                          http.StatusNoContent,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "preflighted, with cors, no origin",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				cors:           &cors{},
			},
			args: args{
				method: "OPTIONS",
			},
			want: want{
				statusCode:                          http.StatusNoContent,
				headerAllow:                         "GET, HEAD, OPTIONS",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "preflighted, with cors, no origin, max age",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				maxAge:         86400,
				cors:           &cors{},
			},
			args: args{
				method: "OPTIONS",
			},
			want: want{
				statusCode:                          http.StatusNoContent,
				headerAllow:                         "GET, HEAD, OPTIONS",
				headerCacheControl:                  "max-age=86400",
				headerAccessControlAllowOrigin:      "",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "preflighted 1, allowedOriginsAll=true, without credentials",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				cors: &cors{
					allowedOriginsAll: true,
				},
			},
			args: args{
				method:       "OPTIONS",
				headerOrigin: "http://test.com",
			},
			want: want{
				statusCode:                          http.StatusBadRequest,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "*",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "GET, HEAD, OPTIONS",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "preflighted 2, allowedOriginsAll=true, without credentials",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				cors: &cors{
					allowedOriginsAll: true,
				},
			},
			args: args{
				method:                           "OPTIONS",
				headerOrigin:                     "http://test.com",
				headerAccessControlRequestMethod: "POST",
			},
			want: want{
				statusCode:                          http.StatusMethodNotAllowed,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "*",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "GET, HEAD, OPTIONS",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "preflighted 3, allowedOriginsAll=true, without credentials",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				cors: &cors{
					allowedOriginsAll: true,
				},
			},
			args: args{
				method:                           "OPTIONS",
				headerOrigin:                     "http://test.com",
				headerAccessControlRequestMethod: "GET",
			},
			want: want{
				statusCode:                          http.StatusNoContent,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "*",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "GET, HEAD, OPTIONS",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "preflighted 1, allowedOriginsAll=false, without credentials",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				cors: &cors{
					allowedOriginsAll: false,
					allowedOrigins:    []string{"http://test.com"},
				},
			},
			args: args{
				method:       "OPTIONS",
				headerOrigin: "http://test.com",
			},
			want: want{
				statusCode:                          http.StatusBadRequest,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "http://test.com",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "GET, HEAD, OPTIONS",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "preflighted 2, allowedOriginsAll=false, without credentials",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				cors: &cors{
					allowedOriginsAll: false,
					allowedOrigins:    []string{"http://test.com"},
				},
			},
			args: args{
				method:                           "OPTIONS",
				headerOrigin:                     "http://test.com",
				headerAccessControlRequestMethod: "POST",
			},
			want: want{
				statusCode:                          http.StatusMethodNotAllowed,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "http://test.com",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "GET, HEAD, OPTIONS",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "preflighted 3, allowedOriginsAll=false, without credentials",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				cors: &cors{
					allowedOriginsAll: false,
					allowedOrigins:    []string{"http://test.com"},
				},
			},
			args: args{
				method:                           "OPTIONS",
				headerOrigin:                     "http://test.com",
				headerAccessControlRequestMethod: "GET",
			},
			want: want{
				statusCode:                          http.StatusNoContent,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "http://test.com",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "GET, HEAD, OPTIONS",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "preflighted, with credentials",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				cors: &cors{
					allowedOriginsAll: true,
					allowCredentials:  true,
				},
			},
			args: args{
				method:       "OPTIONS",
				headerOrigin: "http://test.com",
			},
			want: want{
				statusCode:                          http.StatusBadRequest,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "http://test.com",
				headerAccessControlAllowCredentials: "true",
				headerAccessControlAllowMethods:     "GET, HEAD, OPTIONS",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "preflighted, with max age",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				cors: &cors{
					allowedOriginsAll: true,
					maxAge:            86400,
				},
			},
			args: args{
				method:       "OPTIONS",
				headerOrigin: "http://test.com",
			},
			want: want{
				statusCode:                          http.StatusBadRequest,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "*",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "GET, HEAD, OPTIONS",
				headerAccessControlMaxAge:           "86400",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "preflighted, allowedHeadersAll=true, without credentials",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				cors: &cors{
					allowedOriginsAll: true,
					allowedHeadersAll: true,
				},
			},
			args: args{
				method:       "OPTIONS",
				headerOrigin: "http://test.com",
			},
			want: want{
				statusCode:                          http.StatusBadRequest,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "*",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "GET, HEAD, OPTIONS",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "*",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "preflighted, allowedHeadersAll=true, with credentials",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				cors: &cors{
					allowedOriginsAll: true,
					allowedHeadersAll: true,
					allowCredentials:  true,
				},
			},
			args: args{
				method:                            "OPTIONS",
				headerOrigin:                      "http://test.com",
				headerAccessControlRequestHeaders: "Header-Custom",
			},
			want: want{
				statusCode:                          http.StatusBadRequest,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "http://test.com",
				headerAccessControlAllowCredentials: "true",
				headerAccessControlAllowMethods:     "GET, HEAD, OPTIONS",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "Header-Custom",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "preflighted, allowedHeadersAll=false",
			fields: fields{
				allowedMethods: []string{"GET", "HEAD", "OPTIONS"},
				cors: &cors{
					allowedOriginsAll: true,
					allowedHeadersAll: false,
					allowCredentials:  false,
					allowedHeaders:    []string{"Header-Custom"},
				},
			},
			args: args{
				method:                            "OPTIONS",
				headerOrigin:                      "http://test.com",
				headerAccessControlRequestHeaders: "Header-Custom",
			},
			want: want{
				statusCode:                          http.StatusBadRequest,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "*",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "GET, HEAD, OPTIONS",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "Header-Custom",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "with cors, no origin",
			fields: fields{
				cors: &cors{},
			},
			args: args{
				method: "GET",
			},
			want: want{
				statusCode:                          http.StatusNoContent,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "allowedOriginsAll=true, without credentials",
			fields: fields{
				cors: &cors{
					allowedOriginsAll: true,
				},
			},
			args: args{
				method:       "GET",
				headerOrigin: "http://test.com",
			},
			want: want{
				statusCode:                          http.StatusNoContent,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "*",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "allowedOriginsAll=false, without credentials",
			fields: fields{
				cors: &cors{
					allowedOriginsAll: false,
					allowedOrigins:    []string{"http://test.com"},
				},
			},
			args: args{
				method:       "GET",
				headerOrigin: "http://test.com",
			},
			want: want{
				statusCode:                          http.StatusNoContent,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "http://test.com",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "with credentials",
			fields: fields{
				cors: &cors{
					allowedOriginsAll: true,
					allowCredentials:  true,
				},
			},
			args: args{
				method:       "GET",
				headerOrigin: "http://test.com",
			},
			want: want{
				statusCode:                          http.StatusNoContent,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "http://test.com",
				headerAccessControlAllowCredentials: "true",
				headerAccessControlAllowMethods:     "",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "",
			},
		},
		{
			name: "with expose headers",
			fields: fields{
				cors: &cors{
					allowedOriginsAll: true,
					exposedHeaders:    []string{"Header-Custom"},
				},
			},
			args: args{
				method:       "GET",
				headerOrigin: "http://test.com",
			},
			want: want{
				statusCode:                          http.StatusNoContent,
				headerAllow:                         "",
				headerCacheControl:                  "",
				headerAccessControlAllowOrigin:      "*",
				headerAccessControlAllowCredentials: "",
				headerAccessControlAllowMethods:     "",
				headerAccessControlMaxAge:           "",
				headerAccessControlAllowHeaders:     "",
				headerAccessControlExposeHeaders:    "Header-Custom",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(tt.args.method, "/", nil)
			r.Header.Set("Origin", tt.args.headerOrigin)
			r.Header.Set("Access-Control-Request-Method", tt.args.headerAccessControlRequestMethod)
			r.Header.Set("Access-Control-Request-Headers", tt.args.headerAccessControlRequestHeaders)

			smr := &serveMuxRoute{
				allowedMethods: tt.fields.allowedMethods,
				maxAge:         tt.fields.maxAge,
				cors:           tt.fields.cors,
			}

			next := smr.middlewareCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))
			next.ServeHTTP(w, r)

			if got := w.Result().StatusCode; got != tt.want.statusCode {
				t.Errorf("StatusCode = %v, want %v", got, tt.want.statusCode)
			}
			if got := w.Header().Get("Allow"); got != tt.want.headerAllow {
				t.Errorf("Allow = %v, want %v", got, tt.want.headerAllow)
			}
			if got := w.Header().Get("Cache-Control"); got != tt.want.headerCacheControl {
				t.Errorf("Cache-Control = %v, want %v", got, tt.want.headerCacheControl)
			}
			if got := w.Header().Get("Access-Control-Allow-Origin"); got != tt.want.headerAccessControlAllowOrigin {
				t.Errorf("Access-Control-Allow-Origin = %v, want %v", got, tt.want.headerAccessControlAllowOrigin)
			}
			if got := w.Header().Get("Access-Control-Allow-Credentials"); got != tt.want.headerAccessControlAllowCredentials {
				t.Errorf("Access-Control-Allow-Credentials = %v, want %v", got, tt.want.headerAccessControlAllowCredentials)
			}
			if got := w.Header().Get("Access-Control-Allow-Methods"); got != tt.want.headerAccessControlAllowMethods {
				t.Errorf("Access-Control-Allow-Methods = %v, want %v", got, tt.want.headerAccessControlAllowMethods)
			}
			if got := w.Header().Get("Access-Control-Max-Age"); got != tt.want.headerAccessControlMaxAge {
				t.Errorf("Access-Control-Max-Age = %v, want %v", got, tt.want.headerAccessControlMaxAge)
			}
			if got := w.Header().Get("Access-Control-Allow-Headers"); got != tt.want.headerAccessControlAllowHeaders {
				t.Errorf("Access-Control-Allow-Headers = %v, want %v", got, tt.want.headerAccessControlAllowHeaders)
			}
			if got := w.Header().Get("Access-Control-Expose-Headers"); got != tt.want.headerAccessControlExposeHeaders {
				t.Errorf("Access-Control-Expose-Headers = %v, want %v", got, tt.want.headerAccessControlExposeHeaders)
			}
		})
	}
}
