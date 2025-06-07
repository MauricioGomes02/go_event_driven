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
	_context = loggerZapAdapter.With(
		_context,
		ports.Field{Key: "worker.domain", Value: "product"})

	defer cancel()

	_signal := make(chan os.Signal, 1)
	signal.Notify(_signal, os.Interrupt)

	go func() {
		<-_signal
		cancel()
	}()

	configuration, _error := configurations.LoadConfigurations("product/worker/.env")

	if _error != nil {
		loggerZapAdapter.LogError(_context, "Error loading settings")
		return
	}

	_context = loggerZapAdapter.With(
		_context,
		ports.Field{Key: "worker.environment", Value: configuration.Environment})

	mysqlDatabaseAdapter := adapters.SetupMySqlAdapter(&configuration.DatabaseConfiguration)
	productDatabaseAdapter := adapters.SetupProductMysqlAdapter(mysqlDatabaseAdapter)
	mysqlCriterionBuilderAdapter := adapters.SetupCriterionBuilderMysqlAdapter()
	kafkaEventBusAdapter := adapters.SetupKafkaAdapter(&configuration.KafkaConfiguration)
	sendProductEventService := services.NewSendProductEventService(productDatabaseAdapter, mysqlCriterionBuilderAdapter, kafkaEventBusAdapter)

	sendProductEventService.SendCreatedEvent(_context, loggerZapAdapter)
}
