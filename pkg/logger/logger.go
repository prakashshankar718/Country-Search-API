package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	log  *slog.Logger
	once sync.Once
)

func Init(level slog.Level) {
	once.Do(func() {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})

		log = slog.New(handler)
	})
}

func Log() *slog.Logger {
	if log == nil {
		Init(slog.LevelInfo)
	}
	return log
}
