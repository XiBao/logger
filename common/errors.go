package common

import "github.com/getsentry/sentry-go"

type ErrWithStackTrace struct {
	Stacktrace *sentry.Stacktrace `json:"stacktrace"`
	Err        string             `json:"error"`
}

func Stacktrace() *sentry.Stacktrace {
	const (
		currentModule = "github.com/XiBao/logger"
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
