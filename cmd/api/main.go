package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"interview/internal/api/hello"
	"interview/pkg/echo"
	"interview/pkg/zapwrapper"

	"interview/internal/api"
	"interview/internal/config"
	"interview/internal/container"
)

type (
	ErrorChan *chan error
	StopChan  *chan os.Signal

	appContext struct {
		Logger   *zap.Logger
		ErrChan  ErrorChan
		StopChan StopChan
	}
)

func main() {
	var appCtx appContext

	app := fx.New(
		fx.NopLogger,
		container.Module,
		fx.Provide(
			api.NewRouteConfigurator,
			newStopChannel,
			newErrorChannel,
			hello.NewController,
			hello.NewRouter,
		),
		fx.Populate(&appCtx.ErrChan, &appCtx.StopChan, &appCtx.Logger),
		fx.Invoke(runApi),
	)

	if app.Err() != nil {
		log.Fatalf("failed to construct app: %v", app.Err())
	}

	if err := app.Start(context.Background()); err != nil {
		appCtx.Logger.With(zap.Error(err)).Fatal("failed to start app")
	}

	select {
	case err := <-*appCtx.ErrChan:
		appCtx.Logger.With(zap.Error(err)).Error("error while running")
	case <-*appCtx.StopChan:
		appCtx.Logger.Info("stopping application")
	}

	if err := app.Stop(context.Background()); err != nil {
		appCtx.Logger.With(zap.Error(err)).Error("failed to stop app")
	}
}

func runApi(
	lc fx.Lifecycle,
	cfg *config.Config,
	routeConfigurator echo.RouteConfigurator,
	errorChannel ErrorChan,
	logger *zap.Logger,
) error {
	engine, err := echo.New(&echo.Config{
		UnixSocket: cfg.UnixSocket,
		Router:     routeConfigurator,
		Logger:     zapwrapper.EchoLogger(logger),
		HideBanner: true,
	})
	if err != nil {
		return errors.Wrap(err, "init api engine")
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				*errorChannel <- errors.Wrap(engine.Run(), "run api engine")
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return errors.Wrap(engine.Stop(ctx), "stop engine")
		},
	})

	return nil
}
func newStopChannel() StopChan {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	return &stop
}

func newErrorChannel() ErrorChan {
	errChan := make(chan error, 1)

	return &errChan
}
