package projection

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/wintermonth2298/library-ddd/internal/catalog/domain"
	"github.com/wintermonth2298/library-ddd/internal/catalog/infra/storage/sql/mapping"
)

type slidesRepo interface {
	GetSlidesByCaseID(ctx context.Context, caseID uuid.UUID) ([]domain.Slide, error)
}

type casesRepo interface {
	GetCase(ctx context.Context, caseID uuid.UUID) (domain.Case, error)
}

type CaseProjector struct {
	db         *sqlx.DB
	slidesRepo slidesRepo
	casesRepo  casesRepo
	service    *domain.Service
}

func NewCaseProjectior(
	db *sqlx.DB,
	service *domain.Service,
	slidesRepo slidesRepo,
	casesRepo casesRepo,
) *CaseProjector {
	return &CaseProjector{
		db:         db,
		service:    service,
		slidesRepo: slidesRepo,
		casesRepo:  casesRepo,
	}
}

type CaseProjection struct {
	ID     string `db:"id"`
	Status uint8  `db:"status"`
}

func (p *CaseProjector) HandleSlideCreated(ctx context.Context, event domain.Event) error {
	projection, err := p.buildCaseProjection(ctx, event.CaseID)
	if err != nil {
		return fmt.Errorf("build case projection: %w", err)
	}

	query := `
		INSERT INTO case_projections (id, status)
		VALUES (:id, :status)
		ON CONFLICT (id) DO NOTHING
	`

	_, err = p.db.NamedExecContext(ctx, query, projection)
	if err != nil {
		return fmt.Errorf("insert case projection: %w", err)
	}

	return nil
}

func (p *CaseProjector) HandleSlideUpdated(ctx context.Context, event domain.Event) error {
	projection, err := p.buildCaseProjection(ctx, event.CaseID)
	if err != nil {
		return fmt.Errorf("build case projection: %w", err)
	}

	query := `
		UPDATE case_projections
		SET status = :status
		WHERE id = :id
		  AND status IS DISTINCT FROM :status
	`

	_, err = p.db.NamedExecContext(ctx, query, projection)
	if err != nil {
		return fmt.Errorf("update case projection: %w", err)
	}

	return nil
}

func (p *CaseProjector) buildCaseProjection(ctx context.Context, caseID uuid.UUID) (CaseProjection, error) {
	c, err := p.casesRepo.GetCase(ctx, caseID)
	if err != nil {
		return CaseProjection{}, fmt.Errorf("get case: %w", err)
	}

	slides, err := p.slidesRepo.GetSlidesByCaseID(ctx, caseID)
	if err != nil {
		return CaseProjection{}, fmt.Errorf("get slides by case id: %w", err)
	}

	status := p.service.CasePreparationStatus(ctx, slides)

	return CaseProjection{
		ID:     c.ID.String(),
		Status: mapping.ToModelCasePreparationStatus(status),
	}, nil
}
