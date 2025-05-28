package inputmodels

import (
	decimal "github.com/shopspring/decimal"
)

type CreateProduct struct {
	Name        string
	Description string
	Amount      decimal.Decimal
}
