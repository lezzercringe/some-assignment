package log

import (
	"fmt"
	"log/slog"
	"order-persistor/internal/config"
	"os"
)

var levels = map[string]slog.Leveler{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func New(cfg config.Log) (*slog.Logger, error) {
	level, ok := levels[cfg.Level]
	if !ok {
		return nil, fmt.Errorf("unknown log level: %s", cfg.Level)
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})), nil
}
