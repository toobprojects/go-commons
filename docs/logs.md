# logs — Structured Logging on top of slog 🧯📊

A thin, ergonomic wrapper around Go’s `log/slog` that gives you:
- One-time **global initialization** (text or JSON; colored text for dev)
- Simple **helpers**: `Debug/Info/Warn/Error` and their `*Ctx` variants
- Reusable **scoped loggers** via `With(...)` and `WithGroup(...)`
- Optional redirection to a **log file** with `SetLogFile(...)`

Use JSON for production ingestion (ELK/Loki/etc.) and colored text locally.

---

# Types

## `type Config struct` ⚙️
Controls how the global logger is initialized.

- `Level slog.Leveler` — e.g., `slog.LevelDebug`, `slog.LevelInfo`
- `JSON bool` — `true` uses the JSON handler; `false` uses a human-readable text handler
- `Out io.Writer` — `os.Stdout`, `os.Stderr`, or any writer you provide
- `Color bool` — enable ANSI colors in text mode (ignored for JSON)

**Example**
```go
logs.Init(logs.Config{
    Level: slog.LevelInfo,
    JSON:  false,        // text handler
    Out:   os.Stdout,    // print to console
    Color: true,         // colorized levels
})
```

---

# Initialization & Output

## `Init(cfg Config)`
Initializes (or re-initializes) the global `*slog.Logger`. Safe to call during startup; last call wins.

```go
logs.Init(logs.Config{ Level: slog.LevelDebug, JSON: true, Out: os.Stdout })
```

## `SetLogFile(path string) error`
Redirects output to (and opens) an **append-only** log file, keeping your current JSON/text mode and level.

```go
if err := logs.SetLogFile("app.log"); err != nil {
    panic(err)
}
logs.Info("file logging enabled", "path", "app.log")
```

---

# Emitting Logs

## Level Helpers
- `Debug(msg string, args ...any)`
- `Info(msg string, args ...any)`
- `Warn(msg string, args ...any)`
- `Error(msg string, args ...any)`

```go
logs.Info("server started", "addr", ":8080", "pid", os.Getpid())
logs.Warn("cache miss", "key", key)
logs.Error("db connect failed", "err", err)
```

## Context Variants
- `DebugCtx(ctx context.Context, msg string, args ...any)`
- `InfoCtx(ctx context.Context, msg string, args ...any)`
- `WarnCtx(ctx context.Context, msg string, args ...any)`
- `ErrorCtx(ctx context.Context, msg string, args ...any)`

```go
ctx := context.Background()
logs.InfoCtx(ctx, "handled request", "path", r.URL.Path, "ms", dur.Milliseconds())
```

---

# Scoped / Structured Logging

## `With(args ...any) *slog.Logger`
Returns a logger that always includes the given attributes. Great for adding component/service tags once and reusing.

```go
authLog := logs.With("component", "auth")
authLog.Info("login ok", "user", uid)
authLog.Warn("token near expiry", "sub", sub, "exp", expUnix)
```

## `WithGroup(name string) *slog.Logger`
Groups subsequent attributes under a key when rendered (especially useful in JSON).

```go
reqLog := logs.WithGroup("req").With("id", reqID)
reqLog.Info("validated")
```

---

# Practical Notes 🧠

- Prefer **JSON** mode in production (machine-friendly) and **text+Color** locally.
- Keep logs **structured**: `logs.Info("msg", "key", val, ...)` makes filtering easy.
- Use the `*Ctx` variants if you need to propagate request-scoped data/middleware cancelation (handlers that inspect `ctx`).
- `SetLogFile` is a convenience for single-process apps; for containers, prefer stdout and let the platform aggregate.
- Under the hood, the package holds a lazily-initialized global `*slog.Logger` guarded by a mutex to be concurrency-safe.

---

# Quick Cheat Sheet 📎

```go
// 1) Initialize once
logs.Init(logs.Config{
    Level: slog.LevelInfo,
    JSON:  false,     // switch to true in prod if you prefer JSON
    Out:   os.Stdout,
    Color: true,
})

// 2) Emit logs
logs.Info("booting", "version", buildVersion)
logs.Debug("detail", "cfg", cfg)

// 3) Scoped logger
svc := logs.With("svc", "billing")
svc.Info("reconcile started")

// 4) Grouped attributes
req := logs.WithGroup("req").With("id", reqID)
req.Error("failed", "err", err)

// 5) Context-aware
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
logs.InfoCtx(ctx, "finished job", "count", n)
```