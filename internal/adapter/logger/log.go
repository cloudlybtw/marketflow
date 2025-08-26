package logger

import (
	"log/slog"
	"os"

	"marketflow/internal/adapter/config"
)

var logger *slog.Logger

func Set(config *config.App) {
	logger = slog.New(
		slog.NewTextHandler(os.Stderr, nil),
	)

	slog.SetDefault(logger)
}
