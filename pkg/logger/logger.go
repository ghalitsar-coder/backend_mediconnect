package logger

import (
	"log/slog"
	"os"
)

// Init configures the global structured logger based on the environment.
// In production, emits JSON logs. In development, emits human-readable text logs.
func Init(env string) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		AddSource: env == "production",
		Level:     slog.LevelDebug,
	}

	if env == "production" {
		opts.Level = slog.LevelInfo
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
