//go:build darwin

package fsutil

import "golang.org/x/sys/unix"

func copyFile(src, dst string) error {
	// clonefile creates a copy-on-write clone on APFS at near-zero cost.
	// Falls back to buffered copy if the filesystem doesn't support it
	// (e.g. exFAT SD card mounted as destination) or if dst already exists.
	if err := unix.Clonefile(src, dst, 0); err == nil {
		return nil
	}
	return copyBuffered(src, dst)
}
