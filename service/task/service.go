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
	"github.com/sunshineOfficial/golib/pagination"
)

const kafkaSubscribeTimeout = 2 * time.Minute

type Service struct {
	repository        Repository
	publisher         *Publisher
	subscriberService SubscriberService
}

func NewService(repository Repository, publisher *Publisher, subscriberService SubscriberService) *Service {
	return &Service{
		repository:        repository,
		publisher:         publisher,
		subscriberService: subscriberService,
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

func (s *Service) GetByIDExtended(ctx goctx.Context, id int) (TaskExtended, error) {
	t, err := s.GetByID(ctx, id)
	if err != nil {
		return TaskExtended{}, err
	}

	contracts, err := s.subscriberService.GetLastContractsByObjectIDs(ctx, []int{t.ObjectID})
	if err != nil {
		return TaskExtended{}, fmt.Errorf("get last contract by object id: %w", err)
	}

	contractsByObjectID := make(map[int]Contract, len(contracts))
	for _, c := range contracts {
		contractsByObjectID[c.Object.ID] = c
	}

	return TaskExtended{
		Task:     t,
		Contract: contractsByObjectID[t.ObjectID],
	}, nil
}

func (s *Service) GetByBrigade(ctx goctx.Context, brigadeID int, page pagination.Pagination) ([]Task, error) {
	if err := page.Validate(); err != nil {
		return nil, fmt.Errorf("validate pagination: %w", err)
	}

	tasks, err := s.repository.GetByBrigade(ctx, brigadeID, page)
	if err != nil {
		return nil, fmt.Errorf("get tasks by brigade id %d from db: %w", brigadeID, err)
	}

	return tasks, nil
}

func (s *Service) GetByBrigadeExtended(ctx goctx.Context, brigadeID int, page pagination.Pagination) ([]TaskExtended, error) {
	tasks, err := s.GetByBrigade(ctx, brigadeID, page)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return []TaskExtended{}, nil
	}

	objectIDs := make([]int, 0, len(tasks))
	for _, t := range tasks {
		objectIDs = append(objectIDs, t.ObjectID)
	}

	contracts, err := s.subscriberService.GetLastContractsByObjectIDs(ctx, objectIDs)
	if err != nil {
		return nil, fmt.Errorf("get last contracts by object ids: %w", err)
	}

	contractsByObjectID := make(map[int]Contract, len(contracts))
	for _, c := range contracts {
		contractsByObjectID[c.Object.ID] = c
	}

	result := make([]TaskExtended, 0, len(tasks))
	for _, t := range tasks {
		result = append(result, TaskExtended{
			Task:     t,
			Contract: contractsByObjectID[t.ObjectID],
		})
	}

	return result, nil
}

func (s *Service) GetAll(ctx goctx.Context, page pagination.Pagination) ([]Task, error) {
	if err := page.Validate(); err != nil {
		return nil, fmt.Errorf("validate pagination: %w", err)
	}

	tasks, err := s.repository.GetAll(ctx, page)
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

func (s *Service) AssignToBrigade(ctx goctx.Context, log golog.Logger, request AssignToBrigadeRequest) (Task, error) {
	t, err := s.repository.GetByID(ctx, request.TaskID)
	if err != nil {
		return Task{}, fmt.Errorf("get task %d from db: %w", request.TaskID, err)
	}

	if t.Status != StatusPlanned {
		return Task{}, fmt.Errorf("invalid task status: %v, expected planned", t.Status)
	}

	t, err = s.repository.AssignToBrigade(ctx, request.TaskID, request.BrigadeID)
	if err != nil {
		return Task{}, fmt.Errorf("assign task %d to brigade %d: %w", request.TaskID, request.BrigadeID, err)
	}

	go s.publisher.Publish(ctx, log, EventTypeAssign, t)

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
