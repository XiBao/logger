package logger

import (
	"context"

	"github.com/XiBao/logger/v2/adapters"
	"github.com/XiBao/logger/v2/adapters/dummy"
)

var defaultLogger = (adapters.Logger)(new(dummy.Adapter))

func SetGlobalLogger(logger adapters.Logger) {
	defaultLogger = logger
}

func Logger() adapters.Logger {
	return defaultLogger
}

func Ctx(ctx context.Context) adapters.Logger {
	if l, ok := ctx.Value(adapters.CtxKey{}).(adapters.Logger); ok {
		return l
	}
	return Logger()
}
