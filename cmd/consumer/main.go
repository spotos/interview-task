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

	"interview/internal/config"
	"interview/internal/container"
	"interview/pkg/queues"
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
			newStopChannel,
			newErrorChannel,
		),
		fx.Populate(&appCtx.ErrChan, &appCtx.Logger, &appCtx.StopChan),
		fx.Invoke(runConsumer),
	)

	if app.Err() != nil {
		log.Fatalf("failed to construct app: %v", app.Err())
	}

	if err := app.Start(context.Background()); err != nil {
		appCtx.Logger.With(zap.Error(err)).Fatal("failed to start consumer")
	}

	appCtx.Logger.Info("consumer started")

	select {
	case err := <-*appCtx.ErrChan:
		appCtx.Logger.With(zap.Error(err)).Error("error while running consumer")
	case <-*appCtx.StopChan:
		appCtx.Logger.Info("stopping consumer")
	}

	if err := app.Stop(context.Background()); err != nil {
		appCtx.Logger.With(zap.Error(err)).Error("failed to stop consumer")
	}
}

func runConsumer(
	lc fx.Lifecycle,
	queuesClient queues.Client,
	logger *zap.Logger,
	errorChannel ErrorChan,
) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				callback := func(callbackCtx context.Context, data []byte) error {
					logger.With(zap.ByteString("payload", data)).Info("consumed message")

					return nil
				}

				*errorChannel <- errors.Wrap(
					queuesClient.Consume(ctx, config.Subscription, callback),
					"consuming events",
				)
			}()

			return nil
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
