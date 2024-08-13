package logger

import (
	"github.com/XiBao/logger/adapters"
	"github.com/XiBao/logger/adapters/dummy"
)

var defaultLogger = (adapters.Logger)(new(dummy.Adapter))

func SetGlobalLogger(logger adapters.Logger) {
	defaultLogger = logger
}

func Logger() adapters.Logger {
	return defaultLogger
}
