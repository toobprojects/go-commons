# text — String Helpers & Delimited Parsing ✂️🔡

Lightweight string utilities used across the codebase.  
The package focuses on **blank checks**, **safe comparisons**, **list membership**, **argument lookup**, and a clean, reusable **delimiter-based parser** for splitting text into structured sections.

---

# Overview

This package contains two groups of helpers:

1. **Basic string utilities** (`Blank`, `NotBlank`, `ListContains`, `Trim`, comparisons, etc.)  
2. **Delimited parsing utilities** that extract sections of text using start/end markers.

Both sets keep things fast, allocation-light, and dependency-free.

---

# 1. Basic String Utilities

## `Blank(s string) bool` / `NotBlank(s string) bool` 🧼  
Checks whether a string is empty or whitespace-only.

```go
text.Blank("   ")    // true
text.NotBlank("go")  // true
```

---

## `Trim(s string) string` ✂️
Trims leading and trailing whitespace.

```go
clean := text.Trim("   hello world   ")  // "hello world"
```

---

## `Equals(a, b string) bool` / `NotEquals(a, b string) bool` ⚖️
Strict byte-level equality.

```go
text.Equals("Go", "Go")      // true
text.NotEquals("Go", "Gopher") // true
```

---

## `EqualsIgnoreCase(a, b string) bool` / `NotEqualsIgnoreCase(a, b string) bool` 🫱🏽‍🫲🏾
Case-insensitive comparisons.

```go
text.EqualsIgnoreCase("Go", "go") // true
```

---

## `ListContains(list []string, val string) bool` 🔍
Case-sensitive membership in a list.

```go
langs := []string{"java", "go", "python"}
text.ListContains(langs, "go") // true
```

---

## `GetArg(args []string, index int) string` 🎯
Returns `args[index]`, or `""` if out of bounds.

```go
cmd := []string{"build", "--verbose"}
flag := text.GetArg(cmd, 1) // "--verbose"
missing := text.GetArg(cmd, 5) // ""
```

---

# 2. Delimited Parsing Utilities 📜

These helpers split a long text into **start/end delimited sections**.  
This is useful for config files, extracting fenced blocks, scanning for custom markers, etc.

---

## `FindDelimiterBlock(text, startDelim, endDelim string) (string, bool)` 🔦
Returns the **first** block between `startDelim` and `endDelim`.

- Returns `(blockContent, true)` if found
- Returns `("", false)` if missing or inverted
- Does **not** include the delimiters themselves

```go
raw := `
hello
---BEGIN---
inside block!
---END---
goodbye
`

block, ok := text.FindDelimiterBlock(raw, "---BEGIN---", "---END---")
fmt.Println(block) // "inside block!"
fmt.Println(ok)    // true
```

---

## `FindAllDelimiterBlocks(text, startDelim, endDelim string) []string` 📚
Extracts **all** blocks between `startDelim` and `endDelim` in order.

```go
raw := `
A
<<start>>
one
<<end>>
B
<<start>>
two
<<end>>
C
`
blocks := text.FindAllDelimiterBlocks(raw, "<<start>>", "<<end>>")
fmt.Println(blocks) // []string{"one", "two"}
```

---

## `SplitByDelimiter(text, delimiter string) []string` 🔪
Splits text by a specific delimiter (non-regex), trimming empty edges.

**Example**
```go
parts := text.SplitByDelimiter("a---b---c", "---")
// ["a", "b", "c"]
```

---

# Notes & Gotchas 🧠

- Delimiter functions are deterministic and never panic; missing markers simply return empty results.
- `FindDelimiterBlock` returns the first match only; use `FindAllDelimiterBlocks` for repeated sections.
- All helpers avoid regex for performance and clarity.
- For advanced templating or placeholder replacement, consider extending this package or building higher-level utilities.
