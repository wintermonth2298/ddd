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
	Add(ctx context.Context, events []domain.Event) error
	MarkPublished(ctx context.Context, events []domain.Event) error
	FetchUnpublished(ctx context.Context, limit int) ([]domain.Event, error)
}

type uow interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewUsecases(
	casesRepo casesRepo,
	slidesRepo slidesRepo,
	eventsStorage eventsStorage,
	uow uow,
) *Usecases {
	return &Usecases{
		casesRepo:       casesRepo,
		slidesRepo:      slidesRepo,
		eventsStorage:   eventsStorage,
		eventsProcessor: newEventsProcessor(eventsStorage),
		uow:             uow,
	}
}

type Usecases struct {
	casesRepo       casesRepo
	slidesRepo      slidesRepo
	eventsStorage   eventsStorage
	eventsProcessor *eventsProcessor
	uow             uow
}

func (u *Usecases) CreateCase(ctx context.Context) error {
	c := domain.CreateCase()
	if err := u.casesRepo.SaveCase(ctx, c); err != nil {
		return fmt.Errorf("save case: %w", err)
	}

	return nil
}

func (u *Usecases) AddSlide(ctx context.Context, caseID uuid.UUID) error {
	return u.uow.Do(ctx, func(ctx context.Context) error {
		slide := domain.CreateSlide(caseID)

		_, err := u.casesRepo.GetCase(ctx, caseID)
		if err != nil {
			return fmt.Errorf("get case: %w", err)
		}

		if err := u.slidesRepo.SaveSlide(ctx, slide); err != nil {
			return fmt.Errorf("save slide: %w", err)
		}

		if err := u.eventsStorage.Add(ctx, slide.PullEvents()); err != nil {
			return fmt.Errorf("add events: %w", err)
		}

		return nil
	})
}

func (u *Usecases) FinishSlide(ctx context.Context, slideID uuid.UUID) error {
	return u.uow.Do(ctx, func(ctx context.Context) error {
		slide, err := u.slidesRepo.GetSlide(ctx, slideID)
		if err != nil {
			return fmt.Errorf("get slide: %w", err)
		}

		slide.Finish()

		if err := u.eventsStorage.Add(ctx, slide.PullEvents()); err != nil {
			return fmt.Errorf("add events: %w", err)
		}

		if err := u.slidesRepo.SaveSlide(ctx, slide); err != nil {
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
