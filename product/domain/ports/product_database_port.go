package ports

import (
	entities "go_event_driven/product/domain/entities"
)

type IProductDatabasePort interface {
	AddWithOutboxEvent(product *entities.Product, event *entities.OutboxEvent) (*entities.Product, error)
	GetOutboxEvents(criterion CriterionPort) ([]entities.OutboxEvent, error)
	UpdateOutboxEvent(entity entities.OutboxEvent) error
}
