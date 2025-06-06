package events

import (
	"go_event_driven/product/domain/entities"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductCreatedEvent struct {
	Id          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	CreatedAt   time.Time       `json:"created_at"`
}

func NewProductCreatedEvent(entity *entities.Product) *ProductCreatedEvent {
	return &ProductCreatedEvent{
		Id:          entity.Id,
		Name:        entity.Name,
		Description: entity.Description,
		Amount:      entity.Amount,
		CreatedAt:   entity.CreatedAt,
	}
}
