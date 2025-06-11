package main

import (
	"context"
	"fmt"
	"go_event_driven/product/api/controllers"
	"go_event_driven/product/api/middlewares"
	"go_event_driven/product/api/routes"
	"go_event_driven/product/application/services"
	"go_event_driven/product/configurations"
	"go_event_driven/product/domain/ports"
	"go_event_driven/product/infrastructure/adapters"
	"os"
	"os/signal"

	gin "github.com/gin-gonic/gin"
)

func main() {
	loggerZapAdapter := adapters.SetupLoggerZapAdapter()
	background := context.Background()
	_context, cancel := context.WithCancel(background)
	_context = loggerZapAdapter.With(
		_context,
		ports.Field{Key: "service.name", Value: "app-product-api"})

	defer cancel()

	_signal := make(chan os.Signal, 1)
	signal.Notify(_signal, os.Interrupt)

	go func() {
		<-_signal
		cancel()
	}()

	configuration, _error := configurations.LoadConfigurations("product/api/.env")

	if _error != nil {
		loggerZapAdapter.LogError(
			_context,
			"Error loading settings",
			ports.Field{Key: "error.reason", Value: _error.Error()})
		return
	}

	_context = loggerZapAdapter.With(
		_context,
		ports.Field{Key: "service.environment", Value: configuration.Environment})

	engine := gin.Default()
	engine.Use(middlewares.HandleLoggingMiddleware(_context, loggerZapAdapter))
	engine.Use(gin.Recovery())
	engine.Use(middlewares.HandleErrorMiddleware(loggerZapAdapter))
	engine.SetTrustedProxies(nil)

	mysqlDatabaseAdapter, _error := adapters.SetupMySqlAdapter(&configuration.DatabaseConfiguration)

	if _error != nil {
		loggerZapAdapter.LogError(
			_context,
			"Error connecting to database",
			ports.Field{Key: "error.reason", Value: _error.Error()},
		)
		return
	}

	productDatabaseAdapter := adapters.SetupProductMysqlAdapter(mysqlDatabaseAdapter, loggerZapAdapter)
	productApplicationService := services.NewProductService(productDatabaseAdapter, loggerZapAdapter)
	productController := controllers.NewProductController(productApplicationService, loggerZapAdapter)
	controllers := controllers.Controllers{
		ProductController: *productController,
	}

	routes.ConfigureRoutes(engine, &controllers)

	address := fmt.Sprintf(":%s", configuration.ApiConfiguration.Port)
	engine.Run(address)
}
