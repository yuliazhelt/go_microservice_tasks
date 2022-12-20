package logger

import (
	"fmt"

	"github.com/caarlos0/env"
	"go.uber.org/zap"
)

type Logger = *zap.Logger

type Config struct {
	IsProd bool `env:"IS_PROD" envDefault:"false"`
	SentryAdress   string `env:"SENTRY_ADDRESS" envDefault:"https://a9df2a83e2e6421abd2301d28d0db7a8@o4503908268507136.ingest.sentry.io/4503953648648192"`
}

func New() (Logger, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("read logger configuration failed: %w", err)
	}

	var zapCfg zap.Config
	if cfg.IsProd {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}

	l, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("create logger failed: %w", err)
	}

	zap.ReplaceGlobals(l)

	return l, nil
}
