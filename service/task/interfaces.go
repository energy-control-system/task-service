package task

import (
	"context"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/pagination"
)

type Repository interface {
	Add(ctx context.Context, request AddRequest) (Task, error)
	GetByID(ctx context.Context, id int) (Task, error)
	GetByBrigade(ctx context.Context, brigadeID int, page pagination.Pagination, filter GetAllFilter) ([]Task, error)
	GetAll(ctx context.Context, page pagination.Pagination, filter GetAllFilter) ([]Task, error)
	StartTask(ctx context.Context, id int) (Task, error)
	FinishTask(ctx context.Context, id int) (Task, error)
	AssignToBrigade(ctx context.Context, taskID, brigadeID int) (Task, error)
}

type SubscriberService interface {
	GetLastContractsByObjectIDs(ctx goctx.Context, objectIDs []int) ([]Contract, error)
}
