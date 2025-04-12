package uow

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type TxKey struct{}

type SQLUnitOfWork struct {
	db *sqlx.DB
}

func NewSQLUnitOfWork(db *sqlx.DB) *SQLUnitOfWork {
	return &SQLUnitOfWork{db: db}
}

func (u *SQLUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, TxKey{}, tx)

	if err := fn(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func Executor(ctx context.Context, db *sqlx.DB) sqlxExecutor {
	if tx, ok := ctx.Value(TxKey{}).(*sqlx.Tx); ok && tx != nil {
		return tx
	}
	return db
}

type sqlxExecutor interface {
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
