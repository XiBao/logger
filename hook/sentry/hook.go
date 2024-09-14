package sentry

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

const FlushTimeout = 2 * time.Second

type Hook struct{}

func NewHook() *Hook {
	return new(Hook)
}

func (h Hook) Run(event *zerolog.Event, level zerolog.Level, message string) {
	if level == zerolog.ErrorLevel {
		ctx := event.GetCtx()
		captured := h.convertEvent(event, level, message)
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub()
			if client, scope := hub.Client(), hub.Scope(); client != nil {
				client.CaptureEvent(&captured, &sentry.EventHint{Context: ctx}, scope)
				return
			}
		}
		hub.CaptureEvent(&captured)
	}

	if level == zerolog.FatalLevel || level == zerolog.PanicLevel {
		sentry.Flush(FlushTimeout)
	}
}
