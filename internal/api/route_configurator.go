package api

import (
	echoLib "github.com/labstack/echo/v4"

	"interview/internal/api/hello"
	"interview/pkg/echo"
)

type Router interface {
	Config(engine *echoLib.Group)
}

func NewRouteConfigurator(helloRouter *hello.Router) echo.RouteConfigurator {
	return func(e *echoLib.Echo) {
		v1 := e.Group("/v1")

		helloRouter.Config(v1)
	}
}
