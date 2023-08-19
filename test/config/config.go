package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	Uri struct {
		Reader        string `envconfig:"URI_READER" default:"localhost:50051" required:"true"`
		Resolver      string `envconfig:"URI_RESOLVER" default:"localhost:50052" required:"true"`
		Subscriptions string `envconfig:"URI_SUBSCRIPTIONS" default:"localhost:50053" required:"true"`
		Matches       string `envconfig:"URI_MATCHES" default:"localhost:50054" required:"true"`
		Messages      string `envconfig:"URI_MESSAGES" default:"localhost:50055" required:"true"`
	}
	Test struct {
		Perf struct {
			E2e struct {
				SubCount  int           `envconfig:"TEST_PERF_E2E_SUB_COUNT" default:"1" required:"true"`
				WriteRate float64       `envconfig:"TEST_PERF_E2E_WRITE_RATE" default:"5" required:"true"`
				BatchSize int           `envconfig:"TEST_PERF_E2E_BATCH_SIZE" default:"16"`
				Duration  time.Duration `envconfig:"TEST_PERF_E2E_DURATION" default:"200s"`
			}
		}
	}
}

func NewConfigFromEnv() (cfg Config, err error) {
	err = envconfig.Process("", &cfg)
	return
}
