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
			method:         http.MethodPost,
			want:           1,
		},
		{
			name:           "exists 1",
			allowedMethods: []string{http.MethodOptions},
			method:         http.MethodPost,
			want:           2,
		},
		{
			name:           "found",
			allowedMethods: []string{http.MethodPost},
			method:         http.MethodPost,
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
			name:           "empty",
			allowedMethods: []string{},
			method:         "",
			want:           []string{},
		},
		{
			name:           http.MethodGet,
			allowedMethods: []string{},
			method:         http.MethodGet,
			want:           []string{http.MethodGet, http.MethodHead},
		},
		{
			name:           http.MethodPost,
			allowedMethods: []string{http.MethodPost},
			method:         "",
			want:           []string{http.MethodPost},
		},
		{
			name:           "CUSTOM",
			allowedMethods: []string{},
			method:         "CUSTOM",
			want:           []string{"CUSTOM"},
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

func Test_serveMuxRoute_headerOptions(t *testing.T) {
	type fields struct {
		allowedMethods       []string
		handlerOptionsMaxAge int
	}
	type want struct {
		headerAllow        string
		headerCacheControl string
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
			name: "Allow",
			fields: fields{
				allowedMethods:       []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				handlerOptionsMaxAge: 0,
			},
			want: want{
				headerAllow:        "GET, HEAD, OPTIONS",
				headerCacheControl: "",
			},
		},
		{
			name: "Allow + CacheControl",
			fields: fields{
				allowedMethods:       []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				handlerOptionsMaxAge: 86400,
			},
			want: want{
				headerAllow:        "GET, HEAD, OPTIONS",
				headerCacheControl: "public, max-age=86400",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			smr := &serveMuxRoute{
				allowedMethods:       tt.fields.allowedMethods,
				handlerOptionsMaxAge: tt.fields.handlerOptionsMaxAge,
			}
			smr.headerOptions(w)

			if got := w.Header().Get("Allow"); got != tt.want.headerAllow {
				t.Errorf("Allow = %v, want %v", got, tt.want.headerAllow)
			}
			if got := w.Header().Get("Cache-Control"); got != tt.want.headerCacheControl {
				t.Errorf("Cache-Control = %v, want %v", got, tt.want.headerCacheControl)
			}
		})
	}
}

func Test_serveMuxRoute_middlewareMethodNotAllowed(t *testing.T) {
	type fields struct {
		allowedMethods       []string
		handlerOptionsMaxAge int
	}
	type want struct {
		headerAllow        string
		headerCacheControl string
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
			name: "Allow",
			fields: fields{
				allowedMethods:       []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				handlerOptionsMaxAge: 0,
			},
			want: want{
				headerAllow:        "GET, HEAD, OPTIONS",
				headerCacheControl: "",
			},
		},
		{
			name: "Allow + CacheControl",
			fields: fields{
				allowedMethods:       []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				handlerOptionsMaxAge: 86400,
			},
			want: want{
				headerAllow:        "GET, HEAD, OPTIONS",
				headerCacheControl: "public, max-age=86400",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smr := &serveMuxRoute{
				allowedMethods:       tt.fields.allowedMethods,
				handlerOptionsMaxAge: tt.fields.handlerOptionsMaxAge,
			}

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			m := smr.middlewareMethodNotAllowed(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			m.ServeHTTP(w, r)

			if got := w.Header().Get("Allow"); got != tt.want.headerAllow {
				t.Errorf("Allow = %v, want %v", got, tt.want.headerAllow)
			}
			if got := w.Header().Get("Cache-Control"); got != tt.want.headerCacheControl {
				t.Errorf("Cache-Control = %v, want %v", got, tt.want.headerCacheControl)
			}
		})
	}
}

func Test_serveMuxRoute_headerCorsAll(t *testing.T) {
	type fields struct {
		cors *cors
	}
	type want struct {
		headerAccessControlAllowOrigin      string
		headerAccessControlAllowCredentials string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "OriginsAll: false, Credentials: false",
			fields: fields{
				cors: &cors{
					allowedOriginsAll: false,
					allowCredentials:  false,
				},
			},
			want: want{
				headerAccessControlAllowOrigin:      "http://test.com",
				headerAccessControlAllowCredentials: "",
			},
		},
		{
			name: "OriginsAll: true, Credentials: false",
			fields: fields{
				cors: &cors{
					allowedOriginsAll: true,
					allowCredentials:  false,
				},
			},
			want: want{
				headerAccessControlAllowOrigin:      "*",
				headerAccessControlAllowCredentials: "",
			},
		},
		{
			name: "OriginsAll: false, Credentials: true",
			fields: fields{
				cors: &cors{
					allowedOriginsAll: false,
					allowCredentials:  true,
				},
			},
			want: want{
				headerAccessControlAllowOrigin:      "http://test.com",
				headerAccessControlAllowCredentials: "true",
			},
		},
		{
			name: "OriginsAll: true, Credentials: true",
			fields: fields{
				cors: &cors{
					allowedOriginsAll: true,
					allowCredentials:  true,
				},
			},
			want: want{
				headerAccessControlAllowOrigin:      "http://test.com",
				headerAccessControlAllowCredentials: "true",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			smr := &serveMuxRoute{
				cors: tt.fields.cors,
			}

			smr.headerCorsAll(w, "http://test.com")

			if got := w.Header().Get("Access-Control-Allow-Origin"); got != tt.want.headerAccessControlAllowOrigin {
				t.Errorf("Access-Control-Allow-Origin = %v, want %v", got, tt.want.headerAccessControlAllowOrigin)
			}
			if got := w.Header().Get("Access-Control-Allow-Credentials"); got != tt.want.headerAccessControlAllowCredentials {
				t.Errorf("Access-Control-Allow-Credentials = %v, want %v", got, tt.want.headerAccessControlAllowCredentials)
			}
		})
	}
}

func Test_serveMuxRoute_headerCorsPreflightAllowHeaders(t *testing.T) {
	type fields struct {
		cors *cors
	}
	type args struct {
		headerAccessControlRequestHeaders string
	}
	type want struct {
		headerAccessControlAllowHeaders string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "HeadersAll: false, Credentials: false",
			fields: fields{
				cors: &cors{
					allowedHeadersAll: false,
					allowCredentials:  false,
					allowedHeaders:    []string{"Header1"},
				},
			},
			args: args{
				headerAccessControlRequestHeaders: "Header1, Header2",
			},
			want: want{
				headerAccessControlAllowHeaders: "Header1",
			},
		},
		{
			name: "HeadersAll: true, Credentials: false",
			fields: fields{
				cors: &cors{
					allowedHeadersAll: true,
					allowCredentials:  false,
				},
			},
			args: args{
				headerAccessControlRequestHeaders: "Header1, Header2",
			},
			want: want{
				headerAccessControlAllowHeaders: "*",
			},
		},
		{
			name: "HeadersAll: false, Credentials: true",
			fields: fields{
				cors: &cors{
					allowedHeadersAll: false,
					allowCredentials:  true,
					allowedHeaders:    []string{"Header1"},
				},
			},
			args: args{
				headerAccessControlRequestHeaders: "Header1, Header2",
			},
			want: want{
				headerAccessControlAllowHeaders: "Header1",
			},
		},
		{
			name: "HeadersAll: true, Credentials: true",
			fields: fields{
				cors: &cors{
					allowedHeadersAll: true,
					allowCredentials:  true,
				},
			},
			args: args{
				headerAccessControlRequestHeaders: "Header1, Header2",
			},
			want: want{
				headerAccessControlAllowHeaders: "Header1, Header2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			r, _ := http.NewRequest(http.MethodOptions, "/", nil)
			r.Header.Set("Access-Control-Request-Headers", tt.args.headerAccessControlRequestHeaders)

			smr := &serveMuxRoute{
				cors: tt.fields.cors,
			}

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				smr.headerCorsPreflightAllowHeaders(w, r)
				w.WriteHeader(http.StatusNoContent)
			})
			next.ServeHTTP(w, r)

			if got := w.Header().Get("Access-Control-Allow-Headers"); got != tt.want.headerAccessControlAllowHeaders {
				t.Errorf("Access-Control-Allow-Headers = %v, want %v", got, tt.want.headerAccessControlAllowHeaders)
			}
		})
	}
}

func Test_serveMuxRoute_headerCorsPreflight(t *testing.T) {
	type fields struct {
		allowedMethods []string
		cors           *cors
	}
	type args struct {
		headerAccessControlRequestMethod string
	}
	type want struct {
		statusCode                      int
		headerError                     string
		headerAccessControlAllowMethods string
		headerAccessControlMaxAge       string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "StatusBadRequest",
			fields: fields{
				allowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				cors:           &cors{},
			},
			args: args{
				headerAccessControlRequestMethod: "",
			},
			want: want{
				statusCode:                      http.StatusBadRequest,
				headerError:                     "Access-Control-Request-Method not found",
				headerAccessControlAllowMethods: "GET, HEAD, OPTIONS",
				headerAccessControlMaxAge:       "",
			},
		},
		{
			name: "StatusMethodNotAllowed",
			fields: fields{
				allowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				cors:           &cors{},
			},
			args: args{
				headerAccessControlRequestMethod: "POST",
			},
			want: want{
				statusCode:                      http.StatusMethodNotAllowed,
				headerError:                     "",
				headerAccessControlAllowMethods: "GET, HEAD, OPTIONS",
				headerAccessControlMaxAge:       "",
			},
		},
		{
			name: "StatusNoContent",
			fields: fields{
				allowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				cors:           &cors{},
			},
			args: args{
				headerAccessControlRequestMethod: "GET",
			},
			want: want{
				statusCode:                      http.StatusNoContent,
				headerError:                     "",
				headerAccessControlAllowMethods: "GET, HEAD, OPTIONS",
				headerAccessControlMaxAge:       "",
			},
		},
		{
			name: "with maxAge",
			fields: fields{
				allowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				cors: &cors{
					maxAge: 86400,
				},
			},
			args: args{
				headerAccessControlRequestMethod: "GET",
			},
			want: want{
				statusCode:                      http.StatusNoContent,
				headerError:                     "",
				headerAccessControlAllowMethods: "GET, HEAD, OPTIONS",
				headerAccessControlMaxAge:       "86400",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			r, _ := http.NewRequest(http.MethodOptions, "/", nil)
			r.Header.Set("Access-Control-Request-Method", tt.args.headerAccessControlRequestMethod)

			smr := &serveMuxRoute{
				allowedMethods: tt.fields.allowedMethods,
				cors:           tt.fields.cors,
			}

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				smr.headerCorsPreflight(w, r)
				w.WriteHeader(http.StatusNoContent)
			})
			next.ServeHTTP(w, r)

			if got := w.Result().StatusCode; got != tt.want.statusCode {
				t.Errorf("StatusCode = %v, want %v", got, tt.want.statusCode)
			}

			if got := w.Header().Get("Error"); got != tt.want.headerError {
				t.Errorf("Error = %v, want %v", got, tt.want.headerError)
			}
			if got := w.Header().Get("Access-Control-Allow-Methods"); got != tt.want.headerAccessControlAllowMethods {
				t.Errorf("Access-Control-Allow-Methods = %v, want %v", got, tt.want.headerAccessControlAllowMethods)
			}
			if got := w.Header().Get("Access-Control-Max-Age"); got != tt.want.headerAccessControlMaxAge {
				t.Errorf("Access-Control-Max-Age = %v, want %v", got, tt.want.headerAccessControlMaxAge)
			}
		})
	}
}

func Test_serveMuxRoute_headerCorsActualRequest(t *testing.T) {
	type fields struct {
		cors *cors
	}
	type want struct {
		headerAccessControlExposeHeaders string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty",
			fields: fields{
				cors: &cors{},
			},
			want: want{
				headerAccessControlExposeHeaders: "",
			},
		},
		{
			name: "exposedHeaders",
			fields: fields{
				cors: &cors{
					exposedHeaders: []string{"Header1", "Header2"},
				},
			},
			want: want{
				headerAccessControlExposeHeaders: "Header1, Header2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			smr := &serveMuxRoute{
				cors: tt.fields.cors,
			}
			smr.headerCorsActualRequest(w)

			if got := w.Header().Get("Access-Control-Expose-Headers"); got != tt.want.headerAccessControlExposeHeaders {
				t.Errorf("Access-Control-Expose-Headers = %v, want %v", got, tt.want.headerAccessControlExposeHeaders)
			}
		})
	}
}

func Test_serveMuxRoute_middlewareCors(t *testing.T) {
	type fields struct {
		allowedMethods []string
		cors           *cors
	}
	type args struct {
		method       string
		headerOrigin string
	}
	type want struct {
		headerAllow                      string
		headerVary                       []string
		headerAccessControlAllowOrigin   string
		headerAccessControlAllowMethods  string
		headerAccessControlExposeHeaders string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "OPTIONS, empty Origin",
			fields: fields{
				allowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				cors:           &cors{},
			},
			args: args{
				method:       http.MethodOptions,
				headerOrigin: "",
			},
			want: want{
				headerAllow:                      "GET, HEAD, OPTIONS",
				headerVary:                       []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
				headerAccessControlAllowOrigin:   "",
				headerAccessControlAllowMethods:  "",
				headerAccessControlExposeHeaders: "",
			},
		},
		{
			name: "OPTIONS, invalid Origin",
			fields: fields{
				allowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				cors:           &cors{},
			},
			args: args{
				method:       http.MethodOptions,
				headerOrigin: "http://invalid.com",
			},
			want: want{
				headerAllow:                      "GET, HEAD, OPTIONS",
				headerVary:                       []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
				headerAccessControlAllowOrigin:   "",
				headerAccessControlAllowMethods:  "",
				headerAccessControlExposeHeaders: "",
			},
		},
		{
			name: "OPTIONS, Origin",
			fields: fields{
				allowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				cors: &cors{
					allowedOriginsAll: true,
				},
			},
			args: args{
				method:       http.MethodOptions,
				headerOrigin: "http://test.com",
			},
			want: want{
				headerAllow:                      "",
				headerVary:                       []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
				headerAccessControlAllowOrigin:   "http://test.com",
				headerAccessControlAllowMethods:  "GET, HEAD, OPTIONS",
				headerAccessControlExposeHeaders: "",
			},
		},
		{
			name: "GET, empty Origin",
			fields: fields{
				allowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				cors: &cors{
					allowedOriginsAll: true,
					exposedHeaders:    []string{"Header1"},
				},
			},
			args: args{
				method:       http.MethodGet,
				headerOrigin: "",
			},
			want: want{
				headerAllow:                      "",
				headerVary:                       []string{"Origin"},
				headerAccessControlAllowOrigin:   "",
				headerAccessControlAllowMethods:  "",
				headerAccessControlExposeHeaders: "",
			},
		},
		{
			name: "GET, invalid Origin",
			fields: fields{
				allowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				cors: &cors{
					allowedOriginsAll: true,
					exposedHeaders:    []string{"Header1"},
				},
			},
			args: args{
				method:       http.MethodGet,
				headerOrigin: "http://invalid.com",
			},
			want: want{
				headerAllow:                      "",
				headerVary:                       []string{"Origin"},
				headerAccessControlAllowOrigin:   "",
				headerAccessControlAllowMethods:  "",
				headerAccessControlExposeHeaders: "",
			},
		},
		{
			name: "GET, Origin",
			fields: fields{
				allowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				cors: &cors{
					allowedOriginsAll: true,
					exposedHeaders:    []string{"Header1"},
				},
			},
			args: args{
				method:       http.MethodGet,
				headerOrigin: "http://test.com",
			},
			want: want{
				headerAllow:                      "",
				headerVary:                       []string{"Origin"},
				headerAccessControlAllowOrigin:   "http://test.com",
				headerAccessControlAllowMethods:  "",
				headerAccessControlExposeHeaders: "Header1",
			},
		},
		{
			name: "POST, empty Origin",
			fields: fields{
				allowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				cors: &cors{
					allowedOriginsAll: true,
					exposedHeaders:    []string{"Header1"},
				},
			},
			args: args{
				method:       http.MethodPost,
				headerOrigin: "",
			},
			want: want{
				headerAllow:                      "GET, HEAD, OPTIONS",
				headerVary:                       []string{"Origin"},
				headerAccessControlAllowOrigin:   "",
				headerAccessControlAllowMethods:  "",
				headerAccessControlExposeHeaders: "",
			},
		},
		{
			name: "POST, invalid Origin",
			fields: fields{
				allowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				cors: &cors{
					allowedOriginsAll: true,
					exposedHeaders:    []string{"Header1"},
				},
			},
			args: args{
				method:       http.MethodPost,
				headerOrigin: "http://invalid.com",
			},
			want: want{
				headerAllow:                      "GET, HEAD, OPTIONS",
				headerVary:                       []string{"Origin"},
				headerAccessControlAllowOrigin:   "",
				headerAccessControlAllowMethods:  "",
				headerAccessControlExposeHeaders: "",
			},
		},
		{
			name: "POST, Origin",
			fields: fields{
				allowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
				cors: &cors{
					allowedOriginsAll: true,
					exposedHeaders:    []string{"Header1"},
				},
			},
			args: args{
				method:       http.MethodPost,
				headerOrigin: "http://test.com",
			},
			want: want{
				headerAllow:                      "",
				headerVary:                       []string{"Origin"},
				headerAccessControlAllowOrigin:   "http://test.com",
				headerAccessControlAllowMethods:  "",
				headerAccessControlExposeHeaders: "Header1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(tt.args.method, "/", nil)
			r.Header.Set("Origin", tt.args.headerOrigin)

			smr := &serveMuxRoute{
				allowedMethods: tt.fields.allowedMethods,
				cors: &cors{
					allowedOrigins: []string{"http://test.com"},
					exposedHeaders: []string{"Header1"},
				},
			}

			next := smr.middlewareCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			next.ServeHTTP(w, r)

			if got := w.Header().Get("Allow"); got != tt.want.headerAllow {
				t.Errorf("Allow = %v, want %v", got, tt.want.headerAllow)
			}
			if got := w.Header().Values("Vary"); !reflect.DeepEqual(got, tt.want.headerVary) {
				t.Errorf("Vary = %v, want %v", got, tt.want.headerVary)
			}
			if got := w.Header().Get("Access-Control-Allow-Origin"); got != tt.want.headerAccessControlAllowOrigin {
				t.Errorf("Access-Control-Allow-Origin = %v, want %v", got, tt.want.headerAccessControlAllowOrigin)
			}
			if got := w.Header().Get("Access-Control-Allow-Methods"); got != tt.want.headerAccessControlAllowMethods {
				t.Errorf("Access-Control-Allow-Methods = %v, want %v", got, tt.want.headerAccessControlAllowMethods)
			}
			if got := w.Header().Get("Access-Control-Expose-Headers"); got != tt.want.headerAccessControlExposeHeaders {
				t.Errorf("Access-Control-Expose-Headers = %v, want %v", got, tt.want.headerAccessControlExposeHeaders)
			}
		})
	}
}
