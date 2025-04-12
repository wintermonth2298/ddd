package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventType uint8

const (
	EventTypeUnknown EventType = iota
	EventTypeSlideCreated
	EventTypeSlideUpdated
)

type Event struct {
	ID        uuid.UUID
	SlideID   uuid.UUID
	CaseID    uuid.UUID
	Type      EventType
	CreatedAt time.Time
}
