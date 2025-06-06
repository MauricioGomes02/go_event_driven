package kafka

import (
	"fmt"
	"go_event_driven/product/configurations"

	"github.com/IBM/sarama"
)

type KafkaEventBusAdapter struct {
	broker     string
	producer   sarama.SyncProducer
	eventTopic map[string]string
}

func NewKafkaEventBusAdapter(kafkaConfiguration *configurations.KafkaConfiguration) (*KafkaEventBusAdapter, error) {
	_producer, _error := initProducer(kafkaConfiguration.Broker)
	if _error != nil {
		// TODO
	}

	eventTopic := make(map[string]string)

	eventTopic["product.created"] = kafkaConfiguration.KafkaTopicConfiguration.ProductCreated

	return &KafkaEventBusAdapter{
		broker:     kafkaConfiguration.Broker,
		producer:   _producer,
		eventTopic: eventTopic,
	}, nil
}

func initProducer(broker string) (sarama.SyncProducer, error) {
	settings := sarama.NewConfig()
	settings.Producer.RequiredAcks = sarama.WaitForAll
	settings.Producer.Return.Successes = true

	brokers := []string{broker}
	_producer, _error := sarama.NewSyncProducer(brokers, settings)
	if _error != nil {
		// TODO
		fmt.Println(_error.Error())
	}

	return _producer, nil
}

func (adapter *KafkaEventBusAdapter) Publish(eventType string, message []byte) error {
	topic := adapter.eventTopic[eventType]

	fmt.Println(fmt.Sprintf("Got topic %s", topic))

	_message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, _error := adapter.producer.SendMessage(_message)
	return _error
}
