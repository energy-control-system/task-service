package task

import (
	"context"
	"encoding/json"
	"fmt"
	"task-service/cluster/inspection"
	"time"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/gokafka"
	"github.com/sunshineOfficial/golib/golog"
)

const kafkaSubscribeTimeout = 2 * time.Minute

type Service struct {
	repository Repository
	publisher  *Publisher
}

func NewService(repository Repository, publisher *Publisher) *Service {
	return &Service{
		repository: repository,
		publisher:  publisher,
	}
}

func (s *Service) Add(ctx goctx.Context, log golog.Logger, request AddRequest) (Task, error) {
	t, err := s.repository.Add(ctx, request)
	if err != nil {
		return Task{}, fmt.Errorf("add task to db: %w", err)
	}

	go s.publisher.Publish(ctx, log, EventTypeAdd, t)

	return t, nil
}

func (s *Service) GetByID(ctx goctx.Context, id int) (Task, error) {
	t, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return Task{}, fmt.Errorf("get task %d from db: %w", id, err)
	}

	return t, nil
}

func (s *Service) GetByBrigade(ctx goctx.Context, brigadeID int) ([]Task, error) {
	tasks, err := s.repository.GetByBrigade(ctx, brigadeID)
	if err != nil {
		return nil, fmt.Errorf("get tasks by brigade id %d from db: %w", brigadeID, err)
	}

	return tasks, nil
}

func (s *Service) GetAll(ctx goctx.Context) ([]Task, error) {
	tasks, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all tasks from db: %w", err)
	}

	return tasks, nil
}

func (s *Service) StartTask(ctx goctx.Context, log golog.Logger, id int) (Task, error) {
	t, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return Task{}, fmt.Errorf("get task %d from db: %w", id, err)
	}

	if t.Status != StatusPlanned {
		return Task{}, fmt.Errorf("invalid task status: %v", t.Status)
	}

	t, err = s.repository.StartTask(ctx, id)
	if err != nil {
		return Task{}, fmt.Errorf("start task %d: %w", id, err)
	}

	go s.publisher.Publish(ctx, log, EventTypeStart, t)

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
		case inspection.EventTypeStart:
			err = s.handleStartedInspection(ctx, event.Inspection)
		case inspection.EventTypeFinish:
			err = s.handleFinishedInspection(ctx, log, event.Inspection)
		default:
			err = fmt.Errorf("unknown event type: %v", event.Type)
		}

		if err != nil {
			log.Errorf("failed to handle inspection event (type = %d): %v", event.Type, err)
			return
		}
	}
}

func (s *Service) handleStartedInspection(ctx context.Context, ins inspection.Inspection) error {
	if ins.Status != inspection.StatusInWork {
		return fmt.Errorf("invalid inspection status: %v", ins.Status)
	}

	return nil
}

func (s *Service) handleFinishedInspection(ctx context.Context, log golog.Logger, ins inspection.Inspection) error {
	if ins.Status != inspection.StatusDone {
		return fmt.Errorf("invalid inspection status: %v", ins.Status)
	}

	t, err := s.repository.FinishTask(ctx, ins.TaskID)
	if err != nil {
		return fmt.Errorf("finish task %d: %w", ins.TaskID, err)
	}

	go s.publisher.Publish(goctx.Wrap(ctx), log, EventTypeFinish, t)

	return nil
}
