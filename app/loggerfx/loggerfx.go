package loggerfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Logger struct {
	zap.SugaredLogger
}

// ProvideLogger to fx
func ProvideLogger() *Logger {
	logger, _ := zap.NewProduction()
	slogger := logger.Sugar()

	return &Logger{*slogger}
}

// Module provided to fx
var Module = fx.Options(
	fx.Provide(ProvideLogger),
)
