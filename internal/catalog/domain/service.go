package domain

import (
	"context"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CasePreparationStatus(ctx context.Context, caseSlides []Slide) CasePreparationStatus {
	var (
		anyProcessing = false
		anyErrors     = false
		anyNotStarted = false
	)

	for _, slide := range caseSlides {
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
