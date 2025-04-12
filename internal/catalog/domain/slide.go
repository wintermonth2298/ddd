package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrSlideNotFound = errors.New("slide not found")

type Slide struct {
	ID                uuid.UUID
	CaseID            uuid.UUID
	Version           Version
	PreparationStatus SlidePreparationStatus

	events []Event
}

func CreateSlide(caseID uuid.UUID) Slide {
	id := uuid.New()

	s := Slide{
		ID:                id,
		Version:           0,
		PreparationStatus: SlidePreparationStatusNotStarted,
		CaseID:            caseID,
	}

	s.addEvent(EventTypeSlideCreated)

	return s
}

func (s *Slide) Finish() {
	s.PreparationStatus = SlidePreparationStatusDone

	s.addEvent(EventTypeSlideUpdated)
}

func (s *Slide) PullEvents() []Event {
	out := s.events
	s.events = nil
	return out
}

func (s *Slide) addEvent(eventType EventType) {
	event := Event{
		ID:        uuid.New(),
		SlideID:   s.ID,
		CaseID:    s.CaseID,
		Type:      eventType,
		CreatedAt: time.Now(),
	}
	s.events = append(s.events, event)
}

type SlidePreparationStatus uint8

const (
	SlidePreparationStatusUnkown SlidePreparationStatus = iota
	SlidePreparationStatusNotStarted
	SlidePreparationStatusProcessing
	SlidePreparationStatusDone
	SlidePreparationStatusError
)
