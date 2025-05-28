package outputmodels

import (
	time "time"

	uuid "github.com/google/uuid"
	decimal "github.com/shopspring/decimal"

	outpuModels "go_event_driven/product/application/output_models"
)

type Product struct {
	Id          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	CreatedAt   time.Time       `json:"createdAt"`
}

func ConvertFromApplicationToApi(application *outpuModels.Product) *Product {
	return &Product{
		Id:          application.Id,
		Name:        application.Name,
		Description: application.Description,
		Amount:      application.Amount,
		CreatedAt:   application.CreatedAt,
	}
}
