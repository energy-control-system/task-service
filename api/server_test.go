package api

import (
	"net/http"
	"net/http/httptest"
	"task-service/config"
	"testing"

	"github.com/sunshineOfficial/golib/golog"
)

func TestTaskAuthorizationPolicy(t *testing.T) {
	builder := NewServerBuilder(t.Context(), golog.NewLogger("test"), config.Settings{
		Port: 80,
	})
	builder.AddTasks(nil)

	t.Run("task creation requires authorization", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/tasks", nil)

		builder.router.ServeHTTP(response, request)

		if response.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
		}
	})

	t.Run("get by id allows internal calls without authorization", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)

		builder.router.ServeHTTP(response, request)

		if response.Code == http.StatusUnauthorized {
			t.Fatalf("status = %d, route must stay open for internal service calls", response.Code)
		}
	})
}
