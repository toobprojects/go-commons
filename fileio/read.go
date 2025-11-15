package fileio

import (
	"fmt"
	"io"
	"os"

	"github.com/toobprojects/go-commons/errx"
)

func ReadFile(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, errx.Wrap(err, fmt.Sprintf("read file %q", path))
	}
	return b, nil
}

func ReadString(path string) (string, error) {
	b, err := ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// StreamRead opens the file and returns its bytes; example of using CloseQuietly.
func StreamRead(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errx.Wrap(err, fmt.Sprintf("open %q", path))
	}
	defer errx.CloseQuietly(f, "close file", "path", path)

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, errx.Wrap(err, fmt.Sprintf("read %q", path))
	}
	return b, nil
}
