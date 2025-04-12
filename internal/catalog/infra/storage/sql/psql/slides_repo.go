package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/wintermonth2298/library-ddd/internal/catalog/domain"
	"github.com/wintermonth2298/library-ddd/internal/catalog/infra/storage/sql/mapping"
	"github.com/wintermonth2298/library-ddd/internal/catalog/infra/storage/sql/uow"
)

type SlidesRepo struct {
	db *sqlx.DB
}

func NewSlidesRepo(db *sqlx.DB) *SlidesRepo {
	return &SlidesRepo{db: db}
}

func (r *SlidesRepo) GetSlidesByCaseID(ctx context.Context, caseID uuid.UUID) ([]domain.Slide, error) {
	exec := uow.Executor(ctx, r.db)

	query := `
		SELECT id, version, preparation_status, case_id
		FROM slides
		WHERE case_id = $1
	`

	var models []mapping.SlideModel
	err := exec.SelectContext(ctx, &models, query, caseID.String())
	if err != nil {
		return nil, fmt.Errorf("select slides by caseID: %w", err)
	}

	slides := make([]domain.Slide, 0, len(models))
	for _, model := range models {
		slides = append(slides, mapping.ToDomainSlide(model))
	}

	return slides, nil
}

func (r *SlidesRepo) GetSlide(ctx context.Context, id uuid.UUID) (domain.Slide, error) {
	exec := uow.Executor(ctx, r.db)

	var model mapping.SlideModel
	err := exec.GetContext(ctx, &model, `
		SELECT id, version, preparation_status, case_id
		FROM slides
		WHERE id = $1
	`, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Slide{}, domain.ErrSlideNotFound
		}
		return domain.Slide{}, fmt.Errorf("select slide: %w", err)
	}

	return mapping.ToDomainSlide(model), nil
}

func (r *SlidesRepo) SaveSlide(ctx context.Context, s domain.Slide) error {
	exec := uow.Executor(ctx, r.db)

	model := mapping.ToModelSlide(s)

	if s.Version == 0 {
		insertQuery := `
			INSERT INTO slides (id, version, preparation_status, case_id)
			VALUES (:id, 1, :preparation_status, :case_id)
		`
		_, err := exec.NamedExecContext(ctx, insertQuery, model)
		if err != nil {
			return fmt.Errorf("insert slide: %w", err)
		}
		return nil
	}

	updateQuery := `
		UPDATE slides
		SET version = version + 1,
			preparation_status = :preparation_status
		WHERE id = :id AND version = :version
	`

	result, err := exec.NamedExecContext(ctx, updateQuery, model)
	if err != nil {
		return fmt.Errorf("update slide: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrVersionConflict
	}

	return nil
}
