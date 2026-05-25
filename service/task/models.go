package task

import (
	"task-service/cluster/subscriber"
	"time"
)

type Status int

const (
	StatusUnknown Status = iota
	StatusPlanned
	StatusInWork
	StatusDone
)

type Task struct {
	ID          int        `json:"ID"`
	BrigadeID   *int       `json:"BrigadeID,omitempty"`
	ObjectID    int        `json:"ObjectID"`
	PlanVisitAt *time.Time `json:"PlanVisitAt,omitempty"`
	Status      Status     `json:"Status"`
	Comment     *string    `json:"Comment,omitempty"`
	StartedAt   *time.Time `json:"StartedAt,omitempty"`
	FinishedAt  *time.Time `json:"FinishedAt,omitempty"`
	CreatedAt   time.Time  `json:"CreatedAt"`
	UpdatedAt   time.Time  `json:"UpdatedAt"`
}

type Contract = subscriber.Contract
type Object = subscriber.Object

type TaskExtended struct {
	Task     Task     `json:"Task"`
	Contract Contract `json:"Contract"`
}

type AddRequest struct {
	BrigadeID   *int       `json:"BrigadeID"`
	ObjectID    int        `json:"ObjectID"`
	PlanVisitAt *time.Time `json:"PlanVisitAt"`
	Comment     *string    `json:"Comment"`
}

type AssignToBrigadeRequest struct {
	TaskID    int `json:"TaskID"`
	BrigadeID int `json:"BrigadeID"`
}

type GetAllFilter struct {
	Status   *Status
	DateFrom *time.Time
	DateTo   *time.Time
}
