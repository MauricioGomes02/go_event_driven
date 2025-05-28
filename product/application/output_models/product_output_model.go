package outputmodels

import (
	entities "go_event_driven/product/domain/entities"
	time "time"

	uuid "github.com/google/uuid"
	decimal "github.com/shopspring/decimal"
)

type Product struct {
	Id          uuid.UUID
	Name        string
	Description string
	Amount      decimal.Decimal
	CreatedAt   time.Time
}

func ConvertFromDomainToApplication(entity *entities.Product) *Product {
	if entity == nil {
		return nil
	}

	return &Product{
		Id:          entity.Id,
		Name:        entity.Name,
		Description: entity.Description,
		Amount:      entity.Amount,
		CreatedAt:   entity.CreatedAt,
	}
}
