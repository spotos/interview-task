package hello

import (
	"github.com/labstack/echo/v4"
)

type Router struct {
	controller *Controller
}

func NewRouter(controller *Controller) *Router {
	return &Router{
		controller: controller,
	}
}

func (r *Router) Config(engine *echo.Group) {
	engine.GET("/hello", r.controller.Hello)
}
