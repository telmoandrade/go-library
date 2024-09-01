package httpserver_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/telmoandrade/go-library/httpserver"
)

func handlerId(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		_, err := w.Write([]byte(fmt.Sprintf("%s ID:%s", r.Pattern, id)))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func middlewareEmpty(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func TestServeMux_Use(t *testing.T) {
	tests := []struct {
		name string
		args int
		want string
	}{
		{
			name: "get user id 1",
			args: 1,
			want: "GET /user/{id} ID:1",
		},
		{
			name: "get user id 2",
			args: 2,
			want: "GET /user/{id} ID:2",
		},
	}
	for _, tt := range tests {
		mux := httpserver.NewServeMux()
		mux.Use(middlewareEmpty)
		mux.Get("/user/{id}", handlerId(t))

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run(tt.name, func(t *testing.T) {
			_, got := testRequest(t, ts, "GET", fmt.Sprintf("/user/%d", tt.args), nil)
			if got != tt.want {
				t.Errorf("ServeMux.Use() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServeMux_With(t *testing.T) {
	tests := []struct {
		name string
		args int
		want string
	}{
		{
			name: "get user id 1",
			args: 1,
			want: "GET /user/{id} ID:1",
		},
		{
			name: "get user id 2",
			args: 2,
			want: "GET /user/{id} ID:2",
		},
	}
	for _, tt := range tests {
		mux := httpserver.NewServeMux()
		mux.With(middlewareEmpty).Get("/user/{id}", handlerId(t))

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run(tt.name, func(t *testing.T) {
			_, got := testRequest(t, ts, "GET", fmt.Sprintf("/user/%d", tt.args), nil)
			if got != tt.want {
				t.Errorf("ServeMux.Use() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServeMux_Group(t *testing.T) {
	tests := []struct {
		name string
		args int
		want string
	}{
		{
			name: "get user id 1",
			args: 1,
			want: "GET /user/{id} ID:1",
		},
		{
			name: "get user id 2",
			args: 2,
			want: "GET /user/{id} ID:2",
		},
	}
	for _, tt := range tests {
		mux := httpserver.NewServeMux()
		muxUser := mux.Group("/user")

		muxUser.Get("/{id}", handlerId(t))

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run(tt.name, func(t *testing.T) {
			_, got := testRequest(t, ts, "GET", fmt.Sprintf("/user/%d", tt.args), nil)
			if got != tt.want {
				t.Errorf("ServeMux.Group() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServeMux_Route(t *testing.T) {
	tests := []struct {
		name string
		args int
		want string
	}{
		{
			name: "get user id 1",
			args: 1,
			want: "GET /user/{id} ID:1",
		},
		{
			name: "get user id 2",
			args: 2,
			want: "GET /user/{id} ID:2",
		},
	}
	for _, tt := range tests {
		mux := httpserver.NewServeMux()
		mux.Route("/user", func(muxUser httpserver.Router) {
			muxUser.Get("/{id}", handlerId(t))
		})

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run(tt.name, func(t *testing.T) {
			_, got := testRequest(t, ts, "GET", fmt.Sprintf("/user/%d", tt.args), nil)
			if got != tt.want {
				t.Errorf("ServeMux.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServeMux_Mount(t *testing.T) {
	tests := []struct {
		name string
		args int
		want string
	}{
		{
			name: "get user id 1",
			args: 1,
			want: "GET /{id} ID:1",
		},
		{
			name: "get user id 2",
			args: 2,
			want: "GET /{id} ID:2",
		},
	}
	for _, tt := range tests {
		muxUser := httpserver.NewServeMux()
		muxUser.Get("/{id}", handlerId(t))

		mux := httpserver.NewServeMux()
		mux.Mount("/user", muxUser)

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run(tt.name, func(t *testing.T) {
			_, got := testRequest(t, ts, "GET", fmt.Sprintf("/user/%d", tt.args), nil)
			if got != tt.want {
				t.Errorf("ServeMux.Mount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServeMux_Connect(t *testing.T) {
	type want struct {
		body       string
		statusCode int
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "user id 1",
			args: "/user/1",
			want: want{
				body:       "CONNECT /user/{id} ID:1",
				statusCode: 200,
			},
		},
		{
			name: "user id 2",
			args: "/user/2",
			want: want{
				body:       "CONNECT /user/{id} ID:2",
				statusCode: 200,
			},
		}, {
			name: "client",
			args: "/client",
			want: want{
				body:       "404 page not found\n",
				statusCode: 404,
			},
		},
	}
	for _, tt := range tests {
		mux := httpserver.NewServeMux()
		mux.Connect("/user/{id}", handlerId(t))

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequest(t, ts, "CONNECT", tt.args, nil)
			if body != tt.want.body {
				t.Errorf("body = %v, want %v", body, tt.want.body)
			}
			if statusCode != tt.want.statusCode {
				t.Errorf("statusCode = %v, want %v", statusCode, tt.want.statusCode)
			}
		})
	}
}

func TestServeMux_Delete(t *testing.T) {
	type want struct {
		body       string
		statusCode int
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "user id 1",
			args: "/user/1",
			want: want{
				body:       "DELETE /user/{id} ID:1",
				statusCode: 200,
			},
		},
		{
			name: "user id 2",
			args: "/user/2",
			want: want{
				body:       "DELETE /user/{id} ID:2",
				statusCode: 200,
			},
		}, {
			name: "client",
			args: "/client",
			want: want{
				body:       "404 page not found\n",
				statusCode: 404,
			},
		},
	}
	for _, tt := range tests {
		mux := httpserver.NewServeMux()
		mux.Delete("/user/{id}", handlerId(t))

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequest(t, ts, "DELETE", tt.args, nil)
			if body != tt.want.body {
				t.Errorf("body = %v, want %v", body, tt.want.body)
			}
			if statusCode != tt.want.statusCode {
				t.Errorf("statusCode = %v, want %v", statusCode, tt.want.statusCode)
			}
		})
	}
}

func TestServeMux_Get(t *testing.T) {
	type want struct {
		body       string
		statusCode int
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "user id 1",
			args: "/user/1",
			want: want{
				body:       "GET /user/{id} ID:1",
				statusCode: 200,
			},
		},
		{
			name: "user id 2",
			args: "/user/2",
			want: want{
				body:       "GET /user/{id} ID:2",
				statusCode: 200,
			},
		}, {
			name: "client",
			args: "/client",
			want: want{
				body:       "404 page not found\n",
				statusCode: 404,
			},
		},
	}
	for _, tt := range tests {
		mux := httpserver.NewServeMux()
		mux.Get("/user/{id}", handlerId(t))

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequest(t, ts, "GET", tt.args, nil)
			if body != tt.want.body {
				t.Errorf("body = %v, want %v", body, tt.want.body)
			}
			if statusCode != tt.want.statusCode {
				t.Errorf("statusCode = %v, want %v", statusCode, tt.want.statusCode)
			}
		})
	}
}

func TestServeMux_Head(t *testing.T) {
	type want struct {
		body       string
		statusCode int
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "user id 1",
			args: "/user/1",
			want: want{
				body:       "",
				statusCode: 200,
			},
		},
		{
			name: "user id 2",
			args: "/user/2",
			want: want{
				body:       "",
				statusCode: 200,
			},
		}, {
			name: "client",
			args: "/client",
			want: want{
				body:       "",
				statusCode: 404,
			},
		},
	}
	for _, tt := range tests {
		mux := httpserver.NewServeMux()
		mux.Head("/user/{id}", handlerId(t))

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequest(t, ts, "HEAD", tt.args, nil)
			if body != tt.want.body {
				t.Errorf("body = %v, want %v", body, tt.want.body)
			}
			if statusCode != tt.want.statusCode {
				t.Errorf("statusCode = %v, want %v", statusCode, tt.want.statusCode)
			}
		})
	}
}

func TestServeMux_Options(t *testing.T) {
	type want struct {
		body       string
		statusCode int
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "user id 1",
			args: "/user/1",
			want: want{
				body:       "OPTIONS /user/{id} ID:1",
				statusCode: 200,
			},
		},
		{
			name: "user id 2",
			args: "/user/2",
			want: want{
				body:       "OPTIONS /user/{id} ID:2",
				statusCode: 200,
			},
		}, {
			name: "client",
			args: "/client",
			want: want{
				body:       "404 page not found\n",
				statusCode: 404,
			},
		},
	}
	for _, tt := range tests {
		mux := httpserver.NewServeMux()
		mux.Options("/user/{id}", handlerId(t))

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequest(t, ts, "OPTIONS", tt.args, nil)
			if body != tt.want.body {
				t.Errorf("body = %v, want %v", body, tt.want.body)
			}
			if statusCode != tt.want.statusCode {
				t.Errorf("statusCode = %v, want %v", statusCode, tt.want.statusCode)
			}
		})
	}
}

func TestServeMux_Patch(t *testing.T) {
	type want struct {
		body       string
		statusCode int
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "user id 1",
			args: "/user/1",
			want: want{
				body:       "PATCH /user/{id} ID:1",
				statusCode: 200,
			},
		},
		{
			name: "user id 2",
			args: "/user/2",
			want: want{
				body:       "PATCH /user/{id} ID:2",
				statusCode: 200,
			},
		}, {
			name: "client",
			args: "/client",
			want: want{
				body:       "404 page not found\n",
				statusCode: 404,
			},
		},
	}
	for _, tt := range tests {
		mux := httpserver.NewServeMux()
		mux.Patch("/user/{id}", handlerId(t))

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequest(t, ts, "PATCH", tt.args, nil)
			if body != tt.want.body {
				t.Errorf("body = %v, want %v", body, tt.want.body)
			}
			if statusCode != tt.want.statusCode {
				t.Errorf("statusCode = %v, want %v", statusCode, tt.want.statusCode)
			}
		})
	}
}

func TestServeMux_Post(t *testing.T) {
	type want struct {
		body       string
		statusCode int
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "user id 1",
			args: "/user/1",
			want: want{
				body:       "POST /user/{id} ID:1",
				statusCode: 200,
			},
		},
		{
			name: "user id 2",
			args: "/user/2",
			want: want{
				body:       "POST /user/{id} ID:2",
				statusCode: 200,
			},
		}, {
			name: "client",
			args: "/client",
			want: want{
				body:       "404 page not found\n",
				statusCode: 404,
			},
		},
	}
	for _, tt := range tests {
		mux := httpserver.NewServeMux()
		mux.Post("/user/{id}", handlerId(t))

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequest(t, ts, "POST", tt.args, nil)
			if body != tt.want.body {
				t.Errorf("body = %v, want %v", body, tt.want.body)
			}
			if statusCode != tt.want.statusCode {
				t.Errorf("statusCode = %v, want %v", statusCode, tt.want.statusCode)
			}
		})
	}
}

func TestServeMux_Put(t *testing.T) {
	type want struct {
		body       string
		statusCode int
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "user id 1",
			args: "/user/1",
			want: want{
				body:       "PUT /user/{id} ID:1",
				statusCode: 200,
			},
		},
		{
			name: "user id 2",
			args: "/user/2",
			want: want{
				body:       "PUT /user/{id} ID:2",
				statusCode: 200,
			},
		}, {
			name: "client",
			args: "/client",
			want: want{
				body:       "404 page not found\n",
				statusCode: 404,
			},
		},
	}
	for _, tt := range tests {
		mux := httpserver.NewServeMux()
		mux.Put("/user/{id}", handlerId(t))

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequest(t, ts, "PUT", tt.args, nil)
			if body != tt.want.body {
				t.Errorf("body = %v, want %v", body, tt.want.body)
			}
			if statusCode != tt.want.statusCode {
				t.Errorf("statusCode = %v, want %v", statusCode, tt.want.statusCode)
			}
		})
	}
}

func TestServeMux_Trace(t *testing.T) {
	type want struct {
		body       string
		statusCode int
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "user id 1",
			args: "/user/1",
			want: want{
				body:       "TRACE /user/{id} ID:1",
				statusCode: 200,
			},
		},
		{
			name: "user id 2",
			args: "/user/2",
			want: want{
				body:       "TRACE /user/{id} ID:2",
				statusCode: 200,
			},
		}, {
			name: "client",
			args: "/client",
			want: want{
				body:       "404 page not found\n",
				statusCode: 404,
			},
		},
	}
	for _, tt := range tests {
		mux := httpserver.NewServeMux()
		mux.Trace("/user/{id}", handlerId(t))

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequest(t, ts, "TRACE", tt.args, nil)
			if body != tt.want.body {
				t.Errorf("body = %v, want %v", body, tt.want.body)
			}
			if statusCode != tt.want.statusCode {
				t.Errorf("statusCode = %v, want %v", statusCode, tt.want.statusCode)
			}
		})
	}
}

func TestServeMux_Method(t *testing.T) {
	type want struct {
		body       string
		statusCode int
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "user id 1",
			args: "/user/1",
			want: want{
				body:       "CUSTOM /user/{id} ID:1",
				statusCode: 200,
			},
		},
		{
			name: "user id 2",
			args: "/user/2",
			want: want{
				body:       "CUSTOM /user/{id} ID:2",
				statusCode: 200,
			},
		}, {
			name: "client",
			args: "/client",
			want: want{
				body:       "404 page not found\n",
				statusCode: 404,
			},
		},
	}
	for _, tt := range tests {
		mux := httpserver.NewServeMux()
		mux.Method("CUSTOM", "/user/{id}", handlerId(t))

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequest(t, ts, "CUSTOM", tt.args, nil)
			if body != tt.want.body {
				t.Errorf("body = %v, want %v", body, tt.want.body)
			}
			if statusCode != tt.want.statusCode {
				t.Errorf("statusCode = %v, want %v", statusCode, tt.want.statusCode)
			}
		})
	}
}

func BenchmarkServerMux_Group(b *testing.B) {
	mux := httpserver.NewServeMux()

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
			r, _ := http.NewRequest("GET", path, nil)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				mux.ServeHTTP(w, r)
			}
		})
	}
}

func BenchmarkServerMux_Route(b *testing.B) {
	mux := httpserver.NewServeMux()

	mux.Route("/path1", func(muxRoute1 httpserver.Router) {
		muxRoute1.Get("/", func(w http.ResponseWriter, r *http.Request) {})
	})

	mux.Route("/path2", func(muxRoute2 httpserver.Router) {
		muxRoute2.Get("/", func(w http.ResponseWriter, r *http.Request) {})
	})

	mux.Route("/path3", func(muxRoute3 httpserver.Router) {
		muxRoute3.Get("/route1", func(w http.ResponseWriter, r *http.Request) {})
		muxRoute3.Get("/route2", func(w http.ResponseWriter, r *http.Request) {})
	})

	mux.Route("/path4", func(muxRoute4 httpserver.Router) {
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
			r, _ := http.NewRequest("GET", path, nil)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				mux.ServeHTTP(w, r)
			}
		})
	}
}

func BenchmarkServerMux_Mount(b *testing.B) {
	mux := httpserver.NewServeMux()

	muxMount1 := httpserver.NewServeMux()
	muxMount1.Get("/", func(w http.ResponseWriter, r *http.Request) {})

	muxMount2 := httpserver.NewServeMux()
	muxMount2.Get("/", func(w http.ResponseWriter, r *http.Request) {})

	muxMount3 := httpserver.NewServeMux()
	muxMount3.Get("/route1", func(w http.ResponseWriter, r *http.Request) {})
	muxMount3.Get("/route2", func(w http.ResponseWriter, r *http.Request) {})

	muxMount4 := httpserver.NewServeMux()
	muxMount4.Get("/route1", func(w http.ResponseWriter, r *http.Request) {})
	muxMount4.Get("/route2", func(w http.ResponseWriter, r *http.Request) {})
	muxMount4.Get("/route3", func(w http.ResponseWriter, r *http.Request) {})
	muxMount4.Get("/route4", func(w http.ResponseWriter, r *http.Request) {})

	mux.Mount("/path1", muxMount1)
	mux.Mount("/path2", muxMount2)
	mux.Mount("/path3", muxMount3)
	mux.Mount("/path4", muxMount4)

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
			r, _ := http.NewRequest("GET", path, nil)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				mux.ServeHTTP(w, r)
			}
		})
	}
}
