package task

import (
	"context"
	"reflect"
	"testing"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/pagination"
)

type mockRepository struct {
	task  Task
	tasks []Task
}

func (m *mockRepository) Add(context.Context, AddRequest) (Task, error) {
	return Task{}, nil
}

func (m *mockRepository) GetByID(context.Context, int) (Task, error) {
	return m.task, nil
}

func (m *mockRepository) GetByBrigade(context.Context, int, pagination.Pagination) ([]Task, error) {
	return m.tasks, nil
}

func (m *mockRepository) GetAll(context.Context, pagination.Pagination) ([]Task, error) {
	return nil, nil
}

func (m *mockRepository) StartTask(context.Context, int) (Task, error) {
	return Task{}, nil
}

func (m *mockRepository) FinishTask(context.Context, int) (Task, error) {
	return Task{}, nil
}

func (m *mockRepository) AssignToBrigade(context.Context, int, int) (Task, error) {
	return Task{}, nil
}

type mockSubscriberService struct {
	gotObjectIDs []int
	contracts    []Contract
}

func (m *mockSubscriberService) GetLastContractsByObjectIDs(_ goctx.Context, objectIDs []int) ([]Contract, error) {
	m.gotObjectIDs = append([]int(nil), objectIDs...)
	return m.contracts, nil
}

func TestGetByBrigadeExtendedGetsContractsInBatch(t *testing.T) {
	repository := &mockRepository{tasks: []Task{
		{ID: 1, ObjectID: 10},
		{ID: 2, ObjectID: 20},
	}}
	subscriberService := &mockSubscriberService{contracts: []Contract{
		{ID: 100, Object: Object{ID: 10}},
		{ID: 200, Object: Object{ID: 20}},
	}}
	service := NewService(repository, nil, subscriberService)

	got, err := service.GetByBrigadeExtended(goctx.Wrap(context.Background()), 7, pagination.Pagination{})
	if err != nil {
		t.Fatalf("GetByBrigadeExtended returned error: %v", err)
	}

	if !reflect.DeepEqual(subscriberService.gotObjectIDs, []int{10, 20}) {
		t.Fatalf("subscriber service object ids = %v, want [10 20]", subscriberService.gotObjectIDs)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if got[0].Task.ID != 1 || got[0].Contract.ID != 100 {
		t.Fatalf("first extended task = %+v, want task 1 with contract 100", got[0])
	}
	if got[1].Task.ID != 2 || got[1].Contract.ID != 200 {
		t.Fatalf("second extended task = %+v, want task 2 with contract 200", got[1])
	}
}

func TestGetByIDExtendedGetsTaskContract(t *testing.T) {
	repository := &mockRepository{task: Task{ID: 1, ObjectID: 10}}
	subscriberService := &mockSubscriberService{contracts: []Contract{
		{ID: 100, Object: Object{ID: 10}},
	}}
	service := NewService(repository, nil, subscriberService)

	got, err := service.GetByIDExtended(goctx.Wrap(context.Background()), 1)
	if err != nil {
		t.Fatalf("GetByIDExtended returned error: %v", err)
	}

	if !reflect.DeepEqual(subscriberService.gotObjectIDs, []int{10}) {
		t.Fatalf("subscriber service object ids = %v, want [10]", subscriberService.gotObjectIDs)
	}
	if got.Task.ID != 1 || got.Contract.ID != 100 {
		t.Fatalf("extended task = %+v, want task 1 with contract 100", got)
	}
}
