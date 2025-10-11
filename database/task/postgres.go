package task

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(db *sqlx.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

//go:embed sql/add.sql
var addSQL string

func (p *Postgres) Add(ctx context.Context, task Task) (newTask Task, err error) {
	rows, err := p.db.NamedQueryContext(ctx, addSQL, task)
	if err != nil {
		return Task{}, fmt.Errorf("p.db.NamedQueryContext: %w", err)
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	if !rows.Next() {
		return Task{}, errors.New("rows.Next == false")
	}

	err = rows.StructScan(&newTask)
	if err != nil {
		return Task{}, fmt.Errorf("rows.Scan: %w", err)
	}

	if err = rows.Err(); err != nil {
		return Task{}, fmt.Errorf("rows.Err: %w", err)
	}

	return newTask, err
}

//go:embed sql/get_by_id.sql
var getByIDSQL string

func (p *Postgres) GetByID(ctx context.Context, id int) (task Task, err error) {
	err = p.db.GetContext(ctx, &task, getByIDSQL, id)
	if err != nil {
		return Task{}, fmt.Errorf("p.db.GetContext: %w", err)
	}

	return task, nil
}

//go:embed sql/get_by_brigade.sql
var getByBrigadeSQL string

func (p *Postgres) GetByBrigade(ctx context.Context, brigadeID int) (tasks []Task, err error) {
	err = p.db.SelectContext(ctx, &tasks, getByBrigadeSQL, brigadeID)
	if err != nil {
		return nil, fmt.Errorf("p.db.SelectContext: %w", err)
	}

	return tasks, nil
}

//go:embed sql/start_task.sql
var startTaskSQL string

func (p *Postgres) StartTask(ctx context.Context, id int) (task Task, err error) {
	err = p.db.GetContext(ctx, &task, startTaskSQL, id)
	if err != nil {
		return Task{}, fmt.Errorf("p.db.GetContext: %w", err)
	}

	return task, nil
}

//go:embed sql/finish_task.sql
var finishTaskSQL string

func (p *Postgres) FinishTask(ctx context.Context, id int) (task Task, err error) {
	err = p.db.GetContext(ctx, &task, finishTaskSQL, id)
	if err != nil {
		return Task{}, fmt.Errorf("p.db.GetContext: %w", err)
	}

	return task, nil
}
