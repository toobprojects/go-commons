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
	Color bool         // enable ANSI colors in text mode (ignored for JSON)
}

var (
	mu         sync.RWMutex
	logger     *slog.Logger
	once       sync.Once
	defaultCfg = Config{
		Level: slog.LevelInfo,
		JSON:  false,
		Out:   os.Stdout,
		Color: false,
	}
	// currentCfg holds the last active config (initialized by Init).
	currentCfg = Config{}
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
	colorBlue   = "\033[34m"
)

// colorHandler wraps another slog.Handler and injects ANSI color
// codes into the log message based on the log level.
//
// This is only used when Config.Color is true and JSON is false.
type colorHandler struct {
	h slog.Handler
}

func (c *colorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return c.h.Enabled(ctx, level)
}

func (c *colorHandler) Handle(ctx context.Context, r slog.Record) error {
	// Work on a copy of the record so we don't mutate the original.
	rec := r.Clone()

	var color string
	switch {
	case rec.Level >= slog.LevelError:
		color = colorRed
	case rec.Level >= slog.LevelWarn:
		color = colorYellow
	case rec.Level >= slog.LevelInfo:
		color = colorGreen
	default:
		color = colorBlue
	}

	if color != "" {
		rec.Message = color + rec.Message + colorReset
	}

	return c.h.Handle(ctx, rec)
}

func (c *colorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &colorHandler{h: c.h.WithAttrs(attrs)}
}

func (c *colorHandler) WithGroup(name string) slog.Handler {
	return &colorHandler{h: c.h.WithGroup(name)}
}

func SetLogFile(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}

	// Use the last active config if available, otherwise fall back to defaultCfg.
	cfg := currentCfg
	// If currentCfg is zero (not initialized), use defaultCfg.
	if cfg.Out == nil && cfg.Level == nil {
		cfg = defaultCfg
	}
	cfg.Out = f

	// Reinitialize logger using the preserved config but with new output.
	Init(cfg)

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
		base := slog.NewTextHandler(cfg.Out, &slog.HandlerOptions{Level: cfg.Level})
		if cfg.Color {
			h = &colorHandler{h: base}
		} else {
			h = base
		}
	}

	l := slog.New(h)

	mu.Lock()
	logger = l
	// Remember the active config for future SetLogFile calls
	currentCfg = cfg
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
