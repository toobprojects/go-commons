package fileio

import (
	"fmt"
	"io"
	"os"

	"github.com/toobprojects/go-commons/errx"
)

func CopyFile(src, dst string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return errx.Wrap(err, fmt.Sprintf("open src %q", src))
	}
	defer errx.CloseQuietly(in, "close src", "path", src)

	// create/truncate destination
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return errx.Wrap(err, fmt.Sprintf("create dst %q", dst))
	}
	defer errx.CloseQuietly(out, "close dst", "path", dst)

	if _, err := io.Copy(out, in); err != nil {
		return errx.Wrap(err, fmt.Sprintf("copy %q -> %q", src, dst))
	}
	return nil
}
