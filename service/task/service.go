package task

import (
	"context"
	"encoding/json"
	"fmt"
	"task-service/cluster/inspection"
	"task-service/database/task"
	"time"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/gokafka"
	"github.com/sunshineOfficial/golib/golog"
)

const kafkaSubscribeTimeout = 2 * time.Minute

type Service struct {
	taskRepository *task.Postgres
	taskPublisher  *Publisher
}

func NewService(taskRepository *task.Postgres, taskPublisher *Publisher) *Service {
	return &Service{
		taskRepository: taskRepository,
		taskPublisher:  taskPublisher,
	}
}

func (s *Service) Add(ctx goctx.Context, log golog.Logger, request AddRequest) (Task, error) {
	dbTask, err := s.taskRepository.Add(ctx, task.Task{
		BrigadeID:   request.BrigadeID,
		ObjectID:    request.ObjectID,
		PlanVisitAt: request.PlanVisitAt,
		Status:      int(StatusPlanned),
		Comment:     request.Comment,
	})
	if err != nil {
		return Task{}, fmt.Errorf("add task to db: %w", err)
	}

	t := MapFromDB(dbTask)

	go s.taskPublisher.Publish(ctx, log, EventTypeCreate, t)

	return t, nil
}

func (s *Service) GetByID(ctx goctx.Context, id int) (Task, error) {
	dbTask, err := s.taskRepository.GetByID(ctx, id)
	if err != nil {
		return Task{}, fmt.Errorf("get task %d from db: %w", id, err)
	}

	return MapFromDB(dbTask), nil
}

func (s *Service) GetByBrigade(ctx goctx.Context, brigadeID int) ([]Task, error) {
	dbTasks, err := s.taskRepository.GetByBrigade(ctx, brigadeID)
	if err != nil {
		return nil, fmt.Errorf("get tasks by brigade id %d from db: %w", brigadeID, err)
	}

	return MapSliceFromDB(dbTasks), nil
}

func (s *Service) StartTask(ctx goctx.Context, log golog.Logger, id int) (Task, error) {
	dbTask, err := s.taskRepository.StartTask(ctx, id)
	if err != nil {
		return Task{}, fmt.Errorf("start task %d: %w", id, err)
	}

	t := MapFromDB(dbTask)

	go s.taskPublisher.Publish(ctx, log, EventTypeUpdate, t)

	return t, nil
}

func (s *Service) SubscriberOnInspectionEvent(mainCtx context.Context, log golog.Logger) gokafka.Subscriber {
	return func(message gokafka.Message, err error) {
		ctx, cancel := context.WithTimeout(mainCtx, kafkaSubscribeTimeout)
		defer cancel()

		if err != nil {
			log.Errorf("got error on inspection event: %v", err)
			return
		}

		var event inspection.Event
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			log.Errorf("failed to unmarshal inspection event: %v", err)
			return
		}

		switch event.Type {
		case inspection.EventTypeCreate:
			err = s.handleCreatedInspection(ctx, event.Inspection)
		case inspection.EventTypeUpdate:
			err = s.handleUpdatedInspection(ctx, log, event.Inspection)
		default:
			err = fmt.Errorf("unknown event type: %v", event.Type)
		}

		if err != nil {
			log.Errorf("failed to handle inspection event (type = %d): %v", event.Type, err)
			return
		}
	}
}

func (s *Service) handleCreatedInspection(ctx context.Context, ins inspection.Inspection) error {
	if ins.Status != inspection.StatusInWork {
		return fmt.Errorf("invalid inspection status: %v", ins.Status)
	}

	return nil
}

func (s *Service) handleUpdatedInspection(ctx context.Context, log golog.Logger, ins inspection.Inspection) error {
	if ins.Status != inspection.StatusDone {
		return nil
	}

	dbTask, err := s.taskRepository.FinishTask(ctx, ins.TaskID)
	if err != nil {
		return fmt.Errorf("finish task %d: %w", ins.TaskID, err)
	}

	go s.taskPublisher.Publish(goctx.Wrap(ctx), log, EventTypeUpdate, MapFromDB(dbTask))

	return nil
}
