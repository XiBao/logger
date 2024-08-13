package zerolog

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/XiBao/logger/v2/adapters"
)

// Compile-time check that Adapter and Context implements onelog.Logger and onelog.LoggerContext respectively
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
	// Adapter is a zerolog adapter for adapters. It implements the adapters.Logger interface.
	Adapter struct {
		adapters.Adapter
		adapters *zerolog.Logger
	}

	// Context is the zerolog logging context. It implements the adapters.LoggerContext interface.
	Context struct {
		event *zerolog.Event
	}
)

// NewAdapter creates a new zerolog adapter for adapters.
func NewAdapter(l *zerolog.Logger) adapters.Logger {
	return &Adapter{
		adapters: l,
	}
}

func newContext(event *zerolog.Event) *Context {
	ctx := contextPool.Get().(*Context)
	ctx.event = event
	return ctx
}

func releaseContext(ctx *Context) {
	contextPool.Put(ctx)
}

// With returns the adapters with the given fields.
func (a *Adapter) With(fields ...any) adapters.Logger {
	adapters := a.adapters.With().Fields(fields).Logger()
	return &Adapter{adapters: &adapters}
}

// WithLevel starts a new message with level.
func (a *Adapter) WithLevel(level adapters.Level) adapters.LoggerContext {
	return newContext(a.adapters.WithLevel(zerolog.Level(level)))
}

// Debug returns a LoggerContext for a debug log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Debug() adapters.LoggerContext {
	return newContext(a.adapters.Debug())
}

// Info returns a LoggerContext for a info log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Info() adapters.LoggerContext {
	return newContext(a.adapters.Info())
}

// Warn returns a LoggerContext for a warn log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Warn() adapters.LoggerContext {
	return newContext(a.adapters.Warn())
}

// Error returns a LoggerContext for a error log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Error() adapters.LoggerContext {
	return newContext(a.adapters.Error())
}

// Fatal returns a LoggerContext for a fatal log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Fatal() adapters.LoggerContext {
	return newContext(a.adapters.Fatal())
}

// Fatal returns a LoggerContext for a fatal log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Panic() adapters.LoggerContext {
	return newContext(a.adapters.Panic())
}

// Trace returns a LoggerContext for a trace log. To send the log, use the Msg or Msgf methods.
func (a *Adapter) Trace() adapters.LoggerContext {
	return newContext(a.adapters.Trace())
}

// Bytes adds the field key with val as a []byte to the adapters context.
func (c *Context) Bytes(key string, value []byte) adapters.LoggerContext {
	c.event.Bytes(key, value)

	return c
}

// Hex adds the field key with val as a hex string to the adapters context.
func (c *Context) Hex(key string, value []byte) adapters.LoggerContext {
	c.event.Hex(key, value)

	return c
}

// RawJSON adds the field key with val as a json.RawMessage to the adapters context.
func (c *Context) RawJSON(key string, value []byte) adapters.LoggerContext {
	c.event.RawJSON(key, value)

	return c
}

// Str adds the field key with val as a string to the adapters context.
func (c *Context) Str(key, value string) adapters.LoggerContext {
	c.event.Str(key, value)

	return c
}

// Strs adds the field key with val as a []string to the adapters context.
func (c *Context) Strs(key string, value []string) adapters.LoggerContext {
	c.event.Strs(key, value)

	return c
}

// Stringer adds the field key with val as a fmt.Stringer to the adapters context.
func (c *Context) Stringer(key string, val fmt.Stringer) adapters.LoggerContext {
	c.event.Stringer(key, val)

	return c
}

// Stringers adds the field key with val as a []fmt.Stringer to the adapters context.
func (c *Context) Stringers(key string, vals []fmt.Stringer) adapters.LoggerContext {
	c.event.Stringers(key, vals)

	return c
}

// Int adds the field key with i as a int to the adapters context.
func (c *Context) Int(key string, value int) adapters.LoggerContext {
	c.event.Int(key, value)

	return c
}

// Ints adds the field key with i as a []int to the adapters context.
func (c *Context) Ints(key string, value []int) adapters.LoggerContext {
	c.event.Ints(key, value)

	return c
}

// Int8 adds the field key with i as a int8 to the adapters context.
func (c *Context) Int8(key string, value int8) adapters.LoggerContext {
	c.event.Int8(key, value)

	return c
}

// Ints8 adds the field key with i as a []int8 to the adapters context.
func (c *Context) Ints8(key string, value []int8) adapters.LoggerContext {
	c.event.Ints8(key, value)

	return c
}

// Int16 adds the field key with i as a int16 to the adapters context.
func (c *Context) Int16(key string, value int16) adapters.LoggerContext {
	c.event.Int16(key, value)

	return c
}

// Ints16 adds the field key with i as a []int16 to the adapters context.
func (c *Context) Ints16(key string, value []int16) adapters.LoggerContext {
	c.event.Ints16(key, value)

	return c
}

// Int32 adds the field key with i as a int32 to the adapters context.
func (c *Context) Int32(key string, value int32) adapters.LoggerContext {
	c.event.Int32(key, value)

	return c
}

// Ints32 adds the field key with i as a []int32 to the adapters context.
func (c *Context) Ints32(key string, value []int32) adapters.LoggerContext {
	c.event.Ints32(key, value)

	return c
}

// Int64 adds the field key with i as a int64 to the adapters context.
func (c *Context) Int64(key string, value int64) adapters.LoggerContext {
	c.event.Int64(key, value)

	return c
}

// Ints64 adds the field key with i as a []int64 to the adapters context.
func (c *Context) Ints64(key string, value []int64) adapters.LoggerContext {
	c.event.Ints64(key, value)

	return c
}

// Uint adds the field key with i as a uint to the adapters context.
func (c *Context) Uint(key string, value uint) adapters.LoggerContext {
	c.event.Uint(key, value)

	return c
}

// Uints adds the field key with i as a []uint to the adapters context.
func (c *Context) Uints(key string, value []uint) adapters.LoggerContext {
	c.event.Uints(key, value)

	return c
}

// Uint8 adds the field key with i as a uint8 to the adapters context.
func (c *Context) Uint8(key string, value uint8) adapters.LoggerContext {
	c.event.Uint8(key, value)

	return c
}

// Uints8 adds the field key with i as a []uint8 to the adapters context.
func (c *Context) Uints8(key string, value []uint8) adapters.LoggerContext {
	c.event.Uints8(key, value)

	return c
}

// Uint16 adds the field key with i as a uint16 to the adapters context.
func (c *Context) Uint16(key string, value uint16) adapters.LoggerContext {
	c.event.Uint16(key, value)

	return c
}

// Uints16 adds the field key with i as a []uint16 to the adapters context.
func (c *Context) Uints16(key string, value []uint16) adapters.LoggerContext {
	c.event.Uints16(key, value)

	return c
}

// Uint32 adds the field key with i as a uint32 to the adapters context.
func (c *Context) Uint32(key string, value uint32) adapters.LoggerContext {
	c.event.Uint32(key, value)

	return c
}

// Uints32 adds the field key with i as a []uint32 to the adapters context.
func (c *Context) Uints32(key string, value []uint32) adapters.LoggerContext {
	c.event.Uints32(key, value)

	return c
}

// Uint64 adds the field key with i as a uint64 to the adapters context.
func (c *Context) Uint64(key string, value uint64) adapters.LoggerContext {
	c.event.Uint64(key, value)

	return c
}

// Uints64 adds the field key with i as a []uint64 to the adapters context.
func (c *Context) Uints64(key string, value []uint64) adapters.LoggerContext {
	c.event.Uints64(key, value)

	return c
}

// Float32 adds the field key with f as a float32 to the adapters context.
func (c *Context) Float32(key string, value float32) adapters.LoggerContext {
	c.event.Float32(key, value)

	return c
}

// Floats32 adds the field key with f as a []float32 to the adapters context.
func (c *Context) Floats32(key string, value []float32) adapters.LoggerContext {
	c.event.Floats32(key, value)

	return c
}

// Float64 adds the field key with f as a float64 to the adapters context.
func (c *Context) Float64(key string, value float64) adapters.LoggerContext {
	c.event.Float64(key, value)

	return c
}

// Floats64 adds the field key with f as a []float64 to the adapters context.
func (c *Context) Floats64(key string, value []float64) adapters.LoggerContext {
	c.event.Floats64(key, value)

	return c
}

// Bool adds the field key with b as a bool to the adapters context.
func (c *Context) Bool(key string, value bool) adapters.LoggerContext {
	c.event.Bool(key, value)

	return c
}

// Bools adds the field key with b as a []bool to the adapters context.
func (c *Context) Bools(key string, value []bool) adapters.LoggerContext {
	c.event.Bools(key, value)

	return c
}

// Time adds the field key with t as a time.Time to the adapters context.
func (c *Context) Time(key string, value time.Time) adapters.LoggerContext {
	c.event.Time(key, value)

	return c
}

// Times adds the field key with t as a []time.Time to the adapters context.
func (c *Context) Times(key string, value []time.Time) adapters.LoggerContext {
	c.event.Times(key, value)

	return c
}

// Dur adds the field key with d as a time.Duration to the adapters context.
func (c *Context) Dur(key string, value time.Duration) adapters.LoggerContext {
	c.event.Dur(key, value)

	return c
}

// Durs adds the field key with d as a []time.Duration to the adapters context.
func (c *Context) Durs(key string, value []time.Duration) adapters.LoggerContext {
	c.event.Durs(key, value)

	return c
}

// TimeDiff adds the field key with begin and end as a time.Time to the adapters context.
func (c *Context) TimeDiff(key string, begin, end time.Time) adapters.LoggerContext {
	c.event.TimeDiff(key, begin, end)

	return c
}

// IPAddr adds the field key with ip as a net.IP to the adapters context.
func (c *Context) IPAddr(key string, value net.IP) adapters.LoggerContext {
	c.event.IPAddr(key, value)

	return c
}

// IPPrefix adds the field key with ip as a net.IPNet to the adapters context.
func (c *Context) IPPrefix(key string, value net.IPNet) adapters.LoggerContext {
	c.event.IPPrefix(key, value)

	return c
}

// MACAddr adds the field key with ip as a net.HardwareAddr to the adapters context.
func (c *Context) MACAddr(key string, value net.HardwareAddr) adapters.LoggerContext {
	c.event.MACAddr(key, value)

	return c
}

// Err adds the field "error" with err as a error to the adapters context.
func (c *Context) Err(err error) adapters.LoggerContext {
	c.event.Err(err)

	return c
}

// Errs adds the field key with errs as a []error to the adapters context.
func (c *Context) Errs(key string, errs []error) adapters.LoggerContext {
	c.event.Errs(key, errs)

	return c
}

// AnErr adds the field key with err as a error to the adapters context.
func (c *Context) AnErr(key string, err error) adapters.LoggerContext {
	c.event.AnErr(key, err)

	return c
}

// Any adds the field key with val as a arbitrary value to the adapters context.
func (c *Context) Any(key string, value any) adapters.LoggerContext {
	c.event.Any(key, value)

	return c
}

// Fields adds the fields to the adapters context.
func (c *Context) Fields(fields adapters.Fields) adapters.LoggerContext {
	c.event.Fields(fields)

	return c
}

func (c *Context) Array(key string, value ...any) adapters.LoggerContext {
	arr := zerolog.Arr()
	for _, v := range value {
		arr.Interface(v)
	}
	c.event.Array(key, arr)

	return c
}

// Msg sends the LoggerContext with msg to the adapters.
func (c *Context) Msg(msg string) {
	c.event.Msg(msg)
	releaseContext(c)
}

// Msgf sends the LoggerContext with formatted msg to the adapters.
func (c *Context) Msgf(format string, v ...any) {
	c.event.Msgf(format, v...)
	releaseContext(c)
}

// Send sends the LoggerContext to the adapters.
func (c *Context) Send() {
	c.event.Send()
	releaseContext(c)
}

// Stack enables stack trace printing for the error passed to Err().
func (c *Context) Stack() adapters.LoggerContext {
	c.event.Stack()
	return c
}
