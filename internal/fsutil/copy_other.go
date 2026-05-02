//go:build !darwin && !linux && !windows

package fsutil

func copyFile(src, dst string) error {
	return copyBuffered(src, dst)
}
