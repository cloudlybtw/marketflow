package logger

import (
	"log/slog"
	"marketflow/internal/adapter/config"
	"os"
)

var logger *slog.Logger

func Set(config *config.App) {
	logger = slog.New(
		slog.NewTextHandler(os.Stderr, nil),
	)

	slog.SetDefault(logger)
}
