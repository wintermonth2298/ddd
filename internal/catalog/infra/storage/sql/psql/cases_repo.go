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
)

type CasesRepo struct {
	db *sqlx.DB
}

func NewCasesRepo(db *sqlx.DB) *CasesRepo {
	return &CasesRepo{db: db}
}

func (r *CasesRepo) GetCase(ctx context.Context, id uuid.UUID) (domain.Case, error) {
	exec := executor(ctx, r.db)

	var model mapping.CaseModel
	err := exec.GetContext(ctx, &model, `
		SELECT id, version
		FROM cases
		WHERE id = $1
	`, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Case{}, domain.ErrCaseNotFound
		}
		return domain.Case{}, fmt.Errorf("select case: %w", err)
	}

	return mapping.ToDomainCase(model), nil
}

func (r *CasesRepo) SaveCase(ctx context.Context, c domain.Case) error {
	exec := executor(ctx, r.db)

	model := mapping.ToModelCase(c)

	if c.Version == 0 {
		insertQuery := `
			INSERT INTO cases (id, version)
			VALUES (:id, 1)
		`
		_, err := exec.NamedExecContext(ctx, insertQuery, model)
		if err != nil {
			return fmt.Errorf("insert case: %w", err)
		}
		return nil
	}

	updateQuery := `
		UPDATE cases
		SET version = version + 1
		WHERE id = :id AND version = :version
	`

	result, err := exec.NamedExecContext(ctx, updateQuery, model)
	if err != nil {
		return fmt.Errorf("update case: %w", err)
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
