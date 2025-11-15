// Package config provides decoding helpers for JSON/YAML into user-defined structs.
// It supports parsing from files, readers, bytes, and strings, with optional
// strict mode (error on unknown fields) and environment variable expansion.
//
// Usage:
//
//	type AppCfg struct {
//	  Name string `json:"name" yaml:"name"`
//	  Port int    `json:"port" yaml:"port"`
//	}
//
//	cfg, err := config.ParseFile[AppCfg]("app.yaml", config.WithStrict())
package fileio

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// =====================
// Options
// =====================

type parseOptions struct {
	strict     bool   // fail on unknown fields
	envExpand  bool   // expand ${VAR} before parsing
	readerName string // for error context (e.g., filename)
}

// Option configures parsing behavior.
type Option func(*parseOptions)

// WithStrict makes the parser fail on unknown fields:
// - JSON: DisallowUnknownFields()
// - YAML: KnownFields(true)
func WithStrict() Option { return func(o *parseOptions) { o.strict = true } }

// WithEnvExpand expands ${VAR} occurrences in the raw content prior to decoding.
func WithEnvExpand() Option { return func(o *parseOptions) { o.envExpand = true } }

// =====================
/* Public API */
// =====================

// ParseFile reads a JSON or YAML file into T based on file extension:
//
//	.json => JSON, .yaml/.yml => YAML
func ParseFile[T any](path string, opts ...Option) (T, error) {
	var zero T

	data, err := os.ReadFile(path)
	if err != nil {
		return zero, fmt.Errorf("read %q: %w", path, err)
	}
	return ParseBytes[T](data, filepath.Ext(path), append(opts, withReaderName(path))...)
}

// ParseReader reads from r as JSON/YAML based on ext (".json", ".yaml", ".yml").
// Useful when the source is embed.FS, HTTP, etc.
func ParseReader[T any](r io.Reader, ext string, opts ...Option) (T, error) {
	var zero T

	data, err := io.ReadAll(r)
	if err != nil {
		return zero, fmt.Errorf("read: %w", err)
	}
	return ParseBytes[T](data, ext, opts...)
}

// ParseBytes parses raw bytes as JSON/YAML based on ext.
func ParseBytes[T any](data []byte, ext string, opts ...Option) (T, error) {
	var zero T

	cfg := parseOptions{}
	for _, o := range opts {
		o(&cfg)
	}
	if cfg.envExpand {
		data = []byte(os.ExpandEnv(string(data)))
	}

	switch strings.ToLower(normExt(ext)) {
	case ".json":
		return parseJSON[T](data, cfg)
	case ".yaml", ".yml":
		return parseYAML[T](data, cfg)
	default:
		return zero, fmt.Errorf("%w: %s (expected .json, .yaml, .yml)", ErrUnsupportedExt, ext)
	}
}

// ParseString parses JSON/YAML from a string using the given extension.
func ParseString[T any](content, ext string, opts ...Option) (T, error) {
	return ParseBytes[T]([]byte(content), ext, opts...)
}

// ParseStringAuto sniffs JSON ('{' or '[') vs YAML (otherwise) and parses accordingly.
func ParseStringAuto[T any](content string, opts ...Option) (T, error) {
	ext := sniffExt(content)
	return ParseBytes[T]([]byte(content), ext, opts...)
}

// =====================
// Errors
// =====================

var ErrUnsupportedExt = errors.New("unsupported file extension")

// =====================
// Internals
// =====================

func withReaderName(name string) Option {
	return func(o *parseOptions) { o.readerName = name }
}

func parseJSON[T any](data []byte, cfg parseOptions) (T, error) {
	var out T

	dec := json.NewDecoder(bytes.NewReader(data))
	if cfg.strict {
		dec.DisallowUnknownFields()
	}
	if err := dec.Decode(&out); err != nil {
		return out, wrapWhere("json", cfg.readerName, err)
	}
	// If strict, try to detect trailing garbage beyond the first JSON value.
	if cfg.strict {
		// Consume whitespace and see if there's another token/value.
		if dec.More() {
			return out, wrapWhere("json", cfg.readerName, errors.New("trailing data after JSON payload"))
		}
	}
	return out, nil
}

func parseYAML[T any](data []byte, cfg parseOptions) (T, error) {
	var out T

	dec := yaml.NewDecoder(bytes.NewReader(data))
	if cfg.strict {
		dec.KnownFields(true)
	}
	if err := dec.Decode(&out); err != nil {
		return out, wrapWhere("yaml", cfg.readerName, err)
	}
	return out, nil
}

func wrapWhere(kind, name string, err error) error {
	where := kind
	if name != "" {
		where = fmt.Sprintf("%s(%s)", kind, name)
	}
	return fmt.Errorf("%s: %w", where, err)
}

func normExt(ext string) string {
	e := strings.ToLower(strings.TrimSpace(ext))
	switch e {
	case "json", ".json":
		return ".json"
	case "yaml", ".yaml", "yml", ".yml":
		return ".yaml"
	default:
		return e
	}
}

func sniffExt(s string) string {
	clean := strings.TrimLeft(stripBOM(s), " \t\r\n")
	if len(clean) == 0 {
		return ".yaml" // permissive default
	}
	switch clean[0] {
	case '{', '[':
		return ".json"
	default:
		return ".yaml"
	}
}

func stripBOM(s string) string {
	if strings.HasPrefix(s, "\uFEFF") {
		return s[len("\uFEFF"):]
	}
	return s
}
