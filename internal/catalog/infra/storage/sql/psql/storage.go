package psql

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/wintermonth2298/library-ddd/internal/catalog/domain"
)

type Storage struct {
	casesRepo     *CasesRepo
	slidesRepo    *SlidesRepo
	eventsStorage *EventsStorage
	txManager     *TxManager
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		casesRepo:     NewCasesRepo(db),
		slidesRepo:    NewSlidesRepo(db),
		eventsStorage: NewEventsStorage(db),
		txManager:     NewTxManager(db),
	}
}

func (s *Storage) GetCase(ctx context.Context, caseID uuid.UUID) (domain.Case, error) {
	return s.casesRepo.GetCase(ctx, caseID)
}

func (s *Storage) SaveCase(ctx context.Context, c domain.Case) error {
	return s.casesRepo.SaveCase(ctx, c)
}

func (s *Storage) GetSlide(ctx context.Context, id uuid.UUID) (domain.Slide, error) {
	return s.slidesRepo.GetSlide(ctx, id)
}

func (s *Storage) SaveSlide(ctx context.Context, slide domain.Slide) error {
	return s.slidesRepo.SaveSlide(ctx, slide)
}

func (s *Storage) AddEvent(ctx context.Context, events []domain.Event) error {
	return s.eventsStorage.Add(ctx, events)
}

func (s *Storage) MarkEventPublished(ctx context.Context, events []domain.Event) error {
	return s.eventsStorage.MarkPublished(ctx, events)
}

func (s *Storage) FetchUnpublishedEvents(ctx context.Context, limit int) ([]domain.Event, error) {
	return s.eventsStorage.FetchUnpublished(ctx, limit)
}

func (s *Storage) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return s.txManager.Do(ctx, fn)
}
