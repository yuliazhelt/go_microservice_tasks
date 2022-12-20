package http

import (
	"tasks/internal/ports"
	"tasks/internal/domain/models"
	"tasks/pkg/infra/logger"
	"tasks/internal/adapters/auth"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	ginzap "github.com/gin-contrib/zap"
	"github.com/juju/zaputil/zapctx"
)

type Adapter struct {
	s    *http.Server
	l    net.Listener
	tasks ports.Tasks
}

type Config struct {
	Port int `env:"HTTP_PORT" envDefault:"3000"`
}

func AuthMiddleware(a ports.AuthAdapter) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		err := a.Verify(ctx)
		if err != nil {
			if errors.Is(err, models.ErrForbidden) {
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			} else {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
		ctx.Next()
	}
	return fn
}

func New(tasks ports.Tasks, log logger.Logger) (*Adapter, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse server http adapter configuration failed: %w", err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("server start failed: %w", err)
	}

	authAdapter := auth.New()

	router := gin.Default()

	server := http.Server{
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	a := Adapter{
		s:    &server,
		l:    l,
		tasks: tasks,
	}

	router.Use(func(ctx *gin.Context) {
		lCtx := zapctx.WithLogger(ctx.Request.Context(), log)
		ctx.Request = ctx.Request.WithContext(lCtx)
	})
	router.Use(ginzap.Ginzap(log, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(log, true))

	v1 := router.Group("/tasks/api/v1/tasks")
	{
		v1.POST("/approve/:taskId/:approveInd", a.approve)
		v1.POST("/decline/:taskId/:approveInd", a.decline)
		v1.GET("/:taskId", a.getTaskID) // по id таски
	}
	v2 := router.Group("/tasks/api/v1/tasks").Use(AuthMiddleware(authAdapter))
	v2.POST("/", a.create)
	v2.GET("/", a.getUserTasks)

	return &a, nil
}

func (a *Adapter) Start() error {
	var err error
	go func() {
		err = a.s.Serve(a.l)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			err = fmt.Errorf("server start failed: %w", err)
		}
		err = nil
	}()

	if err != nil {
		return err
	}
	return nil
}

func (a *Adapter) Stop(ctx context.Context) error {
	var (
		err  error
		once sync.Once
	)
	once.Do(func() {
		err = a.s.Shutdown(ctx)
	})
	return err
}
