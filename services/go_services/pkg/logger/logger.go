package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func InitLogger(logLevel string) {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	Log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

}
