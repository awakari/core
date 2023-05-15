package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Uri struct {
		Messages      string `envconfig:"URI_MESSAGES" default:"localhost:50051" required:"true"`
		Subscriptions string `envconfig:"URI_SUBSCRIPTIONS" default:"localhost:50052" required:"true"`
		Writer        string `envconfig:"URI_WRITER" default:"localhost:50053" required:"true"`
	}
}

func NewConfigFromEnv() (cfg Config, err error) {
	err = envconfig.Process("", &cfg)
	return
}
