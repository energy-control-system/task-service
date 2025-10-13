package task

import "time"

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

type AddRequest struct {
	BrigadeID   *int       `json:"BrigadeID"`
	ObjectID    int        `json:"ObjectID"`
	PlanVisitAt *time.Time `json:"PlanVisitAt"`
	Comment     *string    `json:"Comment"`
}
