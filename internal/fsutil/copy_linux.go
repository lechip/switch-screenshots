//go:build linux

package fsutil

import (
	"os"

	"golang.org/x/sys/unix"
)

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	// FICLONE is supported on btrfs/xfs; SD cards are typically exFAT so the
	// fallback path is expected in the common case.
	if err := unix.IoctlFileClone(int(out.Fd()), int(in.Fd())); err == nil {
		return out.Close()
	}
	out.Close()
	return copyBuffered(src, dst)
}
