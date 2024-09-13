package zlogsentry

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
	"unsafe"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
)

type ErrWithStackTrace struct {
	Stacktrace *sentry.Stacktrace `json:"stacktrace"`
	Err        string             `json:"error"`
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
	levels          map[zerolog.Level]struct{}
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

	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Category: category,
		Message:  event.Message,
		Level:    event.Level,
		Data:     event.Extra,
	})
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

	sentry.CaptureEvent(event)

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

	sentry.CaptureEvent(event)
	return
}

func (w *Writer) Close() error {
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
	levels      []zerolog.Level
	breadcrumbs bool
}

// WithLevels configures zerolog levels that have to be sent to Sentry. Default levels are error, fatal, panic
func WithLevels(levels ...zerolog.Level) WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.levels = levels
	})
}

// WithBreadcrumbs enables sentry client breadcrumbs.
func WithBreadcrumbs() WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.breadcrumbs = true
	})
}

func New(opts ...WriterOption) (*Writer, error) {
	cfg := newDefaultConfig()
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.apply(&cfg)
		}
	}

	levels := make(map[zerolog.Level]struct{}, len(cfg.levels))
	for _, lvl := range cfg.levels {
		levels[lvl] = struct{}{}
	}

	return &Writer{
		levels:          levels,
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
	}
}
