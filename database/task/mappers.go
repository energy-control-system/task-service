package task

import "task-service/service/task"

func MapFromDB(t Task) task.Task {
	return task.Task{
		ID:          t.ID,
		BrigadeID:   t.BrigadeID,
		ObjectID:    t.ObjectID,
		PlanVisitAt: t.PlanVisitAt,
		Status:      task.Status(t.Status),
		Comment:     t.Comment,
		StartedAt:   t.StartedAt,
		FinishedAt:  t.FinishedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func MapSliceFromDB(dbTasks []Task) []task.Task {
	tasks := make([]task.Task, 0, len(dbTasks))
	for _, t := range dbTasks {
		tasks = append(tasks, MapFromDB(t))
	}

	return tasks
}

func MapAddRequestToDB(r task.AddRequest) AddRequest {
	return AddRequest{
		BrigadeID:   r.BrigadeID,
		ObjectID:    r.ObjectID,
		PlanVisitAt: r.PlanVisitAt,
		Comment:     r.Comment,
	}
}
