package services

import (
	"context"
	"encoding/json"
	inputModels "go_event_driven/product/application/input_models"
	outputModels "go_event_driven/product/application/output_models"
	entities "go_event_driven/product/domain/entities"
	"go_event_driven/product/domain/events"
	ports "go_event_driven/product/domain/ports"
	time "time"

	"github.com/google/uuid"
)

type IProductService interface {
	Create(_context context.Context, createProduct *inputModels.CreateProduct) (*outputModels.Product, error)
}

type ProductService struct {
	databaseAdapter ports.IProductDatabasePort
	logger          ports.Logger
}

func NewProductService(databaseAdapter ports.IProductDatabasePort, logger ports.Logger) *ProductService {
	return &ProductService{
		databaseAdapter: databaseAdapter,
		logger:          logger,
	}
}

func (service *ProductService) Create(_context context.Context, createProduct *inputModels.CreateProduct) (*outputModels.Product, error) {
	utcNow := time.Now().UTC()
	newId, _error := uuid.NewUUID()

	if _error != nil {
		service.logger.LogError(
			_context,
			"Error generating new id",
			ports.Field{Key: "error.reason", Value: _error.Error()},
		)
		return nil, _error
	}

	outboxEventId, _error := uuid.NewUUID()

	if _error != nil {
		service.logger.LogError(
			_context,
			"Error generating new id",
			ports.Field{Key: "error.reason", Value: _error.Error()},
		)
		return nil, _error
	}

	entity, _error := entities.NewProduct(
		newId,
		createProduct.Name,
		createProduct.Description,
		createProduct.Amount,
		utcNow,
		true,
	)

	if _error != nil {
		service.logger.LogError(
			_context,
			"Error generating new entity",
			ports.Field{Key: "error.reason", Value: _error.Error()},
			ports.Field{Key: "entity.name", Value: "Product"},
		)
		return nil, _error
	}

	productCreatedEvent := events.NewProductCreatedEvent(entity)
	_bytes, _ := json.Marshal(productCreatedEvent)

	outboxEvent := entities.NewOutboxEvent(
		outboxEventId,
		entity.Id,
		"product",
		"product.created",
		_bytes,
		"pending",
		0,
		nil,
		utcNow,
		nil,
	)

	_context = service.logger.With(
		_context,
		ports.Field{Key: "event.entity", Value: "product"},
		ports.Field{Key: "event.type", Value: "created"},
	)

	_entity, _error := service.databaseAdapter.AddWithOutboxEvent(entity, outboxEvent)

	if _error != nil {
		service.logger.LogError(
			_context,
			"Error persisting entity with event",
			ports.Field{Key: "error.reason", Value: _error.Error()},
		)
		return nil, _error
	}

	service.logger.LogInformation(
		_context,
		"Persisted entity and event",
	)

	output := outputModels.ConvertFromDomainToApplication(_entity)
	return output, nil
}
