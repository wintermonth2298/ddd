package application

import (
	"context"
	"fmt"

	"github.com/wintermonth2298/library-ddd/internal/catalog/domain"
)

const maxUnpublishedEventsToFetch = 10

type EventHandler func(ctx context.Context, e domain.Event) error

type eventsProcessor struct {
	storage  eventsStorage
	handlers map[domain.EventType][]EventHandler
}

func newEventsProcessor(storage eventsStorage) *eventsProcessor {
	return &eventsProcessor{
		storage:  storage,
		handlers: make(map[domain.EventType][]EventHandler),
	}
}

func (p *eventsProcessor) Register(t domain.EventType, h EventHandler) {
	p.handlers[t] = append(p.handlers[t], h)
}

func (p *eventsProcessor) Process(ctx context.Context) error {
	events, err := p.storage.FetchUnpublished(ctx, maxUnpublishedEventsToFetch)
	if err != nil {
		return fmt.Errorf("fetch unpublished events: %w", err)
	}

	if err := p.publish(ctx, events); err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	if err := p.storage.MarkPublished(ctx, events); err != nil {
		return fmt.Errorf("mark published: %w", err)
	}

	return nil
}

func (p *eventsProcessor) publish(ctx context.Context, events []domain.Event) error {
	for _, e := range events {
		for _, handler := range p.handlers[e.Type] {
			if err := handler(ctx, e); err != nil {
				return fmt.Errorf("handle event %v: %w", e.Type, err)
			}
		}
	}
	return nil
}
