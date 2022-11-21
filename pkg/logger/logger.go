package logger

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
}

var (
	logger Logger
	once   sync.Once
)

func Get() *Logger {
	once.Do(func() {
		// By default create console writer
		writers := []io.Writer{zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp}}

		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		multiWriters := io.MultiWriter(writers...)

		zeroLogger := zerolog.New(multiWriters).With().Timestamp().Logger()

		logger = Logger{&zeroLogger}
	})

	return &logger
}
