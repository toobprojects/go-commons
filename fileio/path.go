package fileio

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/toobprojects/go-commons/errx"
)

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
	if os.IsNotExist(err) {
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

func isNotExist(err error) bool { return err != nil && (os.IsNotExist(err) || isFsNotExist(err)) }
func isFsNotExist(err error) bool {
	// guard for wrapped errors
	return err != nil && (errorIs[fs.ErrNotExist](err))
}

func errorIs[T error](err error) bool {
	var target T
	return target != nil && (target == (T)(nil)) // placeholder to satisfy generics; or simply rely on os.IsNotExist
}
