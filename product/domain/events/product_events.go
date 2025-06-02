package events

import (
	"go_event_driven/product/domain/entities"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductCreatedEvent struct {
	EventId        uuid.UUID               `json:"event_id"`
	EventType      string                  `json:"event_type"`
	EventTimestamp time.Time               `json:"event_timestamp"`
	Data           ProductCreatedEventData `json:"data"`
}

type ProductCreatedEventData struct {
	Id          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	CreatedAt   time.Time       `json:"created_at"`
}

func NewProductCreatedEvent(id uuid.UUID, eventType string, timestamp time.Time, entity *entities.Product) *ProductCreatedEvent {
	return &ProductCreatedEvent{
		EventId:        id,
		EventType:      eventType,
		EventTimestamp: timestamp,
		Data: ProductCreatedEventData{
			Id:          entity.Id,
			Name:        entity.Name,
			Description: entity.Description,
			Amount:      entity.Amount,
			CreatedAt:   entity.CreatedAt,
		},
	}
}
