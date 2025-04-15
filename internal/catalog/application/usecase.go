package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wintermonth2298/library-ddd/internal/catalog/domain"
)

type casesRepo interface {
	GetCase(ctx context.Context, caseID uuid.UUID) (domain.Case, error)
	SaveCase(ctx context.Context, c domain.Case) error
}

type slidesRepo interface {
	GetSlide(ctx context.Context, id uuid.UUID) (domain.Slide, error)
	SaveSlide(ctx context.Context, s domain.Slide) error
}

type eventsStorage interface {
	AddEvent(ctx context.Context, events []domain.Event) error
	MarkEventPublished(ctx context.Context, events []domain.Event) error
	FetchUnpublishedEvents(ctx context.Context, limit int) ([]domain.Event, error)
}

type storage interface {
	casesRepo
	slidesRepo

	eventsStorage

	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewUsecases(storage storage) *Usecases {
	return &Usecases{
		storage:         storage,
		eventsProcessor: newEventsProcessor(storage),
	}
}

type Usecases struct {
	storage         storage
	eventsProcessor *eventsProcessor
}

func (u *Usecases) CreateCase(ctx context.Context) error {
	c := domain.CreateCase()
	if err := u.storage.SaveCase(ctx, c); err != nil {
		return fmt.Errorf("save case: %w", err)
	}

	return nil
}

func (u *Usecases) AddSlide(ctx context.Context, caseID uuid.UUID) error {
	return u.storage.WithTx(ctx, func(ctx context.Context) error {
		slide := domain.CreateSlide(caseID)

		_, err := u.storage.GetCase(ctx, caseID)
		if err != nil {
			return fmt.Errorf("get case: %w", err)
		}

		if err := u.storage.SaveSlide(ctx, slide); err != nil {
			return fmt.Errorf("save slide: %w", err)
		}

		if err := u.storage.AddEvent(ctx, slide.PullEvents()); err != nil {
			return fmt.Errorf("add events: %w", err)
		}

		return nil
	})
}

func (u *Usecases) FinishSlide(ctx context.Context, slideID uuid.UUID) error {
	return u.storage.WithTx(ctx, func(ctx context.Context) error {
		slide, err := u.storage.GetSlide(ctx, slideID)
		if err != nil {
			return fmt.Errorf("get slide: %w", err)
		}

		slide.Finish()

		if err := u.storage.AddEvent(ctx, slide.PullEvents()); err != nil {
			return fmt.Errorf("add events: %w", err)
		}

		if err := u.storage.SaveSlide(ctx, slide); err != nil {
			return fmt.Errorf("save slide: %w", err)
		}

		return nil
	})
}

func (u *Usecases) RegisterEventHandler(t domain.EventType, h EventHandler) {
	u.eventsProcessor.Register(t, h)
}

func (u *Usecases) StartEventsProcessor(interval time.Duration) {
	ctx := context.TODO()
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := u.eventsProcessor.Process(ctx); err != nil {
					log.Printf("event processing failed: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
