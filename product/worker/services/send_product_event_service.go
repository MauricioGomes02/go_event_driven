package services

import (
	"context"
	"fmt"
	"go_event_driven/product/domain/ports"
	"time"
)

type SendProductEventService struct {
	databaseAdapter         ports.IProductDatabasePort
	criterionBuilderAdapter ports.CriterionBuilderPort
	eventBusAdapter         ports.IProductEventBusPort
}

func NewSendProductEventService(
	databaseAdapter ports.IProductDatabasePort,
	criterionBuilderAdapter ports.CriterionBuilderPort,
	eventBusAdapter ports.IProductEventBusPort,
) *SendProductEventService {
	return &SendProductEventService{
		databaseAdapter:         databaseAdapter,
		criterionBuilderAdapter: criterionBuilderAdapter,
		eventBusAdapter:         eventBusAdapter,
	}
}

func (service *SendProductEventService) SendCreatedEvent(_context context.Context) {
	fmt.Println("Starting SendProductEventService")

	eventType := "product.created"
	for {
		select {
		case <-_context.Done():
			fmt.Println("The context has been closed")
			return
		default:
			defer func() {
				_panic := recover()
				if _panic != nil {
					fmt.Println(_panic)
				}
			}()

			criterion := service.criterionBuilderAdapter.Or(
				service.criterionBuilderAdapter.Where("OutboxEvent", "Status", "=", "error"),
				service.criterionBuilderAdapter.Where("OutboxEvent", "Status", "=", "pending"),
			)
			_entities, _error := service.databaseAdapter.GetOutboxEvents(criterion)

			if _error != nil {
				fmt.Println(_error.Error())
				continue
			}

			fmt.Println(fmt.Sprintf("Got %d outbox_events", len(_entities)))

			for _, entity := range _entities {
				_error := service.eventBusAdapter.Publish(eventType, entity.Payload)

				if _error != nil {
					fmt.Println(_error.Error())

					entity.UpdateStatus("error")
					entity.UpdateErrorMessage(_error.Error())
					entity.UpdateRetries(entity.Retries + 1)

					_error = service.databaseAdapter.UpdateOutboxEvent(entity)

					if _error != nil {
						fmt.Println(_error.Error())
					}

					fmt.Println("Updated the OutboxEvent entity with error information when trying to publish an event")

					continue
				}

				fmt.Println(fmt.Sprintf("Published the %s event", eventType))

				entity.UpdateStatus("sent")
				entity.UpdateSentAt(time.Now().UTC())

				_error = service.databaseAdapter.UpdateOutboxEvent(entity)

				if _error != nil {
					fmt.Println(_error.Error())
				}

				fmt.Println("Updated the OutboxEvent entity with success information when trying to publish an event")
			}

			time.Sleep(1 * time.Second)
		}
	}
}
