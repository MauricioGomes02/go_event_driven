package inputmodels

import (
	applicationInputModels "go_event_driven/product/application/input_models"

	decimal "github.com/shopspring/decimal"
)

type CreateProduct struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
}

func ConvertFromApiToApplication(api *CreateProduct) *applicationInputModels.CreateProduct {
	return &applicationInputModels.CreateProduct{
		Name:        api.Name,
		Description: api.Description,
		Amount:      api.Amount,
	}
}
