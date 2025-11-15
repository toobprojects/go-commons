# CLI — Native Command Runner ⚙️💨

A lightweight utility for running native OS commands with:
- Working-directory control  
- Extra environment variables  
- Captured vs streamed output  
- Optional logging  
- Convenience one-liner helpers  

The goal is to provide a **predictable and simple API** you can use in CLI projects, automation scripts, and internal tooling.

---

# Package Overview

The CLI package exposes three ways to run commands:

1. **`Run`** — the full-control, explicit API  
2. **`RunWithDefaults`** — the easy, “just give the output” API  
3. **`Exec*`** functions — backward-compatible helpers for quick one-liners  

All run through Go’s `os/exec` underneath with consistent behavior and optional logging.

---

# Types

## `type Options struct` 🤖

Controls exactly how a command executes.

- **Dir** — working directory  
- **Env** — extra environment key/value pairs  
- **CaptureOutput** — return combined stdout+stderr as a string  
- **Stdout / Stderr** — destinations when streaming output  
- **LogCommand** — logs command + args before running (uses logs package)  

Example of building `Options`:

```go
opts := cli.Options{
    Dir: "/tmp",
    Env: []string{"FOO=bar"},
    CaptureOutput: true,
    LogCommand: true,
}
```

---

# Functions

## `Run(ctx, command, args, opts) (string, error)` 🏁

**Most powerful API.**

- Logs before running if `opts.LogCommand == true`
- Applies working directory
- Applies environment overrides
- Either returns output or streams it

**Example — Capture Output**

```go
out, err := cli.Run(context.Background(),
    "bash",
    []string{"-lc", "echo $FOO && pwd"},
    cli.Options{
        Dir: "/tmp",
        Env: []string{"FOO=bar"},
        CaptureOutput: true,
        LogCommand: true,
    },
)
fmt.Println(out)
```

**Example — Stream Live Output**

```go
_, err := cli.Run(context.Background(),
    "go",
    []string{"test", "./..."},
    cli.Options{
        CaptureOutput: false,
        LogCommand: true,
    },
)
```

**Example — Timeout**

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

_, err := cli.Run(ctx, "sleep", []string{"10"}, cli.Options{CaptureOutput:false})
```

---

## `RunWithDefaults(ctx, command, args, logCommand)` 🚀

Easy API:

- Always captures output
- Uses current directory
- Inherits process env

**Example**

```go
out, err := cli.RunWithDefaults(context.Background(),
    "git",
    []string{"status", "--porcelain"},
    true,
)
fmt.Println(out)
```

---

## `Exec(command, args, targetPath, returnOutput)` 🧰

Backward-compatible helper:

- Uses `context.Background()`
- Uses `targetPath` as working directory
- Returns **empty string** on error
- Intended for quick scripts

**Example**

```go
ver := cli.Exec("go", []string{"version"}, "", true)
fmt.Println(ver)
```

---

## `ExecWithNativeLog(...)` 📣

Same as `Exec`, but always logs the command before executing.

**Example**

```go
cli.ExecWithNativeLog("docker", []string{"ps"}, "", false)
```

---

## `ExecScriptFile(path, dir, returnOutput)` 📝

Runs a script via `/bin/bash`.

**Example**

```go
cli.ExecScriptFile("./scripts/build.sh", ".", false)
```

---

# Practical Notes 🧠

- Use **`Run`** for robust automation & explicit error handling.
- Use **`RunWithDefaults`** when you just want the output fast.
- Use **`Exec*`** only for small shortcuts (they hide errors!).
- Captured output = stdout + stderr merged.
- For separate output streams → stream them manually.
- Combine with `context.WithTimeout` for long-running commands.

---

# Mini Cheat Sheet 📎

```go
// Simple captured output
out, _ := cli.RunWithDefaults(context.Background(), "uname", []string{"-a"}, true)

// Stream live output
cli.Run(context.Background(), "go", []string{"build"}, cli.Options{CaptureOutput:false})

// Run in specific directory with env
cli.Run(context.Background(), "bash", []string{"-lc", "echo $FOO"}, cli.Options{
    Dir: "/tmp",
    Env: []string{"FOO=bar"},
    CaptureOutput: true,
})

// Compatibility one-liner
cli.Exec("ls", []string{"-la"}, "/etc", false)
```
