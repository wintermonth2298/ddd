package domain

import (
	"errors"

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

func (s *Slide) PullEvents() []Event {
	out := s.events
	s.events = nil
	return out
}

func (s *Slide) addEvent(event Event) {
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
