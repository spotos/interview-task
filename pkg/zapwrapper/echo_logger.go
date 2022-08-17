package zapwrapper

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapio"
)

type echoLogger struct {
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
	prefix        string
}

func EchoLogger(logger *zap.Logger) echo.Logger {
	return &echoLogger{logger: logger, sugaredLogger: logger.Sugar(), prefix: "echo"}
}

func (e *echoLogger) Output() io.Writer {
	return &zapio.Writer{
		Log:   e.logger,
		Level: e.getLogLevel(),
	}
}

func (e *echoLogger) Prefix() string {
	return e.prefix
}

func (e *echoLogger) Level() log.Lvl {
	switch e.getLogLevel() {
	case zapcore.DebugLevel:
		return log.DEBUG
	case zapcore.InfoLevel:
		return log.INFO
	case zapcore.WarnLevel:
		return log.WARN
	case zapcore.ErrorLevel, zapcore.FatalLevel, zapcore.PanicLevel, zapcore.DPanicLevel:
		return log.ERROR
	}

	return log.ERROR
}

func (e *echoLogger) SetOutput(io.Writer) {}
func (e *echoLogger) SetPrefix(string)    {}
func (e *echoLogger) SetLevel(log.Lvl)    {}
func (e *echoLogger) SetHeader(string)    {}

func (e *echoLogger) Print(i ...interface{}) {
	if len(i) == 0 {
		return
	}

	if len(i) == 1 {
		if _, ok := i[0].(error); ok {
			e.sugaredLogger.Error(i[0])
		}
	}

	hasErrors := false

	for _, param := range i {
		if _, ok := param.(error); ok {
			hasErrors = true

			break
		}
	}

	if hasErrors {
		e.sugaredLogger.Error(i...)
	} else {
		e.sugaredLogger.Info(i...)
	}
}

func (e *echoLogger) Printf(format string, args ...interface{}) {
	e.sugaredLogger.Infof(format, args...)
}
func (e *echoLogger) Printj(l log.JSON) {
	e.sugaredLogger.Info(l)
}

func (e *echoLogger) Debug(i ...interface{}) {
	e.sugaredLogger.Debug(i...)
}

func (e echoLogger) Debugf(format string, args ...interface{}) {
	e.sugaredLogger.Debugf(format, args...)
}

func (e *echoLogger) Debugj(l log.JSON) {
	e.sugaredLogger.Debug(l)
}

func (e *echoLogger) Info(i ...interface{}) {
	e.sugaredLogger.Info(i...)
}

func (e *echoLogger) Infof(format string, args ...interface{}) {
	e.sugaredLogger.Infof(format, args...)
}

func (e *echoLogger) Infoj(l log.JSON) {
	e.sugaredLogger.Info(l)
}

func (e *echoLogger) Warn(i ...interface{}) {
	e.sugaredLogger.Warn(i...)
}

func (e *echoLogger) Warnf(format string, args ...interface{}) {
	e.sugaredLogger.Warnf(format, args...)
}

func (e *echoLogger) Warnj(l log.JSON) {
	e.sugaredLogger.Warn(l)
}

func (e *echoLogger) Error(i ...interface{}) {
	e.sugaredLogger.Error(i...)
}

func (e *echoLogger) Errorf(format string, args ...interface{}) {
	e.sugaredLogger.Errorf(format, args...)
}

func (e *echoLogger) Errorj(l log.JSON) {
	e.sugaredLogger.Error(l)
}

func (e *echoLogger) Fatal(i ...interface{}) {
	e.sugaredLogger.Fatal(i...)
}

func (e *echoLogger) Fatalj(l log.JSON) {
	e.sugaredLogger.Fatal(l)
}

func (e *echoLogger) Fatalf(format string, args ...interface{}) {
	e.sugaredLogger.Fatalf(format, args...)
}

func (e *echoLogger) Panic(i ...interface{}) {
	e.sugaredLogger.Panic(i...)
}

func (e *echoLogger) Panicj(l log.JSON) {
	e.sugaredLogger.Panic(l)
}

func (e *echoLogger) Panicf(format string, args ...interface{}) {
	e.sugaredLogger.Panicf(format, args...)
}

func (e *echoLogger) getLogLevel() zapcore.Level {
	levels := []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
	}

	level := zapcore.DebugLevel

	for _, lvl := range levels {
		if e.logger.Core().Enabled(lvl) {
			level = lvl

			break
		}
	}

	return level
}
