package mapping

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wintermonth2298/library-ddd/internal/catalog/domain"
)

type EventModel struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	Payload   []byte    `db:"payload"`
	Published bool      `db:"published"`
	Type      uint8     `db:"type"`
}

func ToModelEvent(e domain.Event, published bool) (EventModel, error) {
	payload := make(map[string]any)

	switch evt := e.(type) {
	case domain.EventSlideCreated:
		payload["case_id"] = evt.CaseID
		payload["case_preparation_status"] = evt.CasePreparationStatus

	case domain.EvenSlideFinished:
		payload["case_id"] = evt.CaseID
		payload["case_preparation_status"] = evt.CasePreparationStatus

	default:
		return EventModel{}, fmt.Errorf("unknown event type: %T", e)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return EventModel{}, fmt.Errorf("marshal payload: %w", err)
	}

	return EventModel{
		ID:        e.EventID(),
		CreatedAt: e.CreatedAt(),
		Payload:   data,
		Type:      toModelEventType(e.EventType()),
		Published: published,
	}, nil
}

func ToDomainEvent(event EventModel) (domain.Event, error) {
	eventType := toDomainEventType(event.Type)

	switch eventType {

	case domain.EventTypeSlideCreated:
		var payload struct {
			CaseID                uuid.UUID                    `json:"case_id"`
			CasePreparationStatus domain.CasePreparationStatus `json:"case_preparation_status"`
		}

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return nil, fmt.Errorf("unmarshal payload for EventSlideCreated: %w", err)
		}

		return domain.EventSlideCreated{
			ID:                    event.ID,
			CreationTime:          event.CreatedAt,
			CaseID:                payload.CaseID,
			CasePreparationStatus: payload.CasePreparationStatus,
		}, nil

	case domain.EventTypeSlideFinished:
		var payload struct {
			CaseID                uuid.UUID                    `json:"case_id"`
			CasePreparationStatus domain.CasePreparationStatus `json:"case_preparation_status"`
		}

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return nil, fmt.Errorf("unmarshal payload for EventSlideFinished: %w", err)
		}

		return domain.EvenSlideFinished{
			ID:                    event.ID,
			CreationTime:          event.CreatedAt,
			CaseID:                payload.CaseID,
			CasePreparationStatus: payload.CasePreparationStatus,
		}, nil

	default:
		return nil, fmt.Errorf("unknown event type: %d", event.Type)
	}
}

func toModelEventType(et domain.EventType) uint8 {
	switch et {
	case domain.EventTypeSlideCreated:
		return 1
	case domain.EventTypeSlideFinished:
		return 2
	}

	return 0
}

func toDomainEventType(et uint8) domain.EventType {
	switch et {
	case 1:
		return domain.EventTypeSlideCreated
	case 2:
		return domain.EventTypeSlideFinished
	}

	return 0
}
