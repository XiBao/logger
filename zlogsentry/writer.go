package zlogsentry

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
	"unsafe"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
)

type ErrWithStackTrace struct {
	Err        string             `json:"error"`
	Stacktrace *sentry.Stacktrace `json:"stacktrace"`
}

var levelsMapping = map[zerolog.Level]sentry.Level{
	zerolog.DebugLevel: sentry.LevelDebug,
	zerolog.InfoLevel:  sentry.LevelInfo,
	zerolog.WarnLevel:  sentry.LevelWarning,
	zerolog.ErrorLevel: sentry.LevelError,
	zerolog.FatalLevel: sentry.LevelFatal,
	zerolog.PanicLevel: sentry.LevelFatal,
}

var _ = io.WriteCloser(new(Writer))

type Writer struct {
	hub *sentry.Hub

	levels          map[zerolog.Level]struct{}
	flushTimeout    time.Duration
	withBreadcrumbs bool
}

// addBreadcrumb adds event as a breadcrumb
func (w *Writer) addBreadcrumb(event *sentry.Event) {
	if !w.withBreadcrumbs {
		return
	}

	// category is totally optional, but it's nice to have
	var category string
	if _, ok := event.Extra["category"]; ok {
		if v, ok := event.Extra["category"].(string); ok {
			category = v
		}
	}

	w.hub.AddBreadcrumb(&sentry.Breadcrumb{
		Category: category,
		Message:  event.Message,
		Level:    event.Level,
		Data:     event.Extra,
	}, nil)
}

func (w *Writer) Write(data []byte) (int, error) {
	n := len(data)
	lvl, err := w.parseLogLevel(data)
	if err != nil {
		return n, nil
	}

	event, ok := w.parseLogEvent(data)
	if !ok {
		return n, nil
	}

	if _, enabled := w.levels[lvl]; !enabled {
		// if the level is not enabled, add event as a breadcrumb
		w.addBreadcrumb(event)
		return n, nil
	}

	w.hub.CaptureEvent(event)
	// should flush before os.Exit
	if event.Level == sentry.LevelFatal {
		w.hub.Flush(w.flushTimeout)
	}

	return len(data), nil
}

// implements zerolog.LevelWriter
func (w *Writer) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	n = len(p)

	event, ok := w.parseLogEvent(p)
	if !ok {
		return
	}
	event.Level, ok = levelsMapping[level]
	if !ok {
		return
	}

	if _, enabled := w.levels[level]; !enabled {
		// if the level is not enabled, add event as a breadcrumb
		w.addBreadcrumb(event)
		return
	}

	w.hub.CaptureEvent(event)
	// should flush before os.Exit
	if event.Level == sentry.LevelFatal {
		w.hub.Flush(w.flushTimeout)
	}
	return
}

func (w *Writer) Close() error {
	w.hub.Flush(w.flushTimeout)
	return nil
}

// parses the log level from the encoded log
func (w *Writer) parseLogLevel(data []byte) (zerolog.Level, error) {
	lvlStr := gjson.GetBytes(data, zerolog.LevelFieldName).String()

	return zerolog.ParseLevel(lvlStr)
}

func (w *Writer) parseLogEvent(data []byte) (*sentry.Event, bool) {
	const logger = "zerolog"

	event := sentry.Event{
		Timestamp: time.Now().UTC(),
		Logger:    logger,
		Contexts:  make(map[string]sentry.Context),
	}

	isStack := false
	var errExept []sentry.Exception
	payload := make(sentry.Context)

	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		switch key.String() {
		// case zerolog.LevelFieldName, zerolog.TimestampFieldName:
		case zerolog.MessageFieldName:
			event.Message = value.String()
		case zerolog.ErrorFieldName:
			errExept = append(errExept, sentry.Exception{
				Value:      value.String(),
				Stacktrace: newStacktrace(),
			})
		case zerolog.ErrorStackFieldName:
			var e ErrWithStackTrace
			err := json.Unmarshal([]byte(value.Raw), &e)
			if err != nil {
				event.Level = sentry.LevelError
				event.Exception = append(event.Exception, sentry.Exception{
					Value:      err.Error(),
					Stacktrace: sentry.ExtractStacktrace(err),
				})
				event.Message = fmt.Sprintf("Error unmarshal: %s", value)
				break
			}
			event.Exception = append(event.Exception, sentry.Exception{
				Value:      e.Err,
				Stacktrace: e.Stacktrace,
			})
			isStack = true
		default:
			payload[string(key.String())] = value.String()
		}
		return true
	})

	if len(payload) != 0 {
		event.Contexts["payload"] = payload
	}
	if !isStack && len(errExept) > 0 {
		event.Exception = errExept
	}

	return &event, true
}

func newStacktrace() *sentry.Stacktrace {
	const (
		currentModule = "github.com/XiBao/logger/zlogsentry"
		zerologModule = "github.com/rs/zerolog"
	)

	st := sentry.NewStacktrace()

	threshold := len(st.Frames) - 1
	// drop current module frames
	for ; threshold > 0 && st.Frames[threshold].Module == currentModule; threshold-- {
	}

outer:
	// try to drop zerolog module frames after logger call point
	for i := threshold; i > 0; i-- {
		if st.Frames[i].Module == zerologModule {
			for j := i - 1; j >= 0; j-- {
				if st.Frames[j].Module != zerologModule {
					threshold = j
					break outer
				}
			}

			break
		}
	}

	st.Frames = st.Frames[:threshold+1]

	return st
}

func bytesToStrUnsafe(data []byte) string {
	// h := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	// return *(*string)(unsafe.Pointer(&reflect.StringHeader{Data: h.Data, Len: h.Len}))
	return unsafe.String(unsafe.SliceData(data), len(data))
}

type WriterOption interface {
	apply(*config)
}

type optionFunc func(*config)

func (fn optionFunc) apply(c *config) { fn(c) }

type config struct {
	levels             []zerolog.Level
	ignoreErrors       []string
	release            string
	environment        string
	serverName         string
	tracesSampler      sentry.TracesSampler
	sampleRate         float64
	tracesSampleRate   float64
	profilesSampleRate float64
	flushTimeout       time.Duration
	maxErrorDepth      int
	breadcrumbs        bool
	enableTracing      bool
	debug              bool
}

// WithLevels configures zerolog levels that have to be sent to Sentry. Default levels are error, fatal, panic
func WithLevels(levels ...zerolog.Level) WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.levels = levels
	})
}

// WithSampleRate configures the sample rate as a percentage of events to be sent in the range of 0.0 to 1.0
func WithSampleRate(rate float64) WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.sampleRate = rate
	})
}

func WithRelease(release string) WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.release = release
	})
}

func WithEnvironment(environment string) WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.environment = environment
	})
}

// WithServerName configures the server name field for events. Default value is OS hostname
func WithServerName(serverName string) WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.serverName = serverName
	})
}

// WithIgnoreErrors configures the list of regexp strings that will be used to match against event's message
// and if applicable, caught errors type and value. If the match is found, then a whole event will be dropped.
func WithIgnoreErrors(reList []string) WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.ignoreErrors = reList
	})
}

// WithBreadcrumbs enables sentry client breadcrumbs.
func WithBreadcrumbs() WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.breadcrumbs = true
	})
}

func WithEnableTracing(enableTracing bool) WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.enableTracing = enableTracing
	})
}

func WithTracesSampleRate(tracesSampleRate float64) WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.tracesSampleRate = tracesSampleRate
	})
}

func WithTracesSampler(tracesSampler sentry.TracesSampler) WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.tracesSampler = tracesSampler
	})
}

func WithProfilesSampleRate(profilesSampleRate float64) WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.profilesSampleRate = profilesSampleRate
	})
}

// WithMaxErrorDepth sets the max depth of error chain.
func WithMaxErrorDepth(maxErrorDepth int) WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.maxErrorDepth = maxErrorDepth
	})
}

// WithDebug enables sentry client debug logs
func WithDebug() WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.debug = true
	})
}

func New(dsn string, opts ...WriterOption) (*Writer, error) {
	cfg := newDefaultConfig()
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.apply(&cfg)
		}
		err := sentry.Init(sentry.ClientOptions{
			Dsn:                dsn,
			SampleRate:         cfg.sampleRate,
			Release:            cfg.release,
			Environment:        cfg.environment,
			ServerName:         cfg.serverName,
			IgnoreErrors:       cfg.ignoreErrors,
			EnableTracing:      cfg.enableTracing,
			TracesSampleRate:   cfg.tracesSampleRate,
			ProfilesSampleRate: cfg.profilesSampleRate,
			TracesSampler:      cfg.tracesSampler,
			MaxErrorDepth:      cfg.maxErrorDepth,
			Debug:              cfg.debug,
		})
		if err != nil {
			return nil, err
		}
	}

	levels := make(map[zerolog.Level]struct{}, len(cfg.levels))
	for _, lvl := range cfg.levels {
		levels[lvl] = struct{}{}
	}

	return &Writer{
		hub:             sentry.CurrentHub(),
		levels:          levels,
		flushTimeout:    cfg.flushTimeout,
		withBreadcrumbs: cfg.breadcrumbs,
	}, nil
}

// NewWithHub creates a writer using an existing sentry Hub and options.
func NewWithHub(hub *sentry.Hub, opts ...WriterOption) (*Writer, error) {
	if hub == nil {
		return nil, errors.New("hub cannot be nil")
	}

	cfg := newDefaultConfig()
	for _, opt := range opts {
		opt.apply(&cfg)
	}

	levels := make(map[zerolog.Level]struct{}, len(cfg.levels))
	for _, lvl := range cfg.levels {
		levels[lvl] = struct{}{}
	}

	return &Writer{
		hub:             hub,
		levels:          levels,
		flushTimeout:    cfg.flushTimeout,
		withBreadcrumbs: cfg.breadcrumbs,
	}, nil
}

func newDefaultConfig() config {
	return config{
		levels: []zerolog.Level{
			zerolog.ErrorLevel,
			zerolog.FatalLevel,
			zerolog.PanicLevel,
		},
		sampleRate:   1.0,
		flushTimeout: 3 * time.Second,
	}
}
