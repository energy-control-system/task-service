package task

import (
	"task-service/database/task"
)

func MapFromDB(t task.Task) Task {
	return Task{
		ID:          t.ID,
		BrigadeID:   t.BrigadeID,
		ObjectID:    t.ObjectID,
		PlanVisitAt: t.PlanVisitAt,
		Status:      t.Status,
		Comment:     t.Comment,
		StartedAt:   t.StartedAt,
		FinishedAt:  t.FinishedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
