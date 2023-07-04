package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	log = zerolog.New(output).With().Timestamp().Logger()
}

func Info(message string) {
	log.Info().Msg(message)
}

func Infof(format string, args ...interface{}) {
	log.Info().Msgf(format, args...)
}

func Warn(err error) {
	log.Warn().Msg(err.Error())
}

func Warnf(format string, args ...interface{}) {
	log.Warn().Msgf(format, args...)
}

func Error(err error) {
	log.Error().Msg(err.Error())
}

func Errorf(format string, args ...interface{}) {
	log.Error().Msgf(format, args...)
}

func Fatal(err error) {
	log.Fatal().Msg(err.Error())
}

func Fatalf(format string, args ...interface{}) {
	log.Fatal().Msgf(format, args...)
}
