package inspection

import "time"

type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeCreate
	EventTypeUpdate
)

type Event struct {
	Type       EventType  `json:"Type"`
	Date       time.Time  `json:"Date"`
	UserID     int        `json:"UserID"`
	Inspection Inspection `json:"Inspection"`
}
