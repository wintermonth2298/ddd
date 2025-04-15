package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventType uint8

const (
	EventTypeUnknown EventType = iota
	EventTypeSlideCreated
	EventTypeSlideFinished
)

type Event interface {
	CreatedAt() time.Time
	Name() string
	EventID() uuid.UUID
	EventType() EventType
}

type EventSlideCreated struct {
	ID                    uuid.UUID
	CreationTime          time.Time
	CaseID                uuid.UUID
	CasePreparationStatus CasePreparationStatus
}

func (e EventSlideCreated) CreatedAt() time.Time {
	return e.CreationTime
}

func (e EventSlideCreated) Name() string {
	return "event(slide created)"
}

func (e EventSlideCreated) EventID() uuid.UUID {
	return e.ID
}

func (e EventSlideCreated) EventType() EventType {
	return EventTypeSlideCreated
}

type EvenSlideFinished struct {
	ID                    uuid.UUID
	CreationTime          time.Time
	CaseID                uuid.UUID
	CasePreparationStatus CasePreparationStatus
}

func (e EvenSlideFinished) EventID() uuid.UUID {
	return e.ID
}

func (e EvenSlideFinished) CreatedAt() time.Time {
	return e.CreationTime
}

func (e EvenSlideFinished) Name() string {
	return "event(slide finished)"
}

func (e EvenSlideFinished) EventType() EventType {
	return EventTypeSlideFinished
}
