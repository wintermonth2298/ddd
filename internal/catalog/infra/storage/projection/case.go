package projection

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/wintermonth2298/library-ddd/internal/catalog/domain"
	"github.com/wintermonth2298/library-ddd/internal/catalog/infra/storage/sql/mapping"
)

type CaseProjector struct {
	db      *sqlx.DB
	service *domain.Service
}

func NewCaseProjectior(db *sqlx.DB, service *domain.Service) *CaseProjector {
	return &CaseProjector{
		db:      db,
		service: service,
	}
}

type CaseProjection struct {
	ID     string `db:"id"`
	Status uint8  `db:"status"`
}

func (p *CaseProjector) HandleSlideCreated(ctx context.Context, event domain.Event) error {
	projection, err := p.buildCaseProjection(event)
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
	projection, err := p.buildCaseProjection(event)
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

func (p *CaseProjector) buildCaseProjection(event domain.Event) (CaseProjection, error) {
	switch e := event.(type) {
	case domain.EventSlideCreated:
		return CaseProjection{
			ID:     e.CaseID.String(),
			Status: mapping.ToModelCasePreparationStatus(e.CasePreparationStatus),
		}, nil
	case domain.EvenSlideFinished:
		return CaseProjection{
			ID:     e.CaseID.String(),
			Status: mapping.ToModelCasePreparationStatus(e.CasePreparationStatus),
		}, nil
	}

	return CaseProjection{}, fmt.Errorf("unknown event: %s", event.Name())
}
