package fileio

import (
	"fmt"
	"os"

	"github.com/toobprojects/go-commons/errx"
)

func IsSymlink(path string) (bool, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errx.Wrap(err, fmt.Sprintf("lstat %q", path))
	}
	return fi.Mode()&os.ModeSymlink != 0, nil
}

func ResolveSymlink(path string) (string, error) {
	dst, err := os.Readlink(path)
	if err != nil {
		return "", errx.Wrap(err, fmt.Sprintf("readlink %q", path))
	}
	return dst, nil
}
