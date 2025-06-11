package services

import (
	"context"
	"encoding/json"
	"go_event_driven/product/domain/ports"
	"time"
)

type SendProductEventService struct {
	databaseAdapter         ports.IProductDatabasePort
	criterionBuilderAdapter ports.CriterionBuilderPort
	eventBusAdapter         ports.IProductEventBusPort
	logger                  ports.Logger
}

func NewSendProductEventService(
	databaseAdapter ports.IProductDatabasePort,
	criterionBuilderAdapter ports.CriterionBuilderPort,
	eventBusAdapter ports.IProductEventBusPort,
	logger ports.Logger,
) *SendProductEventService {
	return &SendProductEventService{
		databaseAdapter:         databaseAdapter,
		criterionBuilderAdapter: criterionBuilderAdapter,
		eventBusAdapter:         eventBusAdapter,
		logger:                  logger,
	}
}

func (service *SendProductEventService) SendCreatedEvent(_context context.Context, logger ports.Logger) {
	eventType := "product.created"
	_context = logger.With(
		_context,
		ports.Field{Key: "event.entity", Value: "product"},
		ports.Field{Key: "event.type", Value: "created"},
	)

	for {
		select {
		case <-_context.Done():
			logger.LogWarning(_context, "The context has been closed")
			return
		default:
			defer func() {
				_panic := recover()
				if _panic != nil {
					logger.LogError(
						_context,
						"Error sending event",
						ports.Field{Key: "panic.reason", Value: _panic.(string)})
				}
			}()

			criterion := service.criterionBuilderAdapter.Or(
				service.criterionBuilderAdapter.Where("OutboxEvent", "Status", "=", "error"),
				service.criterionBuilderAdapter.Where("OutboxEvent", "Status", "=", "pending"),
			)

			_entities, _error := service.databaseAdapter.GetOutboxEvents(criterion)

			if _error != nil {
				logger.LogError(_context, _error.Error())
				continue
			}

			logger.LogInformation(
				_context,
				"Obtained entities",
				ports.Field{Key: "entity", Value: "OutboxEvent"},
				ports.Field{Key: "entity.quantity", Value: len(_entities)})

			for _, entity := range _entities {
				_error := service.eventBusAdapter.Publish(_context, eventType, entity.Payload)

				if _error != nil {
					logger.LogError(
						_context,
						"Error publishing event",
						ports.Field{Key: "error.reason", Value: _error.Error()})

					entity.UpdateStatus("error")
					entity.UpdateErrorMessage(_error.Error())
					entity.UpdateRetries(entity.Retries + 1)

					updatedFields := entity.GetUpdatedFields()

					_error = service.databaseAdapter.UpdateOutboxEvent(entity)

					if _error != nil {
						logger.LogError(
							_context,
							"Error updating entity",
							ports.Field{Key: "entity", Value: "OutboxEvent"},
							ports.Field{
								Key:   "entity.properties.updated",
								Value: updatedFields,
							})
					}

					logger.LogInformation(
						_context,
						"Updated the entity",
						ports.Field{Key: "entity", Value: "OutboxEvent"},
						ports.Field{
							Key:   "entity.properties.updated",
							Value: updatedFields,
						})

					continue
				}

				logger.LogInformation(_context, "Published the event")

				entity.UpdateStatus("sent")
				entity.UpdateSentAt(time.Now().UTC())

				updatedProperties := entity.GetUpdatedFields()
				serializedUpdatedProperties, _ := json.Marshal(updatedProperties)

				_error = service.databaseAdapter.UpdateOutboxEvent(entity)

				if _error != nil {
					logger.LogError(
						_context,
						"Error updating entity",
						ports.Field{Key: "entity", Value: "OutboxEvent"},
						ports.Field{
							Key:   "entity.properties.updated",
							Value: serializedUpdatedProperties,
						})
				}

				logger.LogInformation(
					_context,
					"Updated the entity",
					ports.Field{Key: "entity", Value: "OutboxEvent"},
					ports.Field{
						Key:   "entity.properties.updated",
						Value: serializedUpdatedProperties,
					})
			}

			time.Sleep(1 * time.Second)
		}
	}
}
