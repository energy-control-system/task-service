package handler

import (
	"fmt"
	"net/http"
	"task-service/service/task"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

func AddTask(s *task.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var request task.AddRequest
		if err := c.ReadJson(&request); err != nil {
			return fmt.Errorf("failed to read request: %w", err)
		}

		t, err := s.Add(c.Ctx(), c.Log().WithTags("addTask"), request)
		if err != nil {
			return fmt.Errorf("failed to add task: %w", err)
		}

		return c.WriteJson(http.StatusOK, t)
	}
}

type taskIDVars struct {
	ID int `path:"id"`
}

func GetTaskByID(s *task.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars taskIDVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read task id: %w", err)
		}

		t, err := s.GetByID(c.Ctx(), vars.ID)
		if err != nil {
			return fmt.Errorf("failed to get task by id: %w", err)
		}

		return c.WriteJson(http.StatusOK, t)
	}
}

type brigadeIDVars struct {
	BrigadeID int `path:"brigadeID"`
}

func GetTasksByBrigade(s *task.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars brigadeIDVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read brigade id: %w", err)
		}

		response, err := s.GetByBrigade(c.Ctx(), vars.BrigadeID)
		if err != nil {
			return fmt.Errorf("failed to get tasks by brigade id: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

func StartTask(s *task.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars taskIDVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read task id: %w", err)
		}

		t, err := s.StartTask(c.Ctx(), c.Log().WithTags("startTask"), vars.ID)
		if err != nil {
			return fmt.Errorf("failed to start task: %w", err)
		}

		return c.WriteJson(http.StatusOK, t)
	}
}
