package controllers

import (
	"context"
	apiInputModels "go_event_driven/product/api/input_models"
	outputmodels "go_event_driven/product/api/output_models"
	applicationServices "go_event_driven/product/application/services"
	"go_event_driven/product/domain/ports"
	http "net/http"

	gin "github.com/gin-gonic/gin"
)

type ProductController struct {
	service applicationServices.IProductService
	logger  ports.Logger
}

func NewProductController(service applicationServices.IProductService, logger ports.Logger) *ProductController {
	return &ProductController{
		service: service,
		logger:  logger,
	}
}

func (controller *ProductController) CreateProduct(ginContext *gin.Context) {
	getValue, _ := ginContext.Get("contextLogger")
	contextLogger := getValue.(context.Context)

	var createProductApi *apiInputModels.CreateProduct
	_error := ginContext.ShouldBindJSON(&createProductApi)
	controller.logger.LogInformation(
		contextLogger,
		"Converted body request to api input",
		ports.Field{Key: "api.input.type", Value: "CreateProduct"})

	if _error != nil {
		ginContext.JSON(http.StatusBadRequest, "The request body is invalid")
	}

	createProductApplication := apiInputModels.ConvertFromApiToApplication(createProductApi)
	applicationOutput, _error := controller.service.Create(contextLogger, createProductApplication)

	if _error != nil {
		ginContext.JSON(http.StatusInternalServerError, _error.Error())
	}

	apiOutput := outputmodels.ConvertFromApplicationToApi(applicationOutput)
	CreatedAtRoute(apiOutput.Id, apiOutput, ginContext)
}
