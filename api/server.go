package api

import (
	"context"
	"fmt"
	"task-service/api/handler"
	"task-service/config"
	"task-service/service/task"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
	"github.com/sunshineOfficial/golib/gohttp/gorouter/middleware"
	"github.com/sunshineOfficial/golib/gohttp/gorouter/plugin"
	"github.com/sunshineOfficial/golib/gohttp/goserver"
	"github.com/sunshineOfficial/golib/golog"
)

type ServerBuilder struct {
	server goserver.Server
	router *gorouter.Router
}

func NewServerBuilder(ctx context.Context, log golog.Logger, settings config.Settings) *ServerBuilder {
	return &ServerBuilder{
		server: goserver.NewHTTPServer(ctx, log, fmt.Sprintf(":%d", settings.Port)),
		router: gorouter.NewRouter(log).Use(
			middleware.Metrics(),
			middleware.Recover,
			middleware.LogError,
		),
	}
}

func (s *ServerBuilder) AddDebug() {
	s.router.Install(plugin.NewPProf(), plugin.NewMetrics())
}

func (s *ServerBuilder) AddTasks(service *task.Service) {
	r := s.router.SubRouter("/tasks")
	r.HandlePost("", handler.AddTask(service))
	r.HandleGet("/{id}", handler.GetTaskByID(service))
	r.HandleGet("/brigade/{brigadeID}", handler.GetTasksByBrigade(service))
	r.HandleGet("", handler.GetAllTasks(service))
	r.HandlePatch("/{id}/start", handler.StartTask(service))
}

func (s *ServerBuilder) Build() goserver.Server {
	s.server.UseHandler(s.router)

	return s.server
}
