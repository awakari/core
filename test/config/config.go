package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Uri struct {
		Reader        string `envconfig:"URI_READER" default:"localhost:50051" required:"true"`
		Resolver      string `envconfig:"URI_RESOLVER" default:"localhost:50052" required:"true"`
		Subscriptions string `envconfig:"URI_SUBSCRIPTIONS" default:"localhost:50053" required:"true"`
	}
}

func NewConfigFromEnv() (cfg Config, err error) {
	err = envconfig.Process("", &cfg)
	return
}
