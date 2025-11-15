package fileio

import (
	"fmt"
	"os"

	"github.com/toobprojects/go-commons/errx"
)

// WriteFile writes with the provided permissions (e.g., 0o644).
func WriteFile(path string, data []byte, perm os.FileMode) error {
	if err := os.WriteFile(path, data, perm); err != nil {
		return errx.Wrap(err, fmt.Sprintf("write file %q", path))
	}
	return nil
}

// AppendFile appends, creating the file if needed.
func AppendFile(path string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, perm)
	if err != nil {
		return errx.Wrap(err, fmt.Sprintf("open+append %q", path))
	}
	defer errx.CloseQuietly(f, "close file", "path", path)

	if _, err := f.Write(data); err != nil {
		return errx.Wrap(err, fmt.Sprintf("append %q", path))
	}
	return nil
}
