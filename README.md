# 🚀 Go Commons

A reusable, production-ready toolkit of **Go utilities** for everyday development.  
It eliminates boilerplate, enforces consistency, and gives you clean, battle-tested helpers for CLI tools, services, and automation scripts.

---

# ✨ Overview

**Go Commons** wraps the Go standard library with *ergonomic*, *predictable*, and *developer-friendly* utilities.  
Think of it as your “daily-driver toolbox” — simple, lightweight, and practical.

### 🔧 Included Capabilities

- 📁 **File I/O** — read/write, copy, append, directory helpers, symlinks, path checks  
- 📝 **Config Parsing** — JSON / YAML decoding with strict mode + env variable expansion  
- ⚙️ **CLI Execution** — run native commands with context, working dir, env overrides, logged output  
- 📊 **Structured Logging** — slog-based logger (JSON or colorized text) with global init + helpers  
- ✂️ **Text Helpers** — blank checks, comparisons, argument lookup, delimiter block extraction  

Everything follows a consistent design philosophy:  
**small, no-nonsense helpers that behave the way you expect.**

---

# 🧱 Tech Stack

- **Go:** 1.25.x  
- **Primary dependency:** `gopkg.in/yaml.v3` (YAML parsing)

Zero heavy dependencies. Zero magic. Just clean utilities.

---

# 📦 Packages & Documentation

| Package | Description                                                                         | Documentation                      |
| ------- | ----------------------------------------------------------------------------------- | ---------------------------------- |
| ⚙️ CLI     | Execute native system commands with flexible options (cwd, env, output, logging)   | [docs/cli.md](./docs/cli.md)       |
| 📁 File IO | Read/write, copy, append, paths, symlinks, config parsing (JSON/YAML)             | [docs/fileio.md](./docs/fileio.md) |
| 📊 Logs    | Structured logging on top of slog (JSON or colorized text output)                 | [docs/logs.md](./docs/logs.md)     |
| ✂️ Text    | String utilities: blank checks, comparisons, list search, arg parsing, delimiters | [docs/text.md](./docs/text.md)     |

> **Note:** `machine` and `maven` packages are intentionally excluded from documentation for now.

---

# 🧭 Philosophy

- Keep things **small**, **simple**, and **explicit**  
- Provide consistent behavior across utilities  
- Wrap the standard library without hiding it  
- Handle errors clearly with helpful context  

Use **Go Commons** to accelerate CLI development, config-driven scripts, microservices, or any internal tooling where consistency matters.
