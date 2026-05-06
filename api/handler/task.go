package handler

import (
	"fmt"
	"net/http"
	"task-service/service/task"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
	"github.com/sunshineOfficial/golib/pagination"
)

// AddTask godoc
// @Summary Create task
// @Description Creates an inspection task for a metering object.
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body task.AddRequest true "Task creation payload"
// @Success 200 {object} task.Task
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /tasks [post]
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

// GetTaskByID godoc
// @Summary Get task by ID
// @Description Returns an inspection task by identifier.
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} task.Task
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 404 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /tasks/{id} [get]
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

// GetTasksByBrigade godoc
// @Summary List tasks by brigade
// @Description Returns tasks assigned to a brigade.
// @Tags tasks
// @Produce json
// @Param brigadeID path int true "Brigade ID"
// @Param limit query int false "Maximum number of items to return; 0 means no limit"
// @Param offset query int false "Number of items to skip"
// @Success 200 {array} task.Task
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /tasks/brigade/{brigadeID} [get]
func GetTasksByBrigade(s *task.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars brigadeIDVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read brigade id: %w", err)
		}

		var pageVars pagination.Pagination
		if err := c.Vars(&pageVars); err != nil {
			return fmt.Errorf("failed to read pagination: %w", err)
		}

		response, err := s.GetByBrigade(c.Ctx(), vars.BrigadeID, pageVars)
		if err != nil {
			return fmt.Errorf("failed to get tasks by brigade id: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

// GetAllTasks godoc
// @Summary List tasks
// @Description Returns all inspection tasks.
// @Tags tasks
// @Produce json
// @Param limit query int false "Maximum number of items to return; 0 means no limit"
// @Param offset query int false "Number of items to skip"
// @Success 200 {array} task.Task
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /tasks [get]
func GetAllTasks(s *task.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars pagination.Pagination
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read pagination: %w", err)
		}

		response, err := s.GetAll(c.Ctx(), vars)
		if err != nil {
			return fmt.Errorf("failed to get all tasks: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

// StartTask godoc
// @Summary Start task
// @Description Marks a planned task as started.
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} task.Task
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 404 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /tasks/{id}/start [patch]
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

// AssignTaskToBrigade godoc
// @Summary Assign task to brigade
// @Description Assigns an inspection task to an inspector brigade.
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body task.AssignToBrigadeRequest true "Assignment payload"
// @Success 200 {object} task.Task
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 404 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /tasks/assign [patch]
func AssignTaskToBrigade(s *task.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var request task.AssignToBrigadeRequest
		if err := c.ReadJson(&request); err != nil {
			return fmt.Errorf("failed to read request: %w", err)
		}

		t, err := s.AssignToBrigade(c.Ctx(), c.Log().WithTags("assignTaskToBrigade"), request)
		if err != nil {
			return fmt.Errorf("failed to assign task to brigade: %w", err)
		}

		return c.WriteJson(http.StatusOK, t)
	}
}
