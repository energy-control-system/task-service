package task

import "time"

type Status struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Task struct {
	ID          int        `db:"id"`
	BrigadeID   *int       `db:"brigade_id"`
	ObjectID    int        `db:"object_id"`
	PlanVisitAt *time.Time `db:"plan_visit_at"`
	Status      int        `db:"status"`
	Comment     *string    `db:"comment"`
	StartedAt   *time.Time `db:"started_at"`
	FinishedAt  *time.Time `db:"finished_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

type AddRequest struct {
	BrigadeID   *int       `db:"brigade_id"`
	ObjectID    int        `db:"object_id"`
	PlanVisitAt *time.Time `db:"plan_visit_at"`
	Comment     *string    `db:"comment"`
}
