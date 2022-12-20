package probes

import (
	"tasks/pkg/infra/logger"
	"fmt"
	"net/http"
	"sync"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
)

type Probes struct {
	isReady   bool
	isStarted bool

	readyOnce   sync.Once
	startedOnce sync.Once

	l logger.Logger
}

type Config struct {
	Port int `env:"PROBES_PORT" envDefault:"3030"`
}

func New(l logger.Logger) (*Probes, error) {
	return &Probes{
		l: l,
	}, nil
}

func (p *Probes) Start() error {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("configuration parsing failed: %w", err)
	}

	r := gin.Default()

	r.GET("/ready", func(ctx *gin.Context) {
		if p.isReady {
			ctx.Writer.WriteHeader(http.StatusOK)
		} else {
			ctx.Writer.WriteHeader(http.StatusLocked)
		}
	})

	go func() {
		err := r.Run(fmt.Sprintf(":%d", cfg.Port))
		if err != nil {
			p.l.Sugar().Errorf("start probes failed: %s", err.Error())
		}
	}()

	return nil
}

func (p *Probes) SetReady() {
	p.readyOnce.Do(func() {
		p.isReady = true
	})
}

func (p *Probes) SetStarted() {
	p.startedOnce.Do(func() {
		p.isStarted = true
	})
}
