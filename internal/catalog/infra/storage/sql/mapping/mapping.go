package mapping

import (
	"github.com/google/uuid"
	"github.com/wintermonth2298/library-ddd/internal/catalog/domain"
)

type CaseModel struct {
	ID      string `db:"id"`
	Version int    `db:"version"`
}

func ToDomainCase(model CaseModel) domain.Case {
	uid, _ := uuid.Parse(model.ID)

	return domain.Case{
		ID:      uid,
		Version: domain.Version(model.Version),
	}
}

func ToModelCase(c domain.Case) CaseModel {
	return CaseModel{
		ID:      c.ID.String(),
		Version: int(c.Version),
	}
}

type SlideModel struct {
	ID                string `db:"id"`
	Version           int    `db:"version"`
	PreparationStatus uint8  `db:"preparation_status"`
	CaseID            string `db:"case_id"`
}

func ToDomainSlide(model SlideModel) domain.Slide {
	uid, _ := uuid.Parse(model.ID)
	caseuid, _ := uuid.Parse(model.CaseID)

	return domain.Slide{
		ID:                uid,
		Version:           domain.Version(model.Version),
		PreparationStatus: ToDomainSlidePreparationStatus(model.PreparationStatus),
		CaseID:            caseuid,
	}
}

func ToModelSlide(slide domain.Slide) SlideModel {
	return SlideModel{
		ID:                slide.ID.String(),
		Version:           int(slide.Version),
		PreparationStatus: ToModelSlidePreparationStatus(slide.PreparationStatus),
		CaseID:            slide.CaseID.String(),
	}
}

func ToDomainSlidePreparationStatus(status uint8) domain.SlidePreparationStatus {
	switch status {
	case 1:
		return domain.SlidePreparationStatusNotStarted
	case 2:
		return domain.SlidePreparationStatusProcessing
	case 3:
		return domain.SlidePreparationStatusDone
	case 4:
		return domain.SlidePreparationStatusError
	}
	return domain.SlidePreparationStatusUnkown
}

func ToModelSlidePreparationStatus(status domain.SlidePreparationStatus) uint8 {
	switch status {
	case domain.SlidePreparationStatusNotStarted:
		return 1
	case domain.SlidePreparationStatusProcessing:
		return 2
	case domain.SlidePreparationStatusDone:
		return 3
	case domain.SlidePreparationStatusError:
		return 4
	}
	return 0
}

func ToDomainCasePreparationStatus(status uint8) domain.CasePreparationStatus {
	switch status {
	case 1:
		return domain.CasePreparationStatusNotStarted
	case 2:
		return domain.CasePreparationStatusProcessing
	case 3:
		return domain.CasePreparationStatusDone
	case 4:
		return domain.CasePreparationStatusError
	}
	return domain.CasePreparationStatusUnkown
}

func ToModelCasePreparationStatus(status domain.CasePreparationStatus) uint8 {
	switch status {
	case domain.CasePreparationStatusNotStarted:
		return 1
	case domain.CasePreparationStatusProcessing:
		return 2
	case domain.CasePreparationStatusDone:
		return 3
	case domain.CasePreparationStatusError:
		return 4
	}
	return 0
}
