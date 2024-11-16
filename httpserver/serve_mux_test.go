package httpserver

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"testing"
)

func handlerId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		w.Write([]byte(fmt.Sprintf("%s ID:%s", r.Pattern, id)))
	}
}

func middlewareEmpty(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func middlewareWrite(value string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("BEFORE:" + value + " -> "))
			next.ServeHTTP(w, r)
			w.Write([]byte(" -> AFTER:" + value))
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, payload io.Reader) (int, http.Header, string) {
	req, err := http.NewRequest(method, ts.URL+path, payload)
	if err != nil {
		t.Fatal(err)
		return 0, nil, ""
	}

	req.Host = "www.test.com"

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return 0, nil, ""
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return 0, nil, ""
	}
	defer resp.Body.Close()

	return resp.StatusCode, resp.Header, string(respBody)
}

func TestNewServeMux(t *testing.T) {
	type args struct {
		opts []OptionServeMux
	}
	type want struct {
		middlewares  []func(http.Handler) http.Handler
		patternRoute *patternRoute
		routes       map[string]*serveMuxRoute
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "without option",
			want: want{
				middlewares:  []func(http.Handler) http.Handler{},
				patternRoute: newPatternRoute(""),
				routes:       map[string]*serveMuxRoute{},
			},
		},
		{
			name: "with option",
			args: args{
				opts: []OptionServeMux{
					WithHandlerOptionsMaxAge(0),
				},
			},
			want: want{
				middlewares:  []func(http.Handler) http.Handler{},
				patternRoute: newPatternRoute(""),
				routes:       map[string]*serveMuxRoute{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux(
				tt.args.opts...,
			)
			mux := muxInterface.(*serveMux)

			if mux.serveMux == nil {
				t.Errorf("mux.serveMux = nil, want != nil")
			}

			if got := mux.middlewares; !reflect.DeepEqual(got, tt.want.middlewares) {
				t.Errorf("mux.middlewares = %v, want %v", got, tt.want.middlewares)
			}

			if got := mux.patternRoute; !reflect.DeepEqual(got, tt.want.patternRoute) {
				t.Errorf("mux.patternRoute = %v, want %v", got, tt.want.patternRoute)
			}

			if got := mux.routes; !reflect.DeepEqual(got, tt.want.routes) {
				t.Errorf("mux.routes = %v, want %v", got, tt.want.routes)
			}
		})
	}
}

func Test_serveMux_Use(t *testing.T) {
	type args struct {
		middlewares []func(http.Handler) http.Handler
	}
	type want struct {
		middlewares int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty",
			want: want{
				middlewares: 0,
			},
		},
		{
			name: "middleware",
			args: args{
				middlewares: []func(http.Handler) http.Handler{middlewareEmpty},
			},
			want: want{
				middlewares: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux()
			mux := muxInterface.(*serveMux)

			mux.Use(tt.args.middlewares...)

			if got := len(mux.middlewares); got != tt.want.middlewares {
				t.Errorf("mux.middlewares.len = %v, want %v", got, tt.want.middlewares)
			}
		})
	}
}

func Test_serveMux_With(t *testing.T) {
	type args struct {
		middlewares []func(http.Handler) http.Handler
	}
	type want struct {
		middlewares int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty",
			want: want{
				middlewares: 0,
			},
		},
		{
			name: "middlewareEmpty",
			args: args{
				middlewares: []func(http.Handler) http.Handler{middlewareEmpty},
			},
			want: want{
				middlewares: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux()
			mux := muxInterface.(*serveMux)

			muxWithInterface := mux.With(tt.args.middlewares...)
			muxWith := muxWithInterface.(*serveMux)

			if got := len(muxWith.middlewares); got != tt.want.middlewares {
				t.Errorf("mux.middlewares.len = %v, want %v", got, tt.want.middlewares)
			}
		})
	}
}

func Test_serveMux_Group(t *testing.T) {
	type args struct {
		patternLevel1 string
		patternLevel2 string
	}
	type want struct {
		patternLevel0 string
		patternLevel1 string
		patternLevel2 string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "\\admin \\user",
			args: args{
				patternLevel1: "/admin",
				patternLevel2: "/user",
			},
			want: want{
				patternLevel0: "/",
				patternLevel1: "/admin",
				patternLevel2: "/admin/user",
			},
		},
		{
			name: "\\user \\{id}",
			args: args{
				patternLevel1: "/user",
				patternLevel2: "/{id}",
			},
			want: want{
				patternLevel0: "/",
				patternLevel1: "/user",
				patternLevel2: "/user/{id}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxServeMux := NewServeMux()
			muxLevel0 := muxServeMux.(*serveMux)

			muxGroupLevel1 := muxServeMux.Group(tt.args.patternLevel1)
			muxLevel1 := muxGroupLevel1.(*serveMux)

			muxGroupLevel2 := muxGroupLevel1.Group(tt.args.patternLevel2)
			muxLevel2 := muxGroupLevel2.(*serveMux)

			if got := muxLevel0.patternRoute.String(); !reflect.DeepEqual(got, tt.want.patternLevel0) {
				t.Errorf("level 0 = %v, want %v", got, tt.want.patternLevel0)
			}
			if got := muxLevel1.patternRoute.String(); !reflect.DeepEqual(got, tt.want.patternLevel1) {
				t.Errorf("level 1 = %v, want %v", got, tt.want.patternLevel1)
			}
			if got := muxLevel2.patternRoute.String(); !reflect.DeepEqual(got, tt.want.patternLevel2) {
				t.Errorf("level 2 = %v, want %v", got, tt.want.patternLevel2)
			}
		})
	}
}

func Test_serveMux_Route(t *testing.T) {
	type args struct {
		patternLevel1 string
		patternLevel2 string
	}
	type want struct {
		patternLevel0 string
		patternLevel1 string
		patternLevel2 string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "\\admin \\user",
			args: args{
				patternLevel1: "/admin",
				patternLevel2: "/user",
			},
			want: want{
				patternLevel0: "/",
				patternLevel1: "/admin",
				patternLevel2: "/admin/user",
			},
		},
		{
			name: "\\user \\{id}",
			args: args{
				patternLevel1: "/user",
				patternLevel2: "/{id}",
			},
			want: want{
				patternLevel0: "/",
				patternLevel1: "/user",
				patternLevel2: "/user/{id}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxServeMux := NewServeMux()
			muxLevel0 := muxServeMux.(*serveMux)
			var muxLevel1 *serveMux
			var muxLevel2 *serveMux

			muxServeMux.Route(tt.args.patternLevel1, func(muxGroupLevel1 Router) {
				muxLevel1 = muxGroupLevel1.(*serveMux)
				muxGroupLevel1.Route(tt.args.patternLevel2, func(muxGroupLevel2 Router) {
					muxLevel2 = muxGroupLevel2.(*serveMux)
				})
			})

			if got := muxLevel0.patternRoute.String(); !reflect.DeepEqual(got, tt.want.patternLevel0) {
				t.Errorf("level 0 = %v, want %v", got, tt.want.patternLevel0)
			}
			if got := muxLevel1.patternRoute.String(); !reflect.DeepEqual(got, tt.want.patternLevel1) {
				t.Errorf("level 1 = %v, want %v", got, tt.want.patternLevel1)
			}
			if got := muxLevel2.patternRoute.String(); !reflect.DeepEqual(got, tt.want.patternLevel2) {
				t.Errorf("level 2 = %v, want %v", got, tt.want.patternLevel2)
			}
		})
	}
}

func Test_validateHandler(t *testing.T) {
	type args struct {
		handler http.Handler
	}
	type want struct {
		panicked   bool
		panicError string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "invalid http.Handler",
			args: args{
				handler: nil,
			},
			want: want{
				panicked:   true,
				panicError: "httpserver: nil handler",
			},
		},
		{
			name: "invalid http.HandlerFunc",
			args: args{
				handler: http.HandlerFunc(nil),
			},
			want: want{
				panicked:   true,
				panicError: "httpserver: nil handler",
			},
		},
		{
			name: "valid",
			args: args{
				handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			},
			want: want{
				panicked:   false,
				panicError: "",
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

			validateHandler(tt.args.handler)
		})
	}
}

func Test_serveMux_mountMiddlewares(t *testing.T) {
	type args struct {
		cors        *cors
		middlewares []func(http.Handler) http.Handler
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty without cors",
			want: "DATA",
		},
		{
			name: "empty with cors",
			args: args{
				cors: &cors{},
			},
			want: "DATA",
		},
		{
			name: "with middlewares",
			args: args{
				middlewares: []func(http.Handler) http.Handler{
					middlewareWrite("1"),
					middlewareWrite("2"),
					middlewareWrite("3"),
				},
			},
			want: "BEFORE:1 -> BEFORE:2 -> BEFORE:3 -> DATA -> AFTER:3 -> AFTER:2 -> AFTER:1",
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux()
			mux := muxInterface.(*serveMux)

			mux.Use(tt.args.middlewares...)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "/", nil)

			smr := &serveMuxRoute{
				cors: tt.args.cors,
			}
			handler := mux.mountMiddlewares(smr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("DATA"))
			}))

			handler.ServeHTTP(w, r)

			body := w.Result().Body
			defer body.Close()
			bytesBody, err := io.ReadAll(body)
			if err != nil {
				t.Fatal(err)
			}
			stringBody := string(bytesBody)

			if stringBody != tt.want {
				t.Errorf("body response = %v, want %v", stringBody, tt.want)
			}
		})
	}
}

func Test_serveMux_registerHandle(t *testing.T) {
	type args struct {
		path        string
		pattern     string
		handlerKind string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "\\admin",
			args: args{
				path:        "/admin",
				pattern:     "/admin",
				handlerKind: "handlerKind",
			},
			want: "/admin ID:",
		},
		{
			name: "GET \\admin\\user",
			args: args{
				path:        "/admin/user",
				pattern:     "GET /admin/user",
				handlerKind: "handlerKind",
			},
			want: "GET /admin/user ID:",
		},
		{
			name: "GET \\admin\\user\\{id}",
			args: args{
				path:        "/admin/user/1",
				pattern:     "GET /admin/user/{id}",
				handlerKind: "handlerKind",
			},
			want: "GET /admin/user/{id} ID:1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux()
			mux := muxInterface.(*serveMux)

			ts := httptest.NewServer(mux)
			defer ts.Close()

			mux.registerHandle(tt.args.pattern, "handlerKind", handlerId())

			_, _, got := testRequest(t, ts, http.MethodGet, tt.args.path, nil)
			if got != tt.want {
				t.Errorf("ServeMux.Use() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serveMux_registerServeMuxRoute(t *testing.T) {
	type fields struct {
		routes map[string]*serveMuxRoute
		config *serveMuxConfig
	}
	type args struct {
		pattern string
	}
	type want struct {
		serveMuxRoute *serveMuxRoute
		callCreateFn  bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "not exists",
			fields: fields{
				routes: map[string]*serveMuxRoute{},
				config: &serveMuxConfig{},
			},
			args: args{
				pattern: "/",
			},
			want: want{
				callCreateFn: true,
				serveMuxRoute: &serveMuxRoute{
					allowedMethods:       []string{},
					handlerOptionsMaxAge: 0,
					cors:                 nil,
				},
			},
		},
		{
			name: "exists",
			fields: fields{
				routes: map[string]*serveMuxRoute{
					"/": {},
				},
				config: &serveMuxConfig{},
			},
			args: args{
				pattern: "/",
			},
			want: want{
				callCreateFn:  false,
				serveMuxRoute: &serveMuxRoute{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := &serveMux{
				routes: tt.fields.routes,
				config: tt.fields.config,
			}

			callCreateFn := false
			createFn := func(smr *serveMuxRoute) {
				callCreateFn = true
			}
			got := mux.registerServeMuxRoute(tt.args.pattern, createFn)

			if callCreateFn != tt.want.callCreateFn {
				t.Errorf("callCreateFn = %v, want %v", callCreateFn, tt.want.callCreateFn)
			}
			if !reflect.DeepEqual(got, tt.want.serveMuxRoute) {
				t.Errorf("serveMux.registerServeMuxRoute() = %v, want %v", got, tt.want.serveMuxRoute)
			}
		})
	}
}

func Test_serveMux_addRoute(t *testing.T) {
	type args struct {
		method  string
		pattern string
		cors    *cors
	}
	type want struct {
		routes map[string][]string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "GET \\",
			args: args{
				method:  http.MethodGet,
				pattern: "/",
			},
			want: want{
				routes: map[string][]string{
					"/": {http.MethodGet, http.MethodHead, http.MethodOptions},
				},
			},
		},
		{
			name: "GET \\user",
			args: args{
				method:  http.MethodGet,
				pattern: "/user",
			},
			want: want{
				routes: map[string][]string{
					"/":     {http.MethodOptions},
					"/user": {http.MethodGet, http.MethodHead, http.MethodOptions},
				},
			},
		},
		{
			name: "OPTIONS \\user",
			args: args{
				method:  http.MethodOptions,
				pattern: "/user",
			},
			want: want{
				routes: map[string][]string{
					"/":     {http.MethodOptions},
					"/user": {http.MethodOptions},
				},
			},
		},
		{
			name: "GET \\user with cors",
			args: args{
				method:  http.MethodGet,
				pattern: "/user",
				cors:    &cors{},
			},
			want: want{
				routes: map[string][]string{
					"/":     {http.MethodOptions},
					"/user": {http.MethodGet, http.MethodHead, http.MethodOptions},
				},
			},
		},
		{
			name: "OPTIONS \\user with cors",
			args: args{
				method:  http.MethodOptions,
				pattern: "/user",
				cors:    &cors{},
			},
			want: want{
				routes: map[string][]string{
					"/":     {http.MethodOptions},
					"/user": {http.MethodOptions},
				},
			},
		},
		{
			name: "GET host\\",
			args: args{
				method:  http.MethodGet,
				pattern: "host/",
			},
			want: want{
				routes: map[string][]string{
					"host/": {http.MethodGet, http.MethodHead, http.MethodOptions},
				},
			},
		},
		{
			name: "GET host\\user",
			args: args{
				method:  http.MethodGet,
				pattern: "host/user",
			},
			want: want{
				routes: map[string][]string{
					"host/":     {http.MethodOptions},
					"host/user": {http.MethodGet, http.MethodHead, http.MethodOptions},
				},
			},
		},
		{
			name: "OPTIONS host\\user",
			args: args{
				method:  http.MethodOptions,
				pattern: "host/user",
			},
			want: want{
				routes: map[string][]string{
					"host/":     {http.MethodOptions},
					"host/user": {http.MethodOptions},
				},
			},
		},
		{
			name: "GET host\\user with cors",
			args: args{
				method:  http.MethodGet,
				pattern: "host/user",
				cors:    &cors{},
			},
			want: want{
				routes: map[string][]string{
					"host/":     {http.MethodOptions},
					"host/user": {http.MethodGet, http.MethodHead, http.MethodOptions},
				},
			},
		},
		{
			name: "OPTIONS host\\user with cors",
			args: args{
				method:  http.MethodOptions,
				pattern: "host/user",
				cors:    &cors{},
			},
			want: want{
				routes: map[string][]string{
					"host/":     {http.MethodOptions},
					"host/user": {http.MethodOptions},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux()
			mux := muxInterface.(*serveMux)
			mux.config.cors = tt.args.cors

			mux.addRoute(tt.args.method, tt.args.pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

			for k, smr := range mux.routes {
				if want, ok := tt.want.routes[k]; !ok {
					t.Errorf("serveMux.routes[%v] unexpected route", k)
				} else {
					if got := smr.allowedMethods; !reflect.DeepEqual(got, want) {
						t.Errorf("smr.allowedMethods = %v, want %v", got, want)
					}
				}
			}
		})
	}
}

func Test_serveMux_Connect(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "\\user",
			args: args{
				pattern: "/user",
			},
		},
		{
			name: "\\user\\{id}",
			args: args{
				pattern: "/user/{id}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux()
			mux := muxInterface.(*serveMux)

			mux.Connect(tt.args.pattern, func(w http.ResponseWriter, r *http.Request) {})

			if smr, ok := mux.routes[tt.args.pattern]; !ok {
				t.Errorf("serveMux.routes[%v] not found", tt.args.pattern)
			} else {
				if !slices.Contains(smr.allowedMethods, http.MethodConnect) {
					t.Errorf("serveMux.routes['/'].allowedMethods %v not found", smr.allowedMethods)
				}
			}
		})
	}
}

func Test_serveMux_Delete(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "\\user",
			args: args{
				pattern: "/user",
			},
		},
		{
			name: "\\user\\{id}",
			args: args{
				pattern: "/user/{id}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux()
			mux := muxInterface.(*serveMux)

			mux.Delete(tt.args.pattern, func(w http.ResponseWriter, r *http.Request) {})

			if smr, ok := mux.routes[tt.args.pattern]; !ok {
				t.Errorf("serveMux.routes[%v] not found", tt.args.pattern)
			} else {
				if !slices.Contains(smr.allowedMethods, http.MethodDelete) {
					t.Errorf("serveMux.routes['/'].allowedMethods %v not found", smr.allowedMethods)
				}
			}
		})
	}
}

func Test_serveMux_Get(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "\\user",
			args: args{
				pattern: "/user",
			},
		},
		{
			name: "\\user\\{id}",
			args: args{
				pattern: "/user/{id}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux()
			mux := muxInterface.(*serveMux)

			mux.Get(tt.args.pattern, func(w http.ResponseWriter, r *http.Request) {})

			if smr, ok := mux.routes[tt.args.pattern]; !ok {
				t.Errorf("serveMux.routes[%v] not found", tt.args.pattern)
			} else {
				if !slices.Contains(smr.allowedMethods, http.MethodGet) {
					t.Errorf("serveMux.routes['/'].allowedMethods %v not found", smr.allowedMethods)
				}
			}
		})
	}
}

func Test_serveMux_Head(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "\\user",
			args: args{
				pattern: "/user",
			},
		},
		{
			name: "\\user\\{id}",
			args: args{
				pattern: "/user/{id}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux()
			mux := muxInterface.(*serveMux)

			mux.Head(tt.args.pattern, func(w http.ResponseWriter, r *http.Request) {})

			if smr, ok := mux.routes[tt.args.pattern]; !ok {
				t.Errorf("serveMux.routes[%v] not found", tt.args.pattern)
			} else {
				if !slices.Contains(smr.allowedMethods, http.MethodHead) {
					t.Errorf("serveMux.routes['/'].allowedMethods %v not found", smr.allowedMethods)
				}
			}
		})
	}
}

func Test_serveMux_Options(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "\\user",
			args: args{
				pattern: "/user",
			},
		},
		{
			name: "\\user\\{id}",
			args: args{
				pattern: "/user/{id}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux()
			mux := muxInterface.(*serveMux)

			mux.Options(tt.args.pattern, func(w http.ResponseWriter, r *http.Request) {})

			if smr, ok := mux.routes[tt.args.pattern]; !ok {
				t.Errorf("serveMux.routes[%v] not found", tt.args.pattern)
			} else {
				if !slices.Contains(smr.allowedMethods, http.MethodOptions) {
					t.Errorf("serveMux.routes['/'].allowedMethods %v not found", smr.allowedMethods)
				}
			}
		})
	}
}

func Test_serveMux_Patch(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "\\user",
			args: args{
				pattern: "/user",
			},
		},
		{
			name: "\\user\\{id}",
			args: args{
				pattern: "/user/{id}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux()
			mux := muxInterface.(*serveMux)

			mux.Patch(tt.args.pattern, func(w http.ResponseWriter, r *http.Request) {})

			if smr, ok := mux.routes[tt.args.pattern]; !ok {
				t.Errorf("serveMux.routes[%v] not found", tt.args.pattern)
			} else {
				if !slices.Contains(smr.allowedMethods, http.MethodPatch) {
					t.Errorf("serveMux.routes['/'].allowedMethods %v not found", smr.allowedMethods)
				}
			}
		})
	}
}

func Test_serveMux_Post(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "\\user",
			args: args{
				pattern: "/user",
			},
		},
		{
			name: "\\user\\{id}",
			args: args{
				pattern: "/user/{id}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux()
			mux := muxInterface.(*serveMux)

			mux.Post(tt.args.pattern, func(w http.ResponseWriter, r *http.Request) {})

			if smr, ok := mux.routes[tt.args.pattern]; !ok {
				t.Errorf("serveMux.routes[%v] not found", tt.args.pattern)
			} else {
				if !slices.Contains(smr.allowedMethods, http.MethodPost) {
					t.Errorf("serveMux.routes['/'].allowedMethods %v not found", smr.allowedMethods)
				}
			}
		})
	}
}

func Test_serveMux_Put(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "\\user",
			args: args{
				pattern: "/user",
			},
		},
		{
			name: "\\user\\{id}",
			args: args{
				pattern: "/user/{id}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux()
			mux := muxInterface.(*serveMux)

			mux.Put(tt.args.pattern, func(w http.ResponseWriter, r *http.Request) {})

			if smr, ok := mux.routes[tt.args.pattern]; !ok {
				t.Errorf("serveMux.routes[%v] not found", tt.args.pattern)
			} else {
				if !slices.Contains(smr.allowedMethods, http.MethodPut) {
					t.Errorf("serveMux.routes['/'].allowedMethods %v not found", smr.allowedMethods)
				}
			}
		})
	}
}

func Test_serveMux_Trace(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "\\user",
			args: args{
				pattern: "/user",
			},
		},
		{
			name: "\\user\\{id}",
			args: args{
				pattern: "/user/{id}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux()
			mux := muxInterface.(*serveMux)

			mux.Trace(tt.args.pattern, func(w http.ResponseWriter, r *http.Request) {})

			if smr, ok := mux.routes[tt.args.pattern]; !ok {
				t.Errorf("serveMux.routes[%v] not found", tt.args.pattern)
			} else {
				if !slices.Contains(smr.allowedMethods, http.MethodTrace) {
					t.Errorf("serveMux.routes['/'].allowedMethods %v not found", smr.allowedMethods)
				}
			}
		})
	}
}

func Test_serveMux_Method(t *testing.T) {
	type args struct {
		method  string
		pattern string
	}
	type want struct {
		panicked   bool
		panicError string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "\\user",
			args: args{
				method:  "CUSTOM",
				pattern: "/user",
			},
		},
		{
			name: "\\user\\{id}",
			args: args{
				method:  "CUSTOM",
				pattern: "/user/{id}",
			},
		},
		{
			name: "method not specified",
			args: args{
				method: "",
			},
			want: want{
				panicked:   true,
				panicError: "method not specified",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			muxInterface := NewServeMux()
			mux := muxInterface.(*serveMux)

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

			mux.Method(tt.args.method, tt.args.pattern, func(w http.ResponseWriter, r *http.Request) {})

			if smr, ok := mux.routes[tt.args.pattern]; !ok {
				t.Errorf("serveMux.routes[%v] not found", tt.args.pattern)
			} else {
				if !slices.Contains(smr.allowedMethods, tt.args.method) {
					t.Errorf("serveMux.routes['/'].allowedMethods %v not found", smr.allowedMethods)
				}
			}
		})
	}
}

func Test_serveMux_ServeHTTP(t *testing.T) {
	mux := NewServeMux()

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	mux.ServeHTTP(w, r)

	if statusCode := w.Result().StatusCode; statusCode != 404 {
		t.Errorf("statusCode = %v, want %v", statusCode, 404)
	}
}

func BenchmarkServerMux_Group(b *testing.B) {
	mux := NewServeMux()

	muxGroup1 := mux.Group("/path1")
	muxGroup1.Get("/", func(w http.ResponseWriter, r *http.Request) {})

	muxGroup2 := mux.Group("/path2")
	muxGroup2.Get("/", func(w http.ResponseWriter, r *http.Request) {})

	muxGroup3 := mux.Group("/path3")
	muxGroup3.Get("/route1", func(w http.ResponseWriter, r *http.Request) {})
	muxGroup3.Get("/route2", func(w http.ResponseWriter, r *http.Request) {})

	muxGroup4 := mux.Group("/path4")
	muxGroup4.Get("/route1", func(w http.ResponseWriter, r *http.Request) {})
	muxGroup4.Get("/route2", func(w http.ResponseWriter, r *http.Request) {})
	muxGroup4.Get("/route3", func(w http.ResponseWriter, r *http.Request) {})
	muxGroup4.Get("/route4", func(w http.ResponseWriter, r *http.Request) {})

	routes := []string{
		"/path1",
		"/path2",
		"/path3/route1",
		"/path3/route2",
		"/path4/route1",
		"/path4/route2",
		"/path4/route3",
		"/path4/route4",
	}

	for _, path := range routes {
		b.Run("route:"+path, func(b *testing.B) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, path, nil)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				mux.ServeHTTP(w, r)
			}
		})
	}
}

func BenchmarkServerMux_Route(b *testing.B) {
	mux := NewServeMux()

	mux.Route("/path1", func(muxRoute1 Router) {
		muxRoute1.Get("/", func(w http.ResponseWriter, r *http.Request) {})
	})

	mux.Route("/path2", func(muxRoute2 Router) {
		muxRoute2.Get("/", func(w http.ResponseWriter, r *http.Request) {})
	})

	mux.Route("/path3", func(muxRoute3 Router) {
		muxRoute3.Get("/route1", func(w http.ResponseWriter, r *http.Request) {})
		muxRoute3.Get("/route2", func(w http.ResponseWriter, r *http.Request) {})
	})

	mux.Route("/path4", func(muxRoute4 Router) {
		muxRoute4.Get("/route1", func(w http.ResponseWriter, r *http.Request) {})
		muxRoute4.Get("/route2", func(w http.ResponseWriter, r *http.Request) {})
		muxRoute4.Get("/route3", func(w http.ResponseWriter, r *http.Request) {})
		muxRoute4.Get("/route4", func(w http.ResponseWriter, r *http.Request) {})
	})

	routes := []string{
		"/path1",
		"/path2",
		"/path3/route1",
		"/path3/route2",
		"/path4/route1",
		"/path4/route2",
		"/path4/route3",
		"/path4/route4",
	}

	for _, path := range routes {
		b.Run("route:"+path, func(b *testing.B) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, path, nil)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				mux.ServeHTTP(w, r)
			}
		})
	}
}
