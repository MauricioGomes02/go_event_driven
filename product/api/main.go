package main

import (
	"fmt"
	"go_event_driven/product/api/controllers"
	"go_event_driven/product/api/middlewares"
	"go_event_driven/product/api/routes"
	"go_event_driven/product/application/services"
	"go_event_driven/product/configurations"
	"go_event_driven/product/infrastructure/adapters"

	gin "github.com/gin-gonic/gin"
)

func main() {
	configuration, _error := configurations.LoadConfigurations("product/api/.env")

	if _error != nil {
		fmt.Println("Unable to load settings")
		return
	}

	engine := gin.Default()
	engine.Use(gin.Recovery())
	engine.Use(middlewares.HandleErrorMiddleware())
	engine.SetTrustedProxies(nil)

	mysqlDatabaseAdapter := adapters.SetupMySqlAdapter(&configuration.DatabaseConfiguration)
	productDatabaseAdapter := adapters.SetupProductMysqlAdapter(mysqlDatabaseAdapter)
	productApplicationService := services.NewProductService(productDatabaseAdapter)
	productController := controllers.NewProductController(productApplicationService)
	controllers := controllers.Controllers{
		ProductController: *productController,
	}

	routes.ConfigureRoutes(engine, &controllers)

	address := fmt.Sprintf(":%s", configuration.ApiConfiguration.Port)
	engine.Run(address)
}
