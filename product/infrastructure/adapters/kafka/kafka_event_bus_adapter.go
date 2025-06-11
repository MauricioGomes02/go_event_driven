package kafka

import (
	"context"
	"fmt"
	"go_event_driven/product/configurations"
	"go_event_driven/product/domain/ports"

	"github.com/IBM/sarama"
)

type KafkaEventBusAdapter struct {
	broker     string
	producer   sarama.SyncProducer
	eventTopic map[string]string
	logger     ports.Logger
}

func NewKafkaEventBusAdapter(kafkaConfiguration *configurations.KafkaConfiguration, logger ports.Logger) (*KafkaEventBusAdapter, error) {
	_producer, _error := initProducer(kafkaConfiguration.Broker)
	if _error != nil {
		return nil, _error
	}

	eventTopic := make(map[string]string)

	eventTopic["product.created"] = kafkaConfiguration.KafkaTopicConfiguration.ProductCreated

	return &KafkaEventBusAdapter{
		broker:     kafkaConfiguration.Broker,
		producer:   _producer,
		eventTopic: eventTopic,
		logger:     logger,
	}, nil
}

func initProducer(broker string) (sarama.SyncProducer, error) {
	settings := sarama.NewConfig()
	settings.Producer.RequiredAcks = sarama.WaitForAll
	settings.Producer.Return.Successes = true

	brokers := []string{broker}
	_producer, _error := sarama.NewSyncProducer(brokers, settings)
	if _error != nil {
		return nil, _error
	}

	return _producer, nil
}

func (adapter *KafkaEventBusAdapter) Publish(_context context.Context, eventType string, message []byte) error {
	topic := adapter.eventTopic[eventType]

	_context = adapter.logger.With(
		_context,
		ports.Field{Key: "event.name", Value: eventType},
		ports.Field{Key: "event.topic.name", Value: topic},
	)

	if topic == "" {
		_error := fmt.Errorf("Unable to get topic for event %s", eventType)
		adapter.logger.LogError(
			_context,
			"Error getting topic",
			ports.Field{Key: "error.reason", Value: _error.Error()},
		)

		return _error
	}

	_message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, _error := adapter.producer.SendMessage(_message)

	if _error != nil {
		adapter.logger.LogError(
			_context,
			"Error publishing event",
		)

		return _error
	}

	adapter.logger.LogInformation(
		_context,
		"Published event",
		ports.Field{Key: "event.topic.partition", Value: partition},
		ports.Field{Key: "event.topic.offset", Value: offset},
	)

	return nil
}
