package main

import (
	"context"
	"go_event_driven/product/configurations"
	"go_event_driven/product/domain/ports"
	"go_event_driven/product/infrastructure/adapters"
	"go_event_driven/product/worker/services"
	"os"
	"os/signal"
)

func main() {
	loggerZapAdapter := adapters.SetupLoggerZapAdapter()
	background := context.Background()
	_context, cancel := context.WithCancel(background)

	defer cancel()

	_signal := make(chan os.Signal, 1)
	signal.Notify(_signal, os.Interrupt)

	go func() {
		<-_signal
		cancel()
	}()

	configuration, _error := configurations.LoadConfigurations("product/worker/.env")

	_context = loggerZapAdapter.With(
		_context,
		ports.Field{Key: "service.name", Value: "app-product-worker"})

	if _error != nil {
		loggerZapAdapter.LogError(_context, "Error loading settings")
		return
	}

	_context = loggerZapAdapter.With(
		_context,
		ports.Field{Key: "service.environment", Value: configuration.Environment})

	mysqlDatabaseAdapter, _error := adapters.SetupMySqlAdapter(&configuration.DatabaseConfiguration)

	if _error != nil {
		loggerZapAdapter.LogError(
			_context,
			"Error connecting to database",
			ports.Field{Key: "error.reason", Value: _error.Error()},
		)
	}

	productDatabaseAdapter := adapters.SetupProductMysqlAdapter(mysqlDatabaseAdapter, loggerZapAdapter)
	mysqlCriterionBuilderAdapter := adapters.SetupCriterionBuilderMysqlAdapter()
	kafkaEventBusAdapter, _error := adapters.SetupKafkaAdapter(&configuration.KafkaConfiguration, loggerZapAdapter)

	if _error != nil {
		loggerZapAdapter.LogError(
			_context,
			"Error trying to connect to the database",
			ports.Field{Key: "error.reason", Value: _error.Error()},
		)
	}

	sendProductEventService := services.NewSendProductEventService(productDatabaseAdapter, mysqlCriterionBuilderAdapter, kafkaEventBusAdapter, loggerZapAdapter)

	sendProductEventService.SendCreatedEvent(_context, loggerZapAdapter)
}
