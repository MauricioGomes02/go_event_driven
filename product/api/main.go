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
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	_ = godotenv.Load("product/api/.env")

	var configuration configurations.Configuration
	_error := envconfig.Process("", &configuration)
	if _error != nil {
		return
	}

	engine := gin.Default()
	engine.Use(gin.Recovery())
	engine.Use(middlewares.HandleErrorMiddleware())
	engine.SetTrustedProxies(nil)

	mysqlDatabaseAdapter := adapters.SetupMySqlAdapter(&configuration.DatabaseConfiguration)
	kafkaEventBusAdapter := adapters.SetupKafkaAdapter(&configuration.KafkaConfiguration)
	productDatabaseAdapter := adapters.SetupProductMysqlAdapter(mysqlDatabaseAdapter)
	productApplicationService := services.NewProductService(productDatabaseAdapter, kafkaEventBusAdapter)
	productController := controllers.NewProductController(productApplicationService)
	controllers := controllers.Controllers{
		ProductController: *productController,
	}

	routes.ConfigureRoutes(engine, &controllers)

	address := fmt.Sprintf(":%s", configuration.ApiConfiguration.Port)
	engine.Run(address)
}
