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
	eventBusAdapter ports.IProductEventBusPort
}

func NewProductService(
	databaseAdapter ports.IProductDatabasePort,
	eventBusAdapter ports.IProductEventBusPort) *ProductService {
	return &ProductService{
		databaseAdapter: databaseAdapter,
		eventBusAdapter: eventBusAdapter,
	}
}

func (productService *ProductService) Create(createProduct *inputModels.CreateProduct) (*outputModels.Product, error) {
	utcNow := time.Now().UTC()
	newId, _error := uuid.NewUUID()

	if _error != nil {
		return nil, _error
	}

	entity := entities.NewProduct(
		newId,
		createProduct.Name,
		createProduct.Description,
		createProduct.Amount,
		utcNow,
		true,
	)

	_entity, _error := productService.databaseAdapter.Add(entity)

	eventId, _error := uuid.NewUUID()

	if _error != nil {
		// TODO
	}

	eventType := "product.created"
	timestamp := time.Now().UTC()
	productCreatedEvent := events.NewProductCreatedEvent(eventId, eventType, timestamp, _entity)
	_bytes, _error := json.Marshal(productCreatedEvent)

	if _error != nil {
		// TODO
	}

	productService.eventBusAdapter.Publish(eventType, _bytes)
	output := outputModels.ConvertFromDomainToApplication(_entity)
	return output, _error
}
