package logger

import "github.com/XiBao/logger/adapters"

// Level defines log levels.
const (
	// DebugLevel defines debug log level.
	DebugLevel = adapters.DebugLevel
	// InfoLevel defines info log level.
	InfoLevel = adapters.InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel = adapters.WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel = adapters.ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel = adapters.FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel = adapters.PanicLevel
	// NoLevel defines an absent log level.
	NoLevel = adapters.NoLevel
	// Disabled disables the logger.
	Disabled = adapters.Disabled

	// TraceLevel defines trace log level.
	TraceLevel = adapters.TraceLevel
	// Values less than TraceLevel are handled as numbers.
)
