package echo

import (
	"context"
	"net"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"
)

type RouteConfigurator func(*echo.Echo)

type Engine interface {
	Run() error
	Stop(ctx context.Context) error
}

type Config struct {
	Environment    string
	UnixSocket     string
	Router         RouteConfigurator
	Logger         echo.Logger
	HideBanner     bool
	DisableRecover bool
}

type engine struct {
	echo   *echo.Echo
	config *Config
}

func New(config *Config) (Engine, error) {
	echoEngine := echo.New()
	echoEngine.HideBanner = config.HideBanner

	if !config.DisableRecover {
		echoEngine.Use(middleware.Recover())
	}

	if config.Logger != nil {
		echoEngine.Logger = config.Logger
	} else {
		echoEngine.Logger.SetOutput(colorable.NewColorableStderr())
	}

	config.Router(echoEngine)

	return &engine{
		echo:   echoEngine,
		config: config,
	}, nil
}

func (e *engine) Run() error {
	filename := e.config.UnixSocket
	if fileExists(filename) {
		if err := os.Remove(filename); err != nil {
			return errors.Wrap(err, "remove socket file")
		}
	}

	listener, err := net.Listen("unix", filename)
	if err != nil {
		return errors.Wrap(err, "create net listener")
	}

	if err = os.Chmod(filename, 0o777); err != nil {
		return errors.Wrap(err, "chmod unix socket")
	}

	e.echo.Listener = listener

	return e.echo.Start("")
}

func (e *engine) Stop(ctx context.Context) error {
	return e.echo.Shutdown(ctx)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)

	return !os.IsNotExist(err)
}
