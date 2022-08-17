package container

import (
	"context"
	"flag"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"interview/internal/config"
	"interview/pkg/dotenv"
	"interview/pkg/queues"
	"interview/pkg/zapwrapper"
)

var (
	defaultEnvFile = ".env"
	envFile        = flag.String("env-file", "", "Path to env file")
)

var Module = fx.Options(
	fx.Provide(
		newConfiguration,
		zapwrapper.New,
		newQueuesClient,
		newRedisClient,
	),
)

func newConfiguration() (*config.Config, error) {
	var cfg config.Config

	_, err := conf.Parse("", &cfg, dotenv.FromEnvFiles(defaultEnvFile, *envFile))

	return &cfg, err
}

func newQueuesClient(cfg *config.Config, lc fx.Lifecycle, logger *zap.Logger) (queues.Client, error) {
	client, err := queues.NewClient(cfg.ProjectID, config.Subscriptions, logger)
	if err != nil {
		return nil, errors.Wrap(err, "initialize queues client")
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return errors.Wrap(client.Close(), "close queues client")
		},
	})

	return client, nil
}

func newRedisClient(cfg *config.Config, lc fx.Lifecycle) redis.UniversalClient {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{cfg.Redis.URL},
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return errors.Wrap(
				client.Set(ctx, config.GreetingKey, "Hello candidate", 24*time.Hour).Err(),
				"set greeting to redis",
			)
		},
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client
}
