package application

import (	
	"tasks/pkg/infra/logger"
	"tasks/pkg/infra/metrics"
	"tasks/pkg/infra/probes"
	"tasks/internal/domain/usecases"
	"tasks/internal/adapters/http"
	"context"
)

type App struct{
	l logger.Logger
	shutdownFuncs []func(ctx context.Context) error
}

func New(l logger.Logger) *App{
	return &App{
		l:l,
	}
}

func (app *App) Start() error {
	probes, _ := probes.New(app.l)
	probes.SetStarted()
	err := probes.Start()
	if err != nil {
		app.l.Sugar().Fatalf("probes start failed: %s", err.Error())
	}

	me, err := metrics.New()
	if err != nil {
		app.l.Sugar().Fatalf("Metrics init failed: %s", err.Error())
	}
	app.shutdownFuncs = append(app.shutdownFuncs, me.Stop)
	
	storage := usecases.NewStorage()

	tasks, err := usecases.New(app.l, storage)
	if err != nil {
		app.l.Sugar().Fatalf("create buissness logic failed: %s", err.Error())
	}

	s, err := http.New(tasks, app.l)
	if err != nil {
		app.l.Sugar().Fatalf("server not started %s", err.Error())
	}
	app.shutdownFuncs = append(app.shutdownFuncs, s.Stop)
	err = s.Start()
	if err != nil {
		app.l.Sugar().Fatalf("server not started: %s", err.Error())
	}

	probes.SetReady()

	return nil

} 

func (a *App) Stop(ctx context.Context) error {
	var err error
	for i := len(a.shutdownFuncs)-1; i>=0; i-- {
		err  = a.shutdownFuncs[i](ctx)
		if  err != nil {
			a.l.Sugar().Error(err)
		}
	}

	a.l.Info("app stopped")

	return nil

} 