package logger

import (
	"log/slog"
	"os"
)

// InitJSONLogger configures the default standard logger to output in JSON format.
func InitJSONLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	
	// Create the logger and set it as the global default for the application
	logger := slog.New(handler)
	slog.SetDefault(logger)
}