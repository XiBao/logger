/*
 * telemetry
 * hook.go
 * This file is part of telemetry.
 * Copyright (c) 2024.
 * Last modified at Mon, 22 Jul 2024 16:17:16 -0500 by nick.
 *
 * DISCLAIMER: This software is provided "as is" without warranty of any kind, either expressed or implied. The entire
 * risk as to the quality and performance of the software is with you. In no event will the author be liable for any
 * damages, including any general, special, incidental, or consequential damages arising out of the use or inability
 * to use the software (that includes, but not limited to, loss of data, data being rendered inaccurate, or losses
 * sustained by you or third parties, or a failure of the software to operate with any other programs), even if the
 * author has been advised of the possibility of such damages.
 * If a license file is provided with this software, all use of this software is governed by the terms and conditions
 * set forth in that license file. If no license file is provided, no rights are granted to use, modify, distribute,
 * or otherwise exploit this software.
 */

package otlp

import (
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/log/global"
)

const (
	instrumentName = "github.com/XiBao/logger/hook/otel"
)

// Hook struct defines a logger hook for the zerolog logging library.
type Hook struct{}

// Run is the method that gets called on each log event.
// It converts the zerolog event to an OpenTelemetry log record and emits it using the global logger provider.
// In case of PanicLevel or FatalLevel log events, it also attempts to shut down the logger provider gracefully.
//
// Parameters:
// - event: The zerolog event that contains all the log information.
// - constants: The logging constants of the event (e.g., Info, Warn, Error).
// - message: The log message.
//
// The method extracts the context from the event, converts the event to an OpenTelemetry log record,
// and emits the record using the global logger provider. If the log constants is PanicLevel or FatalLevel,
// it shuts down the logger provider to ensure all logs are flushed before the application exits.
func (h Hook) Run(event *zerolog.Event, level zerolog.Level, message string) {
	ctx := event.GetCtx()

	// Extract context from the event.
	record := h.convertEvent(event, level, message)   // Convert zerolog event to OpenTelemetry log record.
	provider := global.GetLoggerProvider()            // Get the global logger provider.
	provider.Logger(instrumentName).Emit(ctx, record) // Emit the log record.
}
