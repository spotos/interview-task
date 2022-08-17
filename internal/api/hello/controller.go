package hello

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"interview/internal/config"
	"interview/pkg/queues"
)

type Controller struct {
	redis        redis.UniversalClient
	queuesClient queues.Client
}

func NewController(redis redis.UniversalClient, queuesClient queues.Client) *Controller {
	return &Controller{
		redis:        redis,
		queuesClient: queuesClient,
	}
}

func (c *Controller) Hello(ctx echo.Context) error {
	value, err := c.redis.Get(ctx.Request().Context(), config.GreetingKey).Result()
	if err != nil {
		return errors.Wrap(err, "get greeting from redis")
	}

	result := echo.Map{"greeting": value}

	if err := c.queuesClient.Publish(ctx.Request().Context(), config.Topic, result); err != nil {
		return errors.Wrap(err, "publish result to queue")
	}

	return ctx.JSON(http.StatusOK, result)
}
