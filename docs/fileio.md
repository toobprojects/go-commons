# fileio — Files, Paths & Parsing 📁🧭

Utility helpers for working with files and directories, copying, symlinks, path checks, and decoding JSON/YAML into your own structs. All functions add clear error context using `errx.Wrap`, and path checks correctly handle wrapped `fs.ErrNotExist`.

---

# Quick Map

- **Read/Write:** `ReadFile`, `ReadString`, `StreamRead`, `WriteFile`, `AppendFile`
- **Dirs/Paths:** `EnsureDir`, `Home`, `ExpandHome`, `Stat`, `Exists`, `IsDir`, `IsFile`
- **Copy/Symlinks:** `CopyFile`, `IsSymlink`, `ResolveSymlink`
- **Parsing (JSON/YAML):** `ParseFile[T]`, `ParseReader[T]`, `ParseBytes[T]`, `ParseString[T]`, `ParseStringAuto[T]`, with `WithStrict`, `WithEnvExpand`

---

# Reading & Writing 🔤

## `ReadFile(path string) ([]byte, error)`
Reads a whole file. Errors are wrapped like `read file "path": ...`.

```go
b, err := fileio.ReadFile("README.md")
if err != nil { panic(err) }
fmt.Printf("%d bytes\n", len(b))
```

## `ReadString(path string) (string, error)`

Like `ReadFile`, but returns a `string`.

```go
s, err := fileio.ReadString("README.md")
if err != nil { panic(err) }
fmt.Println(s[:80])
```

## `StreamRead(path string) ([]byte, error)`

Opens and streams the contents with `io.ReadAll`; file is closed via `errx.CloseQuietly`.

```go
b, err := fileio.StreamRead("README.md")
if err != nil { panic(err) }
```

## `WriteFile(path string, data []byte, perm os.FileMode) error`

Writes bytes with explicit permissions (e.g., `0o644`).

```go
err := fileio.WriteFile("out.txt", []byte("hello\n"), 0o644)
```

## `AppendFile(path string, data []byte, perm os.FileMode) error`

Appends to a file (creates it if missing).

```go
_ = fileio.AppendFile("app.log", []byte("started\n"), 0o644)
```

---

# Directories & Paths 🗂️

## `EnsureDir(path string, perm os.FileMode) error`

`mkdir -p` behavior with error wrapping.

```go
_ = fileio.EnsureDir("./build/reports", 0o755)
```

## `Home() (string, error)`

Resolves the current user’s home directory, wrapped on error.

```go
home, err := fileio.Home()
```

## `ExpandHome(path string) (string, error)`

Expands `~` or `~/...` using `Home()`; returns input unchanged if it doesn’t start with `~`.

```go
abs, _ := fileio.ExpandHome("~/projects/go-commons")
```

## `Stat(path string) (fs.FileInfo, error)`

`os.Lstat` with wrapped errors (safe for symlinks). Used by other helpers.

```go
fi, err := fileio.Stat("go.mod")
if err == nil { fmt.Println(fi.Mode()) }
```

## `Exists(path string) (bool, error)`

`true` if `os.Lstat` succeeds; returns `(false, nil)` if not-exist (even when wrapped), otherwise wraps and returns error.

```go
ok, err := fileio.Exists("go.mod") // ok=true if present
```

## `IsDir(path string) (bool, error)` / `IsFile(path string) (bool, error)`

Use `Stat` + mode checks; return `(false, nil)` when path doesn’t exist.

```go
isDir, _  := fileio.IsDir("some/path")
isFile, _ := fileio.IsFile("go.mod")
```

---

# Copy & Symlinks 🔗

## `CopyFile(src, dst string, perm os.FileMode) error`

Opens `src`, creates/truncates `dst` with `perm`, then `io.Copy`.

```go
_ = fileio.CopyFile("a.txt", "b.txt", 0o644)
```

## `IsSymlink(path string) (bool, error)`

`os.Lstat` + mode check; `(false, nil)` if the path doesn’t exist.

```go
if yes, _ := fileio.IsSymlink("latest"); yes { fmt.Println("is link") }
```

## `ResolveSymlink(path string) (string, error)`

Returns `os.Readlink(path)` (the symlink target).

```go
dst, err := fileio.ResolveSymlink("latest")
```

---

# JSON/YAML Parsing 🧩

Decode JSON or YAML into your own structs with optional **strict mode** and **env expansion**.

## Options

- `WithStrict()` — error on unknown fields
- `WithEnvExpand()` — expand `${VAR}` before decoding (from process env)

## Public APIs

- `ParseFile[T any](path string, opts ...Option) (T, error)`
  - Uses extension: `.json` → JSON; `.yaml/.yml` → YAML
- `ParseReader[T any](r io.Reader, ext string, opts ...Option) (T, error)`
- `ParseBytes[T any](data []byte, ext string, opts ...Option) (T, error)`
- `ParseString[T any](content, ext string, opts ...Option) (T, error)`
- `ParseStringAuto[T any](content string, opts ...Option) (T, error)`
  - Sniffs: starts with `{` or `[` ⇒ JSON, otherwise YAML

### Example — From File (strict + env expand)

```go
type AppCfg struct {
    Name string `json:"name" yaml:"name"`
    Port int    `json:"port" yaml:"port"`
}

cfg, err := fileio.ParseFile[AppCfg]("app.yaml",
    fileio.WithStrict(),
    fileio.WithEnvExpand(),
)
```

### Example — From String (auto-sniff)

```go
raw := `
name: ${APP_NAME}
port: 8080
`
cfg2, err := fileio.ParseStringAuto[AppCfg](raw, fileio.WithEnvExpand())
```

### Example — From Bytes (force extension)

```go
b := []byte(`{"name":"svc","port":8080}`)
cfg3, err := fileio.ParseBytes[AppCfg](b, ".json")
```

---

# Notes & Gotchas 🧠

- All file errors are wrapped with the operation context (e.g., `read "path"`), making logs/searching easier.
- Existence checks properly handle **wrapped** `fs.ErrNotExist` via `errors.Is`.
- `Parse*` merges env expansion **before** decoding; keep secrets out of logs.
- `ParseStringAuto` uses a simple leading-character heuristic: `{`/`[` ⇒ JSON; else YAML.
- If you need separate stdout/stderr or process execution, see the `cli` package.

```
```markdown
# fileio — Files, Paths & Parsing 📁🧭

Utility helpers for working with files and directories, copying, symlinks, path checks, and decoding JSON/YAML into your own structs. All functions add clear error context using `errx.Wrap`, and path checks correctly handle wrapped `fs.ErrNotExist`.

---

# Quick Map

- **Read/Write:** `ReadFile`, `ReadString`, `StreamRead`, `WriteFile`, `AppendFile`
- **Dirs/Paths:** `EnsureDir`, `Home`, `ExpandHome`, `Stat`, `Exists`, `IsDir`, `IsFile`
- **Copy/Symlinks:** `CopyFile`, `IsSymlink`, `ResolveSymlink`
- **Parsing (JSON/YAML):** `ParseFile[T]`, `ParseReader[T]`, `ParseBytes[T]`, `ParseString[T]`, `ParseStringAuto[T]`, with `WithStrict`, `WithEnvExpand`

---

# Reading & Writing 🔤

## `ReadFile(path string) ([]byte, error)`
Reads a whole file. Errors are wrapped like `read file "path": ...`.

```go
b, err := fileio.ReadFile("README.md")
if err != nil { panic(err) }
fmt.Printf("%d bytes\n", len(b))
```

## `ReadString(path string) (string, error)`

Like `ReadFile`, but returns a `string`.

```go
s, err := fileio.ReadString("README.md")
if err != nil { panic(err) }
fmt.Println(s[:80])
```

## `StreamRead(path string) ([]byte, error)`

Opens and streams the contents with `io.ReadAll`; file is closed via `errx.CloseQuietly`.

```go
b, err := fileio.StreamRead("README.md")
if err != nil { panic(err) }
```

## `WriteFile(path string, data []byte, perm os.FileMode) error`

Writes bytes with explicit permissions (e.g., `0o644`).

```go
err := fileio.WriteFile("out.txt", []byte("hello\n"), 0o644)
```

## `AppendFile(path string, data []byte, perm os.FileMode) error`

Appends to a file (creates it if missing).

```go
_ = fileio.AppendFile("app.log", []byte("started\n"), 0o644)
```

---

# Directories & Paths 🗂️

## `EnsureDir(path string, perm os.FileMode) error`

`mkdir -p` behavior with error wrapping.

```go
_ = fileio.EnsureDir("./build/reports", 0o755)
```

## `Home() (string, error)`

Resolves the current user’s home directory, wrapped on error.

```go
home, err := fileio.Home()
```

## `ExpandHome(path string) (string, error)`

Expands `~` or `~/...` using `Home()`; returns input unchanged if it doesn’t start with `~`.

```go
abs, _ := fileio.ExpandHome("~/projects/go-commons")
```

## `Stat(path string) (fs.FileInfo, error)`

`os.Lstat` with wrapped errors (safe for symlinks). Used by other helpers.

```go
fi, err := fileio.Stat("go.mod")
if err == nil { fmt.Println(fi.Mode()) }
```

## `Exists(path string) (bool, error)`

`true` if `os.Lstat` succeeds; returns `(false, nil)` if not-exist (even when wrapped), otherwise wraps and returns error.

```go
ok, err := fileio.Exists("go.mod") // ok=true if present
```

## `IsDir(path string) (bool, error)` / `IsFile(path string) (bool, error)`

Use `Stat` + mode checks; return `(false, nil)` when path doesn’t exist.

```go
isDir, _  := fileio.IsDir("some/path")
isFile, _ := fileio.IsFile("go.mod")
```

---

# Copy & Symlinks 🔗

## `CopyFile(src, dst string, perm os.FileMode) error`

Opens `src`, creates/truncates `dst` with `perm`, then `io.Copy`.

```go
_ = fileio.CopyFile("a.txt", "b.txt", 0o644)
```

## `IsSymlink(path string) (bool, error)`

`os.Lstat` + mode check; `(false, nil)` if the path doesn’t exist.

```go
if yes, _ := fileio.IsSymlink("latest"); yes { fmt.Println("is link") }
```

## `ResolveSymlink(path string) (string, error)`

Returns `os.Readlink(path)` (the symlink target).

```go
dst, err := fileio.ResolveSymlink("latest")
```

---

# JSON/YAML Parsing 🧩

Decode JSON or YAML into your own structs with optional **strict mode** and **env expansion**.

## Options

- `WithStrict()` — error on unknown fields
- `WithEnvExpand()` — expand `${VAR}` before decoding (from process env)

## Public APIs

- `ParseFile[T any](path string, opts ...Option) (T, error)`
  - Uses extension: `.json` → JSON; `.yaml/.yml` → YAML
- `ParseReader[T any](r io.Reader, ext string, opts ...Option) (T, error)`
- `ParseBytes[T any](data []byte, ext string, opts ...Option) (T, error)`
- `ParseString[T any](content, ext string, opts ...Option) (T, error)`
- `ParseStringAuto[T any](content string, opts ...Option) (T, error)`
  - Sniffs: starts with `{` or `[` ⇒ JSON, otherwise YAML

### Example — From File (strict + env expand)

```go
type AppCfg struct {
    Name string `json:"name" yaml:"name"`
    Port int    `json:"port" yaml:"port"`
}

cfg, err := fileio.ParseFile[AppCfg]("app.yaml",
    fileio.WithStrict(),
    fileio.WithEnvExpand(),
)
```

### Example — From String (auto-sniff)

```go
raw := `
name: ${APP_NAME}
port: 8080
`
cfg2, err := fileio.ParseStringAuto[AppCfg](raw, fileio.WithEnvExpand())
```

### Example — From Bytes (force extension)

```go
b := []byte(`{"name":"svc","port":8080}`)
cfg3, err := fileio.ParseBytes[AppCfg](b, ".json")
```

---

# Notes & Gotchas 🧠

- All file errors are wrapped with the operation context (e.g., `read "path"`), making logs/searching easier.
- Existence checks properly handle **wrapped** `fs.ErrNotExist` via `errors.Is`.
- `Parse*` merges env expansion **before** decoding; keep secrets out of logs.
- `ParseStringAuto` uses a simple leading-character heuristic: `{`/`[` ⇒ JSON; else YAML.
- If you need separate stdout/stderr or process execution, see the `cli` package.


