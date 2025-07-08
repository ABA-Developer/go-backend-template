package utils

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(outputPath string) zerolog.Logger {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02/15:04:05(Z07:00)"}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Logger()

	return logger
}
