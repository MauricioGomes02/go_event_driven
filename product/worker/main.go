package main

import (
	"context"
	"fmt"
	"go_event_driven/product/configurations"
	"go_event_driven/product/infrastructure/adapters"
	"go_event_driven/product/worker/services"
	"os"
	"os/signal"
)

func main() {
	configuration, _error := configurations.LoadConfigurations("product/worker/.env")

	if _error != nil {
		fmt.Println("Unable to load settings")
		return
	}

	background := context.Background()
	_context, cancel := context.WithCancel(background)

	defer cancel()

	_signal := make(chan os.Signal, 1)
	signal.Notify(_signal, os.Interrupt)

	go func() {
		<-_signal
		cancel()
	}()

	mysqlDatabaseAdapter := adapters.SetupMySqlAdapter(&configuration.DatabaseConfiguration)
	productDatabaseAdapter := adapters.SetupProductMysqlAdapter(mysqlDatabaseAdapter)
	mysqlCriterionBuilderAdapter := adapters.SetupCriterionBuilderMysqlAdapter()
	kafkaEventBusAdapter := adapters.SetupKafkaAdapter(&configuration.KafkaConfiguration)
	sendProductEventService := services.NewSendProductEventService(productDatabaseAdapter, mysqlCriterionBuilderAdapter, kafkaEventBusAdapter)

	sendProductEventService.SendCreatedEvent(_context)
}
