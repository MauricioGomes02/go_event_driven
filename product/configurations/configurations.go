package configurations

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Environment           string `envconfig:"ENVIRONMENT"`
	ApiConfiguration      ApiConfiguration
	DatabaseConfiguration DatabaseConfiguration
	KafkaConfiguration    KafkaConfiguration
}

type ApiConfiguration struct {
	Port string `envconfig:"API_PORT"`
}

type DatabaseConfiguration struct {
	Host     string `envconfig:"DATABASE_HOST"`
	Port     int    `envconfig:"DATABASE_PORT"`
	UserName string `envconfig:"DATABASE_USERNAME"`
	Password string `envconfig:"DATABASE_PASSWORD"`
	Name     string `envconfig:"DATABASE_NAME"`
}

type KafkaConfiguration struct {
	Broker                  string `envconfig:"KAFKA_BROKER"`
	KafkaTopicConfiguration KafkaTopicConfiguration
}

type KafkaTopicConfiguration struct {
	ProductCreated string `envconfig:"KAFKA_TOPIC_PRODUCT_CREATED"`
}

func LoadConfigurations(envFilePath string) (*Configuration, error) {
	_ = godotenv.Load(envFilePath)

	var configuration Configuration
	_error := envconfig.Process("", &configuration)
	if _error != nil {
		fmt.Printf("An error occurred while loading settings: %s", _error.Error())
		return nil, _error
	}

	return &configuration, nil
}
