package api

import (
	"net/http"
	"net/http/httptest"
	"task-service/config"
	"testing"

	"github.com/sunshineOfficial/golib/golog"
)

func TestTaskRoutesAllowUnauthenticatedRequests(t *testing.T) {
	builder := NewServerBuilder(t.Context(), golog.NewLogger("test"), config.Settings{
		Port: 80,
	})
	builder.AddTasks(nil)

	routes := []struct {
		method string
		path   string
	}{
		{method: http.MethodPost, path: "/tasks"},
		{method: http.MethodGet, path: "/tasks/1/extended"},
		{method: http.MethodGet, path: "/tasks/1"},
		{method: http.MethodPatch, path: "/tasks/1"},
		{method: http.MethodDelete, path: "/tasks/1"},
		{method: http.MethodGet, path: "/tasks/brigade/1/extended"},
		{method: http.MethodGet, path: "/tasks/brigade/1"},
		{method: http.MethodGet, path: "/tasks"},
		{method: http.MethodPatch, path: "/tasks/1/start"},
		{method: http.MethodPost, path: "/tasks/assign"},
	}

	for _, route := range routes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(route.method, route.path, nil)

			builder.router.ServeHTTP(response, request)

			if response.Code == http.StatusUnauthorized {
				t.Fatalf("status = %d, route must be open without authorization", response.Code)
			}
		})
	}
}
