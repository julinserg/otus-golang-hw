package logger

import (
	"io"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	fileLogger *zerolog.Logger
}

func New(level string, w io.Writer) *Logger {
	if level == "INFO" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if level == "ERROR" {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if level == "WARN" {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if level == "DEBUG" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	fileLogger := zerolog.New(w).With().Timestamp().Logger()
	return &Logger{&fileLogger}
}

func (l Logger) Info(msg string) {
	l.fileLogger.Info().Msg(msg)
	log.Info().Msg(msg)
}

func (l Logger) Error(msg string) {
	l.fileLogger.Error().Msg(msg)
	log.Error().Msg(msg)
}

func (l Logger) Debug(msg string) {
	l.fileLogger.Debug().Msg(msg)
	log.Debug().Msg(msg)
}

func (l Logger) Warn(msg string) {
	l.fileLogger.Warn().Msg(msg)
	log.Warn().Msg(msg)
}
