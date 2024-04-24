package logger

import (
	"log/slog"
	"os"
)

func NewProvider() *slog.Logger {
	// create a logging level variable
	// the level is Info by default
	var loggingLevel = new(slog.LevelVar)

	// Pass the loggingLevel to the new logger being created so we can change it later at any time. Also adding source file information.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: loggingLevel}))

	// set the level to debug
	loggingLevel.Set(slog.LevelDebug)

	// // Set the logger for the application
	// slog.SetDefault(logger)

	return logger
}
