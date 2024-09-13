package sentry

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

const FlushTimeout = 2 * time.Second

type Hook struct{}

func (h Hook) Run(event *zerolog.Event, level zerolog.Level, message string) {
	if level == zerolog.ErrorLevel {
		captured := h.convertEvent(event, level, message)
		hub := sentry.CurrentHub()
		client, scope := hub.Client(), hub.Scope()
		client.CaptureEvent(&captured, &sentry.EventHint{Context: event.GetCtx()}, scope)
	}

	if level == zerolog.FatalLevel || level == zerolog.PanicLevel {
		sentry.Flush(FlushTimeout)
	}
}
