package errx

import (
	"fmt"
	"io"

	"github.com/toobprojects/go-commons/logs"
)

// Wrap adds message context but preserves the original error for errors.Is/As.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// CloseQuietly logs a WARN if Close() fails (use in defers).
func CloseQuietly(c io.Closer, msg string, attrs ...any) {
	if c == nil {
		return
	}
	if err := c.Close(); err != nil {
		logs.Warn(msg, append(attrs, "error", err)...)
	}
}
