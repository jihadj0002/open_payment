package logger

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(level string) {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	log.Logger = zerolog.New(output).
		Level(lvl).
		With().
		Timestamp().
		Caller().
		Logger()
}

func Ctx(ctx context.Context) *zerolog.Logger {
	return log.Ctx(ctx)
}

func Redact(s string) string {
	if len(s) > 6 {
		return s[:4] + strings.Repeat("*", len(s)-6) + s[len(s)-2:]
	}
	return strings.Repeat("*", len(s))
}
