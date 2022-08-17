package zapwrapper

import (
	"os"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.Logger, error) {
	errorLogs := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.ErrorLevel
	})

	infoLogs := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zapcore.ErrorLevel
	})

	encoder := zapcore.NewJSONEncoder(zapdriver.NewDevelopmentEncoderConfig())

	coreWrapper := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), infoLogs),
			zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), errorLogs),
		)
	})

	return zapdriver.NewDevelopmentWithCore(coreWrapper)
}
