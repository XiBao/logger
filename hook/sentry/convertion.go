package sentry

import (
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

func (h Hook) convertEvent(e *zerolog.Event, level zerolog.Level, message string) sentry.Event {
	var record sentry.Event

	record.Level = sentry.Level(level.String())
	record.Message = message
	record.Timestamp = zerolog.TimestampFunc()
	fields := convertFields(e)
	record.Extra = make(map[string]interface{}, len(fields))
	for k, v := range fields {
		if err, ok := v.(error); ok {
			record.Exception = append(record.Exception, sentry.Exception{
				Value:      err.Error(),
				Stacktrace: sentry.ExtractStacktrace(err),
			})
		} else {
			record.Extra[k] = v
		}
	}
	return record
}

func convertFields(e *zerolog.Event) map[string]interface{} {
	kvs := make(map[string]interface{})

	// Extract fields from the event and convert them
	e.Fields(func(key string, value interface{}) {
		kvs[key] = value
	})

	return kvs
}
