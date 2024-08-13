package zap

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

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
	// Adapter is a zap adapter for adapters. It implements the adapters.Logger interface.
	Adapter struct {
		adapters *zap.Logger
	}

	// Context is the zap logging context. It implements the adapters.LoggerContext interface.
	Context struct {
		adapters *zap.Logger
		fields   []zapcore.Field
		level    zapcore.Level
	}

	ctxKey struct{}
)

// NewAdapter creates a new zap adapter for adapters.
func NewAdapter(l *zap.Logger) adapters.Logger {
	return &Adapter{
		adapters: l,
	}
}

func newContext(level zapcore.Level, adapters *zap.Logger) *Context {
	ctx := contextPool.Get().(*Context)
	ctx.level = level
	ctx.adapters = adapters
	ctx.fields = make([]zapcore.Field, 0)
	return ctx
}

func releaseContext(ctx *Context) {
	contextPool.Put(ctx)
}

func (a *Adapter) newContext(level zapcore.Level) adapters.LoggerContext {
	return newContext(level, a.adapters)
}

// Ctx returns the Logger associated with the ctx. If no adapters
// is associated, DefaultContextLogger is returned, unless DefaultContextLogger
// is nil, in which case a disabled adapters is returned.
func (a *Adapter) Ctx(ctx context.Context) adapters.Logger {
	if l, ok := ctx.Value(ctxKey{}).(adapters.Logger); ok {
		return l
	}
	return &Adapter{adapters: zap.L()}
}

func (a *Adapter) WithContext(ctx context.Context) context.Context {
	if _, ok := ctx.Value(ctxKey{}).(adapters.Logger); !ok {
		return ctx
	}
	return context.WithValue(ctx, ctxKey{}, a)
}

// With returns the adapters with the given fields.
func (a *Adapter) With(fields ...any) adapters.Logger {
	return &Adapter{adapters: a.adapters.Sugar().With(fields...).Desugar()}
}

// Debug returns a LoggerContext for a debug log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Debug() adapters.LoggerContext {
	return a.newContext(zap.DebugLevel)
}

// Info returns a LoggerContext for a info log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Info() adapters.LoggerContext {
	return a.newContext(zap.InfoLevel)
}

// Warn returns a LoggerContext for a warn log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Warn() adapters.LoggerContext {
	return a.newContext(zap.WarnLevel)
}

// Error returns a LoggerContext for a error log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Error() adapters.LoggerContext {
	return a.newContext(zap.ErrorLevel)
}

// Trace returns a LoggerContext for a trace log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Trace() adapters.LoggerContext {
	return a.newContext(zap.DebugLevel)
}

// Panic  returns a LoggerContext for a panic log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Panic() adapters.LoggerContext {
	return a.newContext(zap.PanicLevel)
}

// Fatal returns a LoggerContext for a fatal log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Fatal() adapters.LoggerContext {
	return a.newContext(zap.FatalLevel)
}

// WithLevel starts a new message with level.
func (a *Adapter) WithLevel(level adapters.Level) adapters.LoggerContext {
	var zapLevel zapcore.Level
	if level <= adapters.ErrorLevel {
		zapLevel = zapcore.Level(level)
	} else if level == adapters.FatalLevel {
		zapLevel = zapcore.FatalLevel
	} else if level == adapters.PanicLevel {
		zapLevel = zapcore.PanicLevel
	}
	return a.newContext(zapLevel)
}

func (c *Context) reset() {
	c.fields = make([]zapcore.Field, 0)
}

// Bytes adds the field key with val as a []byte to the adapters context.
func (c *Context) Bytes(key string, value []byte) adapters.LoggerContext {
	c.fields = append(c.fields, zap.ByteString(key, value))

	return c
}

// Hex adds the field key with val as a hex string to the adapters context.
func (c *Context) Hex(key string, value []byte) adapters.LoggerContext {
	c.fields = append(c.fields, zap.String(key, fmt.Sprintf("%x", value)))

	return c
}

// RawJSON adds the field key with val as a raw json string to the adapters context.
func (c *Context) RawJSON(key string, value []byte) adapters.LoggerContext {
	c.fields = append(c.fields, zap.ByteString(key, value))

	return c
}

// Str adds the field key with val as a string to the adapters context.
func (c *Context) Str(key string, value string) adapters.LoggerContext {
	c.fields = append(c.fields, zap.String(key, value))

	return c
}

// Strs adds the field key with val as a []string to the adapters context.
func (c *Context) Strs(key string, value []string) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Strings(key, value))

	return c
}

// Stringer adds the field key with val as a fmt.Stringer to the adapters context.
func (c *Context) Stringer(key string, val fmt.Stringer) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Stringer(key, val))

	return c
}

// Stringers adds the field key with val as a []fmt.Stringer to the adapters context.
func (c *Context) Stringers(key string, vals []fmt.Stringer) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Stringers(key, vals))

	return c
}

// Int adds the field key with val as a int to the adapters context.
func (c *Context) Int(key string, value int) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Int(key, value))

	return c
}

// Ints adds the field key with val as a []int to the adapters context.
func (c *Context) Ints(key string, value []int) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Ints(key, value))

	return c
}

// Int8 adds the field key with val as a int8 to the adapters context.
func (c *Context) Int8(key string, value int8) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Int8(key, value))

	return c
}

// Ints8 adds the field key with val as a []int8 to the adapters context.
func (c *Context) Ints8(key string, value []int8) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Int8s(key, value))

	return c
}

// Int16 adds the field key with val as a int16 to the adapters context.
func (c *Context) Int16(key string, value int16) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Int16(key, value))

	return c
}

// Ints16 adds the field key with val as a []int16 to the adapters context.
func (c *Context) Ints16(key string, value []int16) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Int16s(key, value))

	return c
}

// Int32 adds the field key with val as a int32 to the adapters context.
func (c *Context) Int32(key string, value int32) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Int32(key, value))

	return c
}

// Ints32 adds the field key with val as a []int32 to the adapters context.
func (c *Context) Ints32(key string, value []int32) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Int32s(key, value))

	return c
}

// Int64 adds the field key with val as a int64 to the adapters context.
func (c *Context) Int64(key string, value int64) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Int64(key, value))

	return c
}

// Ints64 adds the field key with val as a []int64 to the adapters context.
func (c *Context) Ints64(key string, value []int64) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Int64s(key, value))

	return c
}

// Uint adds the field key with val as a uint to the adapters context.
func (c *Context) Uint(key string, value uint) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Uint(key, value))

	return c
}

// Uints adds the field key with val as a []uint to the adapters context.
func (c *Context) Uints(key string, value []uint) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Uints(key, value))

	return c
}

// Uint8 adds the field key with val as a uint8 to the adapters context.
func (c *Context) Uint8(key string, value uint8) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Uint8(key, value))

	return c
}

// Uints8 adds the field key with val as a []uint8 to the adapters context.
func (c *Context) Uints8(key string, value []uint8) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Uint8s(key, value))

	return c
}

// Uint16 adds the field key with val as a uint16 to the adapters context.
func (c *Context) Uint16(key string, value uint16) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Uint16(key, value))

	return c
}

// Uints16 adds the field key with val as a []uint16 to the adapters context.
func (c *Context) Uints16(key string, value []uint16) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Uint16s(key, value))

	return c
}

// Uint32 adds the field key with val as a uint32 to the adapters context.
func (c *Context) Uint32(key string, value uint32) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Uint32(key, value))

	return c
}

// Uints32 adds the field key with val as a []uint32 to the adapters context.
func (c *Context) Uints32(key string, value []uint32) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Uint32s(key, value))

	return c
}

// Uint64 adds the field key with val as a uint64 to the adapters context.
func (c *Context) Uint64(key string, value uint64) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Uint64(key, value))

	return c
}

// Uints64 adds the field key with val as a []uint64 to the adapters context.
func (c *Context) Uints64(key string, value []uint64) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Uint64s(key, value))

	return c
}

// Float32 adds the field key with val as a float32 to the adapters context.
func (c *Context) Float32(key string, value float32) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Float32(key, value))

	return c
}

// Floats32 adds the field key with val as a []float32 to the adapters context.
func (c *Context) Floats32(key string, value []float32) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Float32s(key, value))

	return c
}

// Float64 adds the field key with val as a float64 to the adapters context.
func (c *Context) Float64(key string, value float64) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Float64(key, value))

	return c
}

// Floats64 adds the field key with val as a []float64 to the adapters context.
func (c *Context) Floats64(key string, value []float64) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Float64s(key, value))

	return c
}

// Bool adds the field key with val as a bool to the adapters context.
func (c *Context) Bool(key string, value bool) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Bool(key, value))

	return c
}

// Bools adds the field key with val as a []bool to the adapters context.
func (c *Context) Bools(key string, value []bool) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Bools(key, value))

	return c
}

// Time adds the field key with val as a time.Time to the adapters context.
func (c *Context) Time(key string, value time.Time) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Time(key, value))

	return c
}

// Times adds the field key with val as a []time.Time to the adapters context.
func (c *Context) Times(key string, value []time.Time) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Times(key, value))

	return c
}

// Dur adds the field key with val as a time.Duration to the adapters context.
func (c *Context) Dur(key string, value time.Duration) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Duration(key, value))

	return c
}

// Durs adds the field key with val as a []time.Duration to the adapters context.
func (c *Context) Durs(key string, value []time.Duration) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Durations(key, value))

	return c
}

// TimeDiff adds the field key with begin and end as a time.Time to the adapters context.
func (c *Context) TimeDiff(key string, begin, end time.Time) adapters.LoggerContext {
	diff := end.Sub(begin)
	c.fields = append(c.fields, zap.Duration(key, diff))

	return c
}

// IPAddr adds the field key with val as a net.IP to the adapters context.
func (c *Context) IPAddr(key string, value net.IP) adapters.LoggerContext {
	c.fields = append(c.fields, zap.String(key, value.String()))

	return c
}

// IPPrefix adds the field key with val as a net.IPNet to the adapters context.
func (c *Context) IPPrefix(key string, value net.IPNet) adapters.LoggerContext {
	c.fields = append(c.fields, zap.String(key, value.String()))

	return c
}

// MACAddr adds the field key with val as a net.HardwareAddr to the adapters context.
func (c *Context) MACAddr(key string, value net.HardwareAddr) adapters.LoggerContext {
	c.fields = append(c.fields, zap.String(key, value.String()))

	return c
}

// AnErr adds the field key with val as a error to the adapters context.
func (c *Context) AnErr(key string, err error) adapters.LoggerContext {
	c.fields = append(c.fields, zap.NamedError(key, err))

	return c
}

// Err adds the field key with val as a error to the adapters context.
func (c *Context) Err(err error) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Error(err))

	return c
}

// Errs adds the field key with val as a []error to the adapters context.
func (c *Context) Errs(key string, errs []error) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Errors(key, errs))

	return c
}

// Any adds the field key with val as a arbitrary value to the adapters context.
func (c *Context) Any(key string, value any) adapters.LoggerContext {
	c.fields = append(c.fields, zap.Any(key, value))

	return c
}

// Array adds the field key with val as arbitrary array value to the adapters context.
func (c *Context) Array(key string, value ...any) adapters.LoggerContext {
	return c
}

func (c *Context) Fields(fields adapters.Fields) adapters.LoggerContext {
	for k, v := range fields {
		c.fields = append(c.fields, zap.Any(k, v))
	}

	return c
}

// Msg sends the LoggerContext with msg to the adapters.
func (c *Context) Msg(msg string) {
	c.adapters.Log(c.level, msg, c.fields...)
	releaseContext(c)
}

func (c *Context) Stack() adapters.LoggerContext { return c }

// Msgf sends the LoggerContext with formatted msg to the adapters.
func (c *Context) Msgf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	c.Msg(msg)
}

// Send sends the LoggerContext with empty msg to the adapters.
func (c *Context) Send() {
	c.Msg("")
}
