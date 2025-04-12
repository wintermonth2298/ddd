package domain

import (
	"errors"

	"github.com/google/uuid"
)

var ErrCaseNotFound = errors.New("case not found")

type Case struct {
	ID      uuid.UUID
	Version Version
}

func CreateCase() Case {
	return Case{
		ID:      uuid.New(),
		Version: 0,
	}
}

type CasePreparationStatus uint8

const (
	CasePreparationStatusUnkown CasePreparationStatus = iota
	CasePreparationStatusNotStarted
	CasePreparationStatusProcessing
	CasePreparationStatusDone
	CasePreparationStatusError
)
