package fileio

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/toobprojects/go-commons/errx"
)

func Home() (string, error) {
	h, err := os.UserHomeDir()
	if err != nil {
		return "", errx.Wrap(err, "resolve home directory")
	}
	return h, nil
}

func ExpandHome(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	home, err := Home()
	if err != nil {
		return "", err
	}

	if path == "~" {
		return home, nil
	}
	return filepath.Join(home, path[2:]), nil
}

func Stat(path string) (os.FileInfo, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return nil, errx.Wrap(err, fmt.Sprintf("lstat %q", path))
	}
	return fi, nil
}

func Exists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if isNotExist(err) {
		return false, nil
	}
	return false, errx.Wrap(err, fmt.Sprintf("lstat %q", path))
}

func IsDir(path string) (bool, error) {
	fi, err := Stat(path)
	if err != nil {
		if isNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return fi.IsDir(), nil
}

func IsFile(path string) (bool, error) {
	fi, err := Stat(path)
	if err != nil {
		if isNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return fi.Mode().IsRegular(), nil
}

// isNotExist correctly detects "not found", including wrapped errors.
func isNotExist(err error) bool {
	if err == nil {
		return false
	}
	// Modern, wrap-aware check:
	if errors.Is(err, fs.ErrNotExist) {
		return true
	}
	// Compatibility check (also handles some OS-specific errors):
	return os.IsNotExist(err)
}
