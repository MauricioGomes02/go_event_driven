package configurations

type Configuration struct {
	ApiConfiguration      ApiConfiguration
	DatabaseConfiguration DatabaseConfiguration
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
