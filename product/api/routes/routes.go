package routes

import (
	controllers "go_event_driven/product/api/controllers"

	gin "github.com/gin-gonic/gin"
)

func ConfigureRoutes(context *gin.Engine, controllers *controllers.Controllers) {
	version1 := context.Group("/v1")
	productsV1 := version1.Group("/products")

	productsV1.POST("", controllers.ProductController.CreateProduct)
}
