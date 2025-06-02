package ports

import (
	entities "go_event_driven/product/domain/entities"
)

type IProductDatabasePort interface {
	Add(product *entities.Product) (*entities.Product, error)
}
