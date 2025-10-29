package task

import "context"

type Repository interface {
	Add(ctx context.Context, request AddRequest) (Task, error)
	GetByID(ctx context.Context, id int) (Task, error)
	GetByBrigade(ctx context.Context, brigadeID int) ([]Task, error)
	GetAll(ctx context.Context) ([]Task, error)
	StartTask(ctx context.Context, id int) (Task, error)
	FinishTask(ctx context.Context, id int) (Task, error)
}
