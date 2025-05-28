package controllers

import (
	apiInputModels "go_event_driven/product/api/input_models"
	outputmodels "go_event_driven/product/api/output_models"
	applicationServices "go_event_driven/product/application/services"
	http "net/http"

	gin "github.com/gin-gonic/gin"
)

type ProductController struct {
	service applicationServices.IProductService
}

func NewProductController(service applicationServices.IProductService) *ProductController {
	return &ProductController{
		service: service,
	}
}

func (controller *ProductController) CreateProduct(context *gin.Context) {
	var createProductApi *apiInputModels.CreateProduct
	_error := context.ShouldBindJSON(&createProductApi)

	if _error != nil {
		context.JSON(http.StatusBadRequest, "The request body is invalid")
	}

	createProductApplication := apiInputModels.ConvertFromApiToApplication(createProductApi)
	applicationOutput, _error := controller.service.Create(createProductApplication)

	if _error != nil {
		context.JSON(http.StatusInternalServerError, _error.Error())
	}

	apiOutput := outputmodels.ConvertFromApplicationToApi(applicationOutput)
	CreatedAtRoute(apiOutput.Id, apiOutput, context)
}
