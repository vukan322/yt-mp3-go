package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func Setup(env string) {
	var logLevel slog.Level
	switch env {
	case "production":
		logLevel = slog.LevelInfo
	default:
		logLevel = slog.LevelDebug
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      logLevel,
		TimeFormat: "02.01.2006 / 15:04:05",
	}))

	slog.SetDefault(logger)
	slog.Debug("logger initialized", "environment", env, "level", logLevel.String())
}
