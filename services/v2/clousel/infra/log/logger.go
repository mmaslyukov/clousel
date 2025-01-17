package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func New() *zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()
	return &logger

}
