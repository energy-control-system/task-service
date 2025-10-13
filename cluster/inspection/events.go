package inspection

import "time"

type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeStart
	EventTypeFinish
)

type Event struct {
	Type       EventType  `json:"Type"`
	Date       time.Time  `json:"Date"`
	UserID     int        `json:"UserID"`
	Inspection Inspection `json:"Inspection"`
}
