package configurations

type Configuration struct {
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
