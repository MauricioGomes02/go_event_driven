package entities

import (
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
	Active      bool
}

func NewProduct(
	id uuid.UUID,
	name string,
	description string,
	amount decimal.Decimal,
	createdAt time.Time,
	active bool,
) *Product {
	// validations ...

	return &Product{
		Id:          id,
		Name:        name,
		Description: description,
		Amount:      amount,
		CreatedAt:   createdAt,
		Active:      active,
	}
}
