package fileio

import (
	"fmt"
	"os"

	"github.com/toobprojects/go-commons/errx"
)

func EnsureDir(path string, perm os.FileMode) error {
	if err := os.MkdirAll(path, perm); err != nil {
		return errx.Wrap(err, fmt.Sprintf("mkdir -p %q", path))
	}
	return nil
}
