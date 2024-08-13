package slog

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/exp/slog"

	"github.com/XiBao/logger/adapters"
)

// Compile-time check that Adapter and Context implements adapters.Logger and adapters.LoggerContext respectively
var (
	_           adapters.Logger        = (*Adapter)(nil)
	_           adapters.LoggerContext = (*Context)(nil)
	contextPool                        = sync.Pool{
		New: func() any {
			return new(Context)
		},
	}
)

type (
	// Adapter is a slog adapter for adapters. It implements the adapters.Logger interface.
	Adapter struct {
		adapters *slog.Logger
	}

	// Context is the slog logging context. It implements the adapters.LoggerContext interface.
	Context struct {
		adapters *slog.Logger
		fields   []any
		level    slog.Level
	}

	ctxKey struct{}
)

// NewAdapter creates a new slog adapter for adapters.
func NewAdapter(l *slog.Logger) adapters.Logger {
	return &Adapter{
		adapters: l,
	}
}

func newContext(level slog.Level, adapters *slog.Logger) *Context {
	ctx := contextPool.Get().(*Context)
	ctx.level = level
	ctx.adapters = adapters
	ctx.fields = make([]any, 0)
	return ctx
}

func releaseContext(ctx *Context) {
	contextPool.Put(ctx)
}

func (a *Adapter) newContext(level slog.Level) *Context {
	return newContext(level, a.adapters)
}

// Ctx returns the Logger associated with the ctx. If no adapters
// is associated, DefaultContextLogger is returned, unless DefaultContextLogger
// is nil, in which case a disabled adapters is returned.
func (a *Adapter) Ctx(ctx context.Context) adapters.Logger {
	if l, ok := ctx.Value(ctxKey{}).(adapters.Logger); ok {
		return l
	}
	return &Adapter{adapters: slog.Default()}
}

func (a *Adapter) WithContext(ctx context.Context) context.Context {
	if _, ok := ctx.Value(ctxKey{}).(adapters.Logger); !ok {
		return ctx
	}
	return context.WithValue(ctx, ctxKey{}, a)
}

// With returns the adapters with the given fields.
func (a *Adapter) With(fields ...any) adapters.Logger {
	return &Adapter{adapters: a.adapters.With(fields...)}
}

// Debug returns a LoggerContext for a debug log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Debug() adapters.LoggerContext {
	return a.newContext(slog.LevelDebug)
}

// Info returns a LoggerContext for an info log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Info() adapters.LoggerContext {
	return a.newContext(slog.LevelInfo)
}

// Warn returns a LoggerContext for a warn log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Warn() adapters.LoggerContext {
	return a.newContext(slog.LevelWarn)
}

// Error returns a LoggerContext for an error log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Error() adapters.LoggerContext {
	return a.newContext(slog.LevelError)
}

// Trace returns a LoggerContext for a trace log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Trace() adapters.LoggerContext {
	return a.newContext(slog.LevelDebug) // Using Error level here because Fatal is not supported by slog
}

// Fatal returns a LoggerContext for a fatal log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Fatal() adapters.LoggerContext {
	return a.newContext(slog.LevelError) // Using Error level here because Fatal is not supported by slog
}

// Panic returns a LoggerContext for a panic log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Panic() adapters.LoggerContext {
	return a.newContext(slog.LevelError) // Using Error level here because Fatal is not supported by slog
}

// WithLevel starts a new message with level.
func (a *Adapter) WithLevel(level adapters.Level) adapters.LoggerContext {
	var slogLevel slog.Level
	switch level {
	case adapters.DebugLevel:
		slogLevel = slog.LevelDebug
	case adapters.InfoLevel:
		slogLevel = slog.LevelInfo
	case adapters.WarnLevel:
		slogLevel = slog.LevelWarn
	case adapters.ErrorLevel, adapters.FatalLevel, adapters.PanicLevel:
		slogLevel = slog.LevelError
	}
	return a.newContext(slogLevel)
}

// Bytes adds the field key with val as a []byte to the adapters context.
func (c *Context) Bytes(key string, value []byte) adapters.LoggerContext {
	c.fields = append(c.fields, slog.String(key, string(value)))

	return c
}

// Hex adds the field key with val as a hex string to the adapters context.
func (c *Context) Hex(key string, value []byte) adapters.LoggerContext {
	c.fields = append(c.fields, slog.String(key, fmt.Sprintf("%x", value)))

	return c
}

// RawJSON adds the field key with val as a raw JSON string to the adapters context.
func (c *Context) RawJSON(key string, value []byte) adapters.LoggerContext {
	c.fields = append(c.fields, slog.String(key, string(value)))

	return c
}

// Str adds the field key with val as a string to the adapters context.
func (c *Context) Str(key string, value string) adapters.LoggerContext {
	c.fields = append(c.fields, slog.String(key, value))

	return c
}

// Strs adds the field key with val as a []string to the adapters context.
func (c *Context) Strs(key string, value []string) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// Stringer adds the field key with val as a fmt.Stringer to the adapters context.
func (c *Context) Stringer(key string, value fmt.Stringer) adapters.LoggerContext {
	c.fields = append(c.fields, slog.String(key, value.String()))

	return c
}

// Stringers adds the field key with val as a []fmt.Stringer to the adapters context.
func (c *Context) Stringers(key string, value []fmt.Stringer) adapters.LoggerContext {
	// Todo: Better way to do this?
	strs := make([]string, len(value))
	for i, str := range value {
		strs[i] = str.String()
	}
	c.fields = append(c.fields, slog.Any(key, strs))

	return c
}

// Int adds the field key with val as a int to the adapters context.
func (c *Context) Int(key string, value int) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Int(key, value))

	return c
}

// Ints adds the field key with val as a []int to the adapters context.
func (c *Context) Ints(key string, value []int) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// Int8 adds the field key with val as a int8 to the adapters context.
func (c *Context) Int8(key string, value int8) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Int64(key, int64(value)))

	return c
}

// Ints8 adds the field key with val as a []int8 to the adapters context.
func (c *Context) Ints8(key string, value []int8) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// Int16 adds the field key with val as a int16 to the adapters context.
func (c *Context) Int16(key string, value int16) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Int64(key, int64(value)))

	return c
}

// Ints16 adds the field key with val as a []int16 to the adapters context.
func (c *Context) Ints16(key string, value []int16) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// Int32 adds the field key with val as a int32 to the adapters context.
func (c *Context) Int32(key string, value int32) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Int64(key, int64(value)))

	return c
}

// Ints32 adds the field key with val as a []int32 to the adapters context.
func (c *Context) Ints32(key string, value []int32) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// Int64 adds the field key with val as a int64 to the adapters context.
func (c *Context) Int64(key string, value int64) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Int64(key, value))

	return c
}

// Ints64 adds the field key with val as a []int64 to the adapters context.
func (c *Context) Ints64(key string, value []int64) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// Uint adds the field key with val as a uint to the adapters context.
func (c *Context) Uint(key string, value uint) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Uint64(key, uint64(value)))

	return c
}

// Uints adds the field key with val as a []uint to the adapters context.
func (c *Context) Uints(key string, value []uint) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// Uint8 adds the field key with val as a uint8 to the adapters context.
func (c *Context) Uint8(key string, value uint8) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Uint64(key, uint64(value)))

	return c
}

// Uints8 adds the field key with val as a []uint8 to the adapters context.
func (c *Context) Uints8(key string, value []uint8) adapters.LoggerContext {
	// Todo: Better way to do this?
	// Convert []uint8 to []uint64
	uints := make([]uint64, len(value))
	for i, v := range value {
		uints[i] = uint64(v)
	}

	c.fields = append(c.fields, slog.Any(key, uints))

	return c
}

// Uint16 adds the field key with val as a uint16 to the adapters context.
func (c *Context) Uint16(key string, value uint16) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Uint64(key, uint64(value)))

	return c
}

// Uints16 adds the field key with val as a []uint16 to the adapters context.
func (c *Context) Uints16(key string, value []uint16) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// Uint32 adds the field key with val as a uint32 to the adapters context.
func (c *Context) Uint32(key string, value uint32) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Uint64(key, uint64(value)))

	return c
}

// Uints32 adds the field key with val as a []uint32 to the adapters context.
func (c *Context) Uints32(key string, value []uint32) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// Uint64 adds the field key with val as a uint64 to the adapters context.
func (c *Context) Uint64(key string, value uint64) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Uint64(key, value))

	return c
}

// Uints64 adds the field key with val as a []uint64 to the adapters context.
func (c *Context) Uints64(key string, value []uint64) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// Float32 adds the field key with val as a float32 to the adapters context.
func (c *Context) Float32(key string, value float32) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Float64(key, float64(value)))

	return c
}

// Floats32 adds the field key with val as a []float32 to the adapters context.
func (c *Context) Floats32(key string, value []float32) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// Float64 adds the field key with val as a float64 to the adapters context.
func (c *Context) Float64(key string, value float64) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Float64(key, value))

	return c
}

// Floats64 adds the field key with val as a []float64 to the adapters context.
func (c *Context) Floats64(key string, value []float64) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// Bool adds the field key with val as a bool to the adapters context.
func (c *Context) Bool(key string, value bool) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Bool(key, value))

	return c
}

// Bools adds the field key with val as a []bool to the adapters context.
func (c *Context) Bools(key string, value []bool) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// Time adds the field key with val as a time.Time to the adapters context.
func (c *Context) Time(key string, value time.Time) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Time(key, value))

	return c
}

// Times adds the field key with val as a []time.Time to the adapters context.
func (c *Context) Times(key string, value []time.Time) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// Dur adds the field key with val as a time.Duration to the adapters context.
func (c *Context) Dur(key string, value time.Duration) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Duration(key, value))

	return c
}

// Durs adds the field key with val as a []time.Duration to the adapters context.
func (c *Context) Durs(key string, value []time.Duration) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// TimeDiff adds the field key with begin and end as a time.Time to the adapters context.
func (c *Context) TimeDiff(key string, begin, end time.Time) adapters.LoggerContext {
	diff := end.Sub(begin)
	c.fields = append(c.fields, slog.Duration(key, diff))

	return c
}

// IPAddr adds the field key with val as a net.IPAddr to the adapters context.
func (c *Context) IPAddr(key string, value net.IP) adapters.LoggerContext {
	c.fields = append(c.fields, slog.String(key, value.String()))

	return c
}

// IPPrefix adds the field key with val as a net.IPPrefix to the adapters context.
func (c *Context) IPPrefix(key string, value net.IPNet) adapters.LoggerContext {
	c.fields = append(c.fields, slog.String(key, value.String()))

	return c
}

// MACAddr adds the field key with val as a net.HardwareAddr to the adapters context.
func (c *Context) MACAddr(key string, value net.HardwareAddr) adapters.LoggerContext {
	c.fields = append(c.fields, slog.String(key, value.String()))

	return c
}

// AnErr adds the field key with val as a error to the adapters context.
func (c *Context) AnErr(key string, value error) adapters.LoggerContext {
	c.fields = append(c.fields, slog.String(key, value.Error()))

	return c
}

// Err adds the field "error" with val as a error to the adapters context.
func (c *Context) Err(value error) adapters.LoggerContext {
	c.AnErr("error", value)

	return c
}

// Errs adds the field "error" with val as a []error to the adapters context.
func (c *Context) Errs(key string, value []error) adapters.LoggerContext {
	// Todo: Better way to do this?
	// Convert []error to []string. If we don't do this, slog prints empty objects
	errs := make([]string, len(value))
	for i, err := range value {
		errs[i] = err.Error()
	}

	c.fields = append(c.fields, slog.Any(key, errs))

	return c
}

// Any adds the field key with val as a arbitrary value to the adapters context.
func (c *Context) Any(key string, value any) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Any(key, value))

	return c
}

// Any adds the field key with val as a arbitrary value to the adapters context.
func (c *Context) Array(key string, value ...any) adapters.LoggerContext {
	c.fields = append(c.fields, slog.Group(key, value))

	return c
}

func (c *Context) Fields(fields adapters.Fields) adapters.LoggerContext {
	for key, value := range fields {
		c.fields = append(c.fields, slog.Any(key, value))
	}

	return c
}

func (c *Context) Stack() adapters.LoggerContext { return c }

// Msg sends the LoggerContext with msg to the adapters.
func (c *Context) Msg(msg string) {
	//nolint:staticcheck // passing a nil context is fine, check slog.Logger.Info implementation for example
	c.adapters.Log(context.TODO(), c.level, msg, c.fields...)
	c.fields = make([]any, 0) // reset fields
	releaseContext(c)
}

// Msgf sends the LoggerContext with formatted msg to the adapters.
func (c *Context) Msgf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	c.Msg(msg)
}

// Send sends the LoggerContext with empty msg to the adapters.
func (c *Context) Send() {
	c.Msg("")
}
