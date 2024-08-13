package dummy

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/XiBao/logger/v2/adapters"
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
	Adapter struct{ adapters.Adapter }

	Context struct{}
)

func newContext() *Context {
	return contextPool.Get().(*Context)
}

func releaseContext(ctx *Context) {
	contextPool.Put(ctx)
}

// NewAdapter returns a new adapter. The nop adapter does not log anything and can be used as a placeholder or fallback.
func NewAdapter() adapters.Logger { return &Adapter{} }

func (a *Adapter) With(_ ...any) adapters.Logger                     { return a }
func (a *Adapter) WithLevel(_ adapters.Level) adapters.LoggerContext { return newContext() }
func (a *Adapter) Debug() adapters.LoggerContext                     { return newContext() }
func (a *Adapter) Info() adapters.LoggerContext                      { return newContext() }
func (a *Adapter) Warn() adapters.LoggerContext                      { return newContext() }
func (a *Adapter) Trace() adapters.LoggerContext                     { return newContext() }
func (a *Adapter) Error() adapters.LoggerContext                     { return newContext() }
func (a *Adapter) Fatal() adapters.LoggerContext                     { return newContext() }
func (a *Adapter) Panic() adapters.LoggerContext                     { return newContext() }

func (c *Context) Bytes(_ string, _ []byte) adapters.LoggerContext                    { return c }
func (c *Context) Hex(_ string, _ []byte) adapters.LoggerContext                      { return c }
func (c *Context) RawJSON(_ string, _ []byte) adapters.LoggerContext                  { return c }
func (c *Context) Str(_, _ string) adapters.LoggerContext                             { return c }
func (c *Context) Strs(_ string, _ []string) adapters.LoggerContext                   { return c }
func (c *Context) Stringer(_ string, _ fmt.Stringer) adapters.LoggerContext           { return c }
func (c *Context) Stringers(_ string, _ []fmt.Stringer) adapters.LoggerContext        { return c }
func (c *Context) Int(_ string, _ int) adapters.LoggerContext                         { return c }
func (c *Context) Ints(_ string, _ []int) adapters.LoggerContext                      { return c }
func (c *Context) Int8(_ string, _ int8) adapters.LoggerContext                       { return c }
func (c *Context) Ints8(_ string, _ []int8) adapters.LoggerContext                    { return c }
func (c *Context) Int16(_ string, _ int16) adapters.LoggerContext                     { return c }
func (c *Context) Ints16(_ string, _ []int16) adapters.LoggerContext                  { return c }
func (c *Context) Int32(_ string, _ int32) adapters.LoggerContext                     { return c }
func (c *Context) Ints32(_ string, _ []int32) adapters.LoggerContext                  { return c }
func (c *Context) Int64(_ string, _ int64) adapters.LoggerContext                     { return c }
func (c *Context) Ints64(_ string, _ []int64) adapters.LoggerContext                  { return c }
func (c *Context) Uint(_ string, _ uint) adapters.LoggerContext                       { return c }
func (c *Context) Uints(_ string, _ []uint) adapters.LoggerContext                    { return c }
func (c *Context) Uint8(_ string, _ uint8) adapters.LoggerContext                     { return c }
func (c *Context) Uints8(_ string, _ []uint8) adapters.LoggerContext                  { return c }
func (c *Context) Uint16(_ string, _ uint16) adapters.LoggerContext                   { return c }
func (c *Context) Uints16(_ string, _ []uint16) adapters.LoggerContext                { return c }
func (c *Context) Uint32(_ string, _ uint32) adapters.LoggerContext                   { return c }
func (c *Context) Uints32(_ string, _ []uint32) adapters.LoggerContext                { return c }
func (c *Context) Uint64(_ string, _ uint64) adapters.LoggerContext                   { return c }
func (c *Context) Uints64(_ string, _ []uint64) adapters.LoggerContext                { return c }
func (c *Context) Float32(_ string, _ float32) adapters.LoggerContext                 { return c }
func (c *Context) Floats32(_ string, _ []float32) adapters.LoggerContext              { return c }
func (c *Context) Float64(_ string, _ float64) adapters.LoggerContext                 { return c }
func (c *Context) Floats64(_ string, _ []float64) adapters.LoggerContext              { return c }
func (c *Context) Bool(_ string, _ bool) adapters.LoggerContext                       { return c }
func (c *Context) Bools(_ string, _ []bool) adapters.LoggerContext                    { return c }
func (c *Context) Time(_ string, _ time.Time) adapters.LoggerContext                  { return c }
func (c *Context) Times(_ string, _ []time.Time) adapters.LoggerContext               { return c }
func (c *Context) Dur(_ string, _ time.Duration) adapters.LoggerContext               { return c }
func (c *Context) Durs(_ string, _ []time.Duration) adapters.LoggerContext            { return c }
func (c *Context) TimeDiff(_ string, _ time.Time, _ time.Time) adapters.LoggerContext { return c }
func (c *Context) IPAddr(_ string, _ net.IP) adapters.LoggerContext                   { return c }
func (c *Context) IPPrefix(_ string, _ net.IPNet) adapters.LoggerContext              { return c }
func (c *Context) MACAddr(_ string, _ net.HardwareAddr) adapters.LoggerContext        { return c }
func (c *Context) Err(_ error) adapters.LoggerContext                                 { return c }
func (c *Context) Errs(_ string, _ []error) adapters.LoggerContext                    { return c }
func (c *Context) AnErr(_ string, _ error) adapters.LoggerContext                     { return c }
func (c *Context) Any(_ string, _ any) adapters.LoggerContext                         { return c }
func (c *Context) Array(_ string, _ ...any) adapters.LoggerContext                    { return c }
func (c *Context) Fields(_ adapters.Fields) adapters.LoggerContext                    { return c }
func (c *Context) Stack() adapters.LoggerContext                                      { return c }

func (c *Context) Msg(_ string)            { releaseContext(c) }
func (c *Context) Msgf(_ string, _ ...any) { releaseContext(c) }
func (c *Context) Send()                   { releaseContext(c) }
