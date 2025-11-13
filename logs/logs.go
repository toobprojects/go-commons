package logs

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
)

type Config struct {
	Level slog.Leveler // slog.LevelDebug, slog.LevelInfo, etc.
	JSON  bool         // true = JSON handler, false = human-readable text
	Out   io.Writer    // usually os.Stdout or os.Stderr
}

var (
	mu         sync.RWMutex
	logger     *slog.Logger
	once       sync.Once
	defaultCfg = Config{
		Level: slog.LevelInfo,
		JSON:  false,
		Out:   os.Stdout,
	}
)

func SetLogFile(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}

	// Reinitialize logger using default config but override output
	Init(Config{
		Level: defaultCfg.Level,
		JSON:  defaultCfg.JSON,
		Out:   f,
	})

	return nil
}

// Init initializes the global logger.
// Safe to call multiple times, last call wins.
func Init(cfg Config) {
	if cfg.Out == nil {
		cfg.Out = defaultCfg.Out
	}
	if cfg.Level == nil {
		cfg.Level = defaultCfg.Level
	}

	var h slog.Handler
	if cfg.JSON {
		h = slog.NewJSONHandler(cfg.Out, &slog.HandlerOptions{Level: cfg.Level})
	} else {
		h = slog.NewTextHandler(cfg.Out, &slog.HandlerOptions{Level: cfg.Level})
	}

	l := slog.New(h)

	mu.Lock()
	logger = l
	mu.Unlock()
}

// get returns the current global logger, lazily initialized.
func get() *slog.Logger {
	once.Do(func() {
		Init(defaultCfg)
	})
	mu.RLock()
	defer mu.RUnlock()
	return logger
}

// With returns a new logger with additional attributes.
func With(args ...any) *slog.Logger {
	return get().With(args...)
}

// WithGroup creates a grouped logger for logical scoping (e.g. "cli", "fileio").
func WithGroup(name string) *slog.Logger {
	return get().WithGroup(name)
}

// --- Helper functions for convenience ---

func Debug(msg string, args ...any) {
	get().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	get().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	get().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	get().Error(msg, args...)
}

func DebugCtx(ctx context.Context, msg string, args ...any) {
	get().DebugContext(ctx, msg, args...)
}

func InfoCtx(ctx context.Context, msg string, args ...any) {
	get().InfoContext(ctx, msg, args...)
}

func WarnCtx(ctx context.Context, msg string, args ...any) {
	get().WarnContext(ctx, msg, args...)
}

func ErrorCtx(ctx context.Context, msg string, args ...any) {
	get().ErrorContext(ctx, msg, args...)
}
