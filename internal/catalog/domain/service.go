package domain

import (
	"time"

	"github.com/google/uuid"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateSlide(caseID uuid.UUID, caseSlides []Slide) Slide {
	slideID := uuid.New()

	slide := Slide{
		ID:                slideID,
		CaseID:            caseID,
		Version:           0,
		PreparationStatus: SlidePreparationStatusNotStarted,
	}

	slides := append(caseSlides, slide)
	slide.addEvent(EventSlideCreated{
		ID:                    uuid.New(),
		CreationTime:          time.Now(),
		CaseID:                caseID,
		CasePreparationStatus: s.casePreparationStatus(slides),
	})

	return slide
}

func (s *Service) FinishSlide(slide Slide, caseID uuid.UUID, caseSlides []Slide) Slide {
	sl := Slide{
		ID:                slide.ID,
		CaseID:            slide.CaseID,
		Version:           slide.Version,
		PreparationStatus: SlidePreparationStatusDone,
	}

	slides := append(caseSlides, slide)
	sl.addEvent(EvenSlideFinished{
		ID:                    uuid.New(),
		CreationTime:          time.Now(),
		CaseID:                caseID,
		CasePreparationStatus: s.casePreparationStatus(slides),
	})

	return sl
}

func (s *Service) casePreparationStatus(slides []Slide) CasePreparationStatus {
	var (
		anyProcessing = false
		anyErrors     = false
		anyNotStarted = false
	)

	for _, slide := range slides {
		switch slide.PreparationStatus {
		case SlidePreparationStatusProcessing:
			anyProcessing = true
		case SlidePreparationStatusError:
			anyErrors = true
		case SlidePreparationStatusNotStarted:
			anyNotStarted = true
		}
	}

	if anyErrors {
		return CasePreparationStatusError
	}
	if anyProcessing || anyNotStarted {
		return CasePreparationStatusProcessing
	}

	return CasePreparationStatusDone
}
