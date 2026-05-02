package fsutil

import (
	"io"
	"os"
	"path/filepath"
)

// Copy copies src to dst, creating parent directories as needed.
// A fast reflink/clone is attempted on darwin (clonefile) and linux (FICLONE
// ioctl); all other platforms and fallbacks use a buffered io.Copy.
func Copy(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return copyFile(src, dst)
}

// copyBuffered is the portable fallback called by platform-specific files.
func copyBuffered(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	if _, err = io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
