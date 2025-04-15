package psql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/wintermonth2298/library-ddd/internal/catalog/domain"
	"github.com/wintermonth2298/library-ddd/internal/catalog/infra/storage/sql/mapping"
)

type EventsStorage struct {
	db *sqlx.DB
}

func NewEventsStorage(db *sqlx.DB) *EventsStorage {
	return &EventsStorage{db: db}
}

func (s *EventsStorage) MarkPublished(ctx context.Context, events []domain.Event) error {
	exec := executor(ctx, s.db)

	if len(events) == 0 {
		return nil
	}

	ids := make([]string, 0, len(events))
	for _, e := range events {
		ids = append(ids, e.ID.String())
	}

	const tmpl = `
		UPDATE events
		SET published = true
		WHERE id IN (?)
	`

	query, args, err := sqlx.In(tmpl, ids)
	if err != nil {
		return fmt.Errorf("prepare mark published: %w", err)
	}
	query = s.db.Rebind(query)

	_, err = exec.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("mark published: %w", err)
	}
	return nil
}

func (s *EventsStorage) FetchUnpublished(ctx context.Context, limit int) ([]domain.Event, error) {
	exec := executor(ctx, s.db)

	const query = `
		SELECT id, type, created_at, published, case_id, slide_id
		FROM events
		WHERE NOT published
		ORDER BY created_at
		LIMIT $1
	`

	var events []mapping.EventModel
	err := exec.SelectContext(ctx, &events, query, limit)
	if err != nil {
		return nil, fmt.Errorf("fetch outbox: %w", err)
	}

	domainEvents := make([]domain.Event, 0, len(events))
	for _, e := range events {
		domainEvents = append(domainEvents, mapping.ToDomainEvent(e))
	}

	return domainEvents, nil
}

func (s *EventsStorage) Add(ctx context.Context, events []domain.Event) error {
	exec := executor(ctx, s.db)

	eventModels := make([]mapping.EventModel, 0, len(events))
	for _, event := range events {
		eventModels = append(eventModels, mapping.ToModelEvent(event, false))
	}

	const query = `
		INSERT INTO events (id, type, created_at, published, case_id, slide_id)
		VALUES (:id, :type, :created_at, :published, :case_id, :slide_id)
	`

	for _, e := range eventModels {
		if _, err := exec.NamedExecContext(ctx, query, e); err != nil {
			return fmt.Errorf("insert outbox event: %w", err)
		}
	}

	return nil
}
