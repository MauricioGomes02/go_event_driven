package services

import (
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
	Create(createProduct *inputModels.CreateProduct) (*outputModels.Product, error)
}

type ProductService struct {
	databaseAdapter ports.IProductDatabasePort
}

func NewProductService(databaseAdapter ports.IProductDatabasePort) *ProductService {
	return &ProductService{
		databaseAdapter: databaseAdapter,
	}
}

func (productService *ProductService) Create(createProduct *inputModels.CreateProduct) (*outputModels.Product, error) {
	utcNow := time.Now().UTC()
	newId, _error := uuid.NewUUID()

	if _error != nil {
		return nil, _error
	}

	outboxEventId, _error := uuid.NewUUID()

	if _error != nil {
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
		return nil, _error
	}

	productCreatedEvent := events.NewProductCreatedEvent(entity)
	_bytes, _error := json.Marshal(productCreatedEvent)

	if _error != nil {
		return nil, _error
	}

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

	_entity, _error := productService.databaseAdapter.AddWithOutboxEvent(entity, outboxEvent)

	if _error != nil {
		return nil, _error
	}

	output := outputModels.ConvertFromDomainToApplication(_entity)
	return output, _error
}
