package task

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"task-service/service/task"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

//go:embed sql/add.sql
var addSQL string

func (r *Repository) Add(ctx context.Context, request task.AddRequest) (task.Task, error) {
	rows, err := r.db.NamedQueryContext(ctx, addSQL, MapAddRequestToDB(request))
	if err != nil {
		return task.Task{}, fmt.Errorf("r.db.NamedQueryContext: %w", err)
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	if !rows.Next() {
		return task.Task{}, errors.New("rows.Next == false")
	}

	var t Task
	err = rows.StructScan(&t)
	if err != nil {
		return task.Task{}, fmt.Errorf("rows.Scan: %w", err)
	}

	if err = rows.Err(); err != nil {
		return task.Task{}, fmt.Errorf("rows.Err: %w", err)
	}

	return MapFromDB(t), err
}

//go:embed sql/get_by_id.sql
var getByIDSQL string

func (r *Repository) GetByID(ctx context.Context, id int) (task.Task, error) {
	var t Task
	err := r.db.GetContext(ctx, &t, getByIDSQL, id)
	if err != nil {
		return task.Task{}, fmt.Errorf("r.db.GetContext: %w", err)
	}

	return MapFromDB(t), nil
}

//go:embed sql/get_by_brigade.sql
var getByBrigadeSQL string

func (r *Repository) GetByBrigade(ctx context.Context, brigadeID int) ([]task.Task, error) {
	var tasks []Task
	err := r.db.SelectContext(ctx, &tasks, getByBrigadeSQL, brigadeID)
	if err != nil {
		return nil, fmt.Errorf("r.db.SelectContext: %w", err)
	}

	return MapSliceFromDB(tasks), nil
}

//go:embed sql/start_task.sql
var startTaskSQL string

func (r *Repository) StartTask(ctx context.Context, id int) (task.Task, error) {
	var t Task
	err := r.db.GetContext(ctx, &t, startTaskSQL, id)
	if err != nil {
		return task.Task{}, fmt.Errorf("r.db.GetContext: %w", err)
	}

	return MapFromDB(t), nil
}

//go:embed sql/finish_task.sql
var finishTaskSQL string

func (r *Repository) FinishTask(ctx context.Context, id int) (task.Task, error) {
	var t Task
	err := r.db.GetContext(ctx, &t, finishTaskSQL, id)
	if err != nil {
		return task.Task{}, fmt.Errorf("r.db.GetContext: %w", err)
	}

	return MapFromDB(t), nil
}
