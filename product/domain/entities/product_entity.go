package entities

import (
	"errors"
	commonerrors "go_event_driven/product/domain/errors"
	"time"

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
) (*Product, error) {
	var _errors []error

	if name == "" {
		_errors = append(_errors, commonerrors.NewInvalidEntityError("Product name must be filled in"))
	}

	if description == "" {
		_errors = append(_errors, commonerrors.NewInvalidEntityError("Product description must be filled in"))
	}

	if amount.LessThanOrEqual(decimal.New(0, 0)) {
		_errors = append(_errors, commonerrors.NewInvalidEntityError("Product amount cannot be less than or equal to 0"))
	}

	if len(_errors) > 0 {
		return nil, errors.Join(_errors...)
	}

	return &Product{
		Id:          id,
		Name:        name,
		Description: description,
		Amount:      amount,
		CreatedAt:   createdAt,
		Active:      active,
	}, nil
}
