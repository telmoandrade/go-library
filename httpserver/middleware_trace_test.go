package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/telmoandrade/go-library/logger"
)

func TestMiddlewareTrace(t *testing.T) {
	tests := []struct {
		name   string
		status int
		args   uuid.UUID
		want   int
	}{
		{
			name:   "without logId",
			status: http.StatusOK,
			args:   uuid.Nil,
			want:   http.StatusOK,
		},
		{
			name:   "with logId",
			status: http.StatusOK,
			args:   uuid.Must(uuid.NewV7()),
			want:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			r.Pattern = "GET /"
			w := httptest.NewRecorder()

			ctx := context.WithValue(r.Context(), logger.ContextLogID, tt.args)

			m := MiddlewareTrace(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
			}))

			m.ServeHTTP(w, r.WithContext(ctx))

			if w.Code != tt.want {
				t.Errorf("Code() = %v, want %v", w.Code, tt.want)
			}
		})
	}
}
