package services

import (
	inputModels "go_event_driven/product/application/input_models"
	outputModels "go_event_driven/product/application/output_models"
	entities "go_event_driven/product/domain/entities"
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

	entity := entities.NewProduct(
		newId,
		createProduct.Name,
		createProduct.Description,
		createProduct.Amount,
		utcNow,
		true,
	)

	_entity, _error := productService.databaseAdapter.Add(entity)

	output := outputModels.ConvertFromDomainToApplication(_entity)
	return output, _error
}
